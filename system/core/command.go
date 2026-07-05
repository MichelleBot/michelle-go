package core

import (
	"fmt"
	"sort"
	"time"
)

type HandlerFunc func(ctx *Ptz) error

type CommandQuota struct {
	Enabled bool
	Cost    int
}

func PerUserQuota(cost int) *CommandQuota {
	return &CommandQuota{Enabled: true, Cost: cost}
}

func EnsureQuotaAvailable(ptz *Ptz, amount int) error {
	if amount <= 0 || ptz == nil || ptz.Bot == nil || ptz.Bot.Users == nil {
		return nil
	}
	if ptz.IsOwner() {
		return nil
	}

	userID := ptz.GetPhoneJID().User
	if ptz.Bot.Users.IsPremium(userID) {
		return nil
	}

	profile := ptz.Bot.Users.GetUserProfile(userID)
	if profile.LimitBalance < amount {
		return ptz.ReplyText("Limit kamu habis gunakan command buylimit atau hubungi owner")
	}
	return nil
}

func ConsumeQuota(ptz *Ptz, amount int) error {
	if amount <= 0 || ptz == nil || ptz.Bot == nil || ptz.Bot.Users == nil {
		return nil
	}
	if ptz.IsOwner() {
		return nil
	}

	userID := ptz.GetPhoneJID().User
	if ptz.Bot.Users.IsPremium(userID) {
		return nil
	}

	ok, err := ptz.Bot.Users.ConsumeLimit(userID, amount)
	if err != nil {
		return err
	}
	if !ok {
		return ptz.ReplyText("Limit kamu habis gunakan command buylimit atau hubungi owner")
	}
	return nil
}

type CommandLimit struct {
	Enabled bool
	Max     int
	Window  time.Duration
}

func PerUserLimit(max int, window time.Duration) *CommandLimit {
	return &CommandLimit{
		Enabled: true,
		Max:     max,
		Window:  window,
	}
}

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Category    string
	OwnerOnly   bool
	GroupOnly   bool
	AdminOnly   bool
	BotAdmin    bool
	Quota       *CommandQuota
	Limit       *CommandLimit
	Handler     HandlerFunc
}

type Registry struct {
	commands map[string]*Command
}

func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]*Command)}
}

func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

func (r *Registry) Get(name string) (*Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *Registry) All() map[string]*Command {
	return r.commands
}

func (r *Registry) ByCategory() map[string][]*Command {
	seen := map[string]struct{}{}
	result := map[string][]*Command{}
	for _, cmd := range r.commands {
		if _, ok := seen[cmd.Name]; ok {
			continue
		}
		seen[cmd.Name] = struct{}{}
		result[cmd.Category] = append(result[cmd.Category], cmd)
	}
	for cat := range result {
		sort.Slice(result[cat], func(i, j int) bool {
			return result[cat][i].Name < result[cat][j].Name
		})
	}
	return result
}

func (r *Registry) Categories() []string {
	cats := map[string]struct{}{}
	for _, cmd := range r.commands {
		cats[cmd.Category] = struct{}{}
	}
	result := make([]string, 0, len(cats))
	for c := range cats {
		result = append(result, c)
	}
	sort.Strings(result)
	return result
}

var globalRegistry = NewRegistry()

func GlobalRegistry() *Registry {
	return globalRegistry
}

func Use(cmd *Command) {
	globalRegistry.Register(cmd)
}

func (c *Command) Execute(ptz *Ptz) error {
	if c.GroupOnly && !ptz.IsGroup {
		return nil
	}

	if c.OwnerOnly && !ptz.IsOwner() {
		return nil
	}

	if c.AdminOnly || c.BotAdmin {
		if err := ptz.LoadGroupInfo(); err != nil {
			ptz.Bot.Log.Warnf("LoadGroupInfo error for %s: %v", ptz.Chat, err)
			return nil
		}
	}

	if c.AdminOnly && !ptz.IsAdmin() && !ptz.IsOwner() {
		return nil
	}

	if c.BotAdmin && !ptz.IsBotAdmin() {
		ptz.ReplyText("❌ Bot harus jadi admin dulu.")
		return nil
	}

	if c.Quota != nil && c.Quota.Enabled && c.Quota.Cost > 0 && ptz.Bot != nil && ptz.Bot.Users != nil && !ptz.IsOwner() {
		if err := ConsumeQuota(ptz, c.Quota.Cost); err != nil {
			ptz.Bot.Log.Errorf("consume limit failed on %s: %v", c.Name, err)
			return nil
		}
	}

	if c.Limit != nil && c.Limit.Enabled && ptz.Bot != nil && ptz.Bot.CommandLimiter != nil {
		allowed, retryAfter := ptz.Bot.CommandLimiter.Allow(c.Name, ptz.Sender.User, c.Limit.Max, c.Limit.Window)
		if !allowed {
			ptz.ReplyText(fmt.Sprintf("Limit command %s habis coba lagi dalam %s", c.Name, formatRetryAfter(retryAfter)))
			return nil
		}
	}

	if ptz.Bot != nil && ptz.Bot.Users != nil {
		userID := ptz.GetPhoneJID().User
		if _, _, err := ptz.Bot.Users.TrackInteraction(userID); err != nil {
			ptz.Bot.Log.Errorf("track interaction failed for %s: %v", userID, err)
		}
	}

	return c.Handler(ptz)
}

func formatRetryAfter(d time.Duration) string {
	if d <= 0 {
		return "beberapa detik"
	}

	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	if mins > 0 {
		return fmt.Sprintf("%d menit %d detik", mins, secs)
	}
	return fmt.Sprintf("%d detik", secs)
}
