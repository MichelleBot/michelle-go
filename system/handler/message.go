package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

var antiLinkRegex = regexp.MustCompile(`(?i)\b(?:https?://)?(?:chat\.whatsapp\.com/[a-zA-Z0-9]+|wa\.me/[0-9*~_]+|whatsapp\.com/channel/[a-zA-Z0-9]+)`)

func (h *EventHandler) handleAntiToxic(ptz *core.Ptz) {
	if ptz.IsFromMe || !ptz.IsGroup {
		return
	}

	var groupSettingsJSON string
	err := h.bot.DB.Conn.QueryRow("SELECT filter, toxic, member_data FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(nil, nil, &groupSettingsJSON) // Simplify to just getting what we need. Wait, need filter and toxic.
	// Actually query all for now
	var filterEnabled bool
	var toxicJSON string
	var memberDataJSON string
	err = h.bot.DB.Conn.QueryRow("SELECT filter, toxic, member_data FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&filterEnabled, &toxicJSON, &memberDataJSON)
	if err != nil || !filterEnabled {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	toxicWords := []string{}
	json.Unmarshal([]byte(toxicJSON), &toxicWords)

	text := core.ExtractBody(ptz.Message)
	if !utils.DetectBadword(text, toxicWords) {
		return
	}

	type MemberInfo struct {
		Warning int `json:"warning"`
	}
	memberData := make(map[string]*MemberInfo)
	if memberDataJSON != "" {
		json.Unmarshal([]byte(memberDataJSON), &memberData)
	}

	senderID := ptz.Sender.User
	if _, ok := memberData[senderID]; !ok {
		memberData[senderID] = &MemberInfo{}
	}
	memberData[senderID].Warning += 1
	warning := memberData[senderID].Warning

	if warning >= 5 {
		ptz.ReplyText("🚩 Warning : [ 5 / 5 ]")
		h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
		delete(memberData, senderID)
	} else {
		p := fmt.Sprintf("乂  *W A R N I N G* \n\nKamu mendapat +1 poin peringatan : [ %d / 5 ]\n\n> Jika kamu mendapatkan 5 poin peringatan, Kamu akan dikeluarkan dari grup ini.", warning)
		ptz.ReplyText(p)
	}

	updatedJSON, _ := json.Marshal(memberData)
	h.bot.DB.Conn.Exec("UPDATE groups SET member_data = ? WHERE jid = ?", string(updatedJSON), ptz.Chat.String())

	h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
}

func (h *EventHandler) handleAntiBot(ptz *core.Ptz) {
	if ptz.IsFromMe || !ptz.IsGroup {
		return
	}

	var antibot bool
	err := h.bot.DB.Conn.QueryRow("SELECT antibot FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antibot)
	if err != nil || !antibot {
		return
	}

	if ptz.IsOwner() {
		return
	}

	if utils.IsBot(ptz.Info.ID) && ptz.IsBotAdmin() {
		ptz.ReplyText("⚠ Bot lain tidak diperbolehkan di grup ini.")
		
		time.Sleep(1200 * time.Millisecond)
		
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
		h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
	}
}

func (h *EventHandler) handleAntiLink(ptz *core.Ptz) {
	if !ptz.IsGroup {
		return
	}

	var antilink bool
	err := h.bot.DB.Conn.QueryRow("SELECT antilink FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antilink)
	if err != nil {
		return
	}

	if err := ptz.LoadGroupInfo(); err != nil {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	text := core.ExtractBody(ptz.Message)
	match := antiLinkRegex.FindAllString(text, -1)
	if len(match) == 0 {
		return
	}

	for range match {
		// Delete
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))

		// Kick if antilink enabled
		if antilink {
			h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
		}
	}
}

func (h *EventHandler) handleAntiForward(ptz *core.Ptz) {
	if !ptz.IsGroup {
		return
	}

	var antiforward bool
	err := h.bot.DB.Conn.QueryRow("SELECT antiforward FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antiforward)
	if err != nil || !antiforward {
		return
	}

	if err := ptz.LoadGroupInfo(); err != nil {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	// Check if message is forwarded
	isForwarded := false
	if ptz.Message.ExtendedTextMessage != nil && ptz.Message.ExtendedTextMessage.ContextInfo != nil && ptz.Message.ExtendedTextMessage.ContextInfo.GetIsForwarded() {
		isForwarded = true
	} else if ptz.Message.ImageMessage != nil && ptz.Message.ImageMessage.ContextInfo != nil && ptz.Message.ImageMessage.ContextInfo.GetIsForwarded() {
		isForwarded = true
	} else if ptz.Message.VideoMessage != nil && ptz.Message.VideoMessage.ContextInfo != nil && ptz.Message.VideoMessage.ContextInfo.GetIsForwarded() {
		isForwarded = true
	}

	if isForwarded {
		// Delete
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
		// Kick
		h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
	}
}

func (h *EventHandler) handleAntiSticker(ptz *core.Ptz) {
	if !ptz.IsGroup {
		return
	}

	var antisticker bool
	err := h.bot.DB.Conn.QueryRow("SELECT antisticker FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antisticker)
	if err != nil || !antisticker {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	if ptz.Message.StickerMessage != nil {
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
	}
}

func (h *EventHandler) handleAntiVirtex(ptz *core.Ptz) {
	if !ptz.IsGroup {
		return
	}

	var antivirtex bool
	err := h.bot.DB.Conn.QueryRow("SELECT antivirtex FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antivirtex)
	if err != nil || !antivirtex {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	text := core.ExtractBody(ptz.Message)
	virtexRegex := regexp.MustCompile(`(?i)(৭৭৭৭৭৭৭৭|๒๒๒๒๒๒๒๒|๑๑๑๑๑๑๑๑|ดุท้่เึางืผิดุท้่เึางื)`)

	if virtexRegex.MatchString(text) || len(text) > 10000 {
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
		h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
	}
}

func (h *EventHandler) handleAntiTagSW(ptz *core.Ptz) {
	if !ptz.IsGroup {
		return
	}

	var antitagsw bool
	err := h.bot.DB.Conn.QueryRow("SELECT antitagsw FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&antitagsw)
	if err != nil || !antitagsw {
		return
	}

	if ptz.IsAdmin() || ptz.IsOwner() {
		return
	}

	// This is a placeholder for detection, as there is no direct GroupStatus in Go
	if ptz.Message.ProtocolMessage == nil {
		return
	}

	// Fetch member_data
	var memberDataJSON string
	err = h.bot.DB.Conn.QueryRow("SELECT member_data FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&memberDataJSON)
	if err != nil {
		return
	}

	type MemberInfo struct {
		Warning int `json:"warning"`
	}
	memberData := make(map[string]*MemberInfo)
	if memberDataJSON != "" {
		json.Unmarshal([]byte(memberDataJSON), &memberData)
	}

	senderID := ptz.Sender.User
	if _, ok := memberData[senderID]; !ok {
		memberData[senderID] = &MemberInfo{}
	}
	memberData[senderID].Warning += 1
	warning := memberData[senderID].Warning

	updatedJSON, _ := json.Marshal(memberData)
	h.bot.DB.Conn.Exec("UPDATE groups SET member_data = ? WHERE jid = ?", string(updatedJSON), ptz.Chat.String())

	if warning >= 3 {
		ptz.ReplyText(fmt.Sprintf("🚩 Warning : [ 3 / 3 ]"))
		h.bot.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{ptz.Sender}, whatsmeow.ParticipantChangeRemove)
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
	} else {
		p := fmt.Sprintf("乂  *W A R N I N G* \n\nKamu mendapat +1 poin peringatan : [ %d / 3 ]\n\n> Jika kamu mendapatkan 3 poin peringatan, Kamu akan dikeluarkan dari grup ini.", warning)
		ptz.ReplyText(p)
		h.bot.Client.SendMessage(context.Background(), ptz.Chat, h.bot.Client.BuildRevoke(ptz.Chat, ptz.Sender, ptz.Info.ID))
	}
}

func (h *EventHandler) handleMessageEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Message == nil || msg.Event == nil {
		return
	}

	if msg.IsFromMe {
		return
	}

    // Auto read
    h.bot.Client.MarkRead(context.Background(), []types.MessageID{msg.Info.ID}, msg.Info.Timestamp, msg.Chat, msg.Sender)

	h.logNormalizedMessage(msg)
	h.trackMessage(msg)

	ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)
	
	h.handleAntiToxic(ptz)
	h.handleAntiBot(ptz)
	h.handleAntiLink(ptz)
	h.handleAntiForward(ptz)
	h.handleAntiSticker(ptz)
	h.handleAntiVirtex(ptz)
	h.handleAntiTagSW(ptz)
	
	if !h.shouldProcessCommand(ptz) {
		return
	}

	h.executeCommand(ptz)
}

func (h *EventHandler) trackMessage(msg *core.NormalizedMessage) {
	if msg == nil || h.bot == nil || h.bot.Messages == nil {
		return
	}

	stored := core.NewStoredMessage(msg)
	if stored == nil {
		return
	}
	h.enrichStoredMessage(stored, msg)

	switch msg.Kind {
	case core.MessageText, core.MessageImage, core.MessageVideo, core.MessageDocument, core.MessageAudio, core.MessageVoice, core.MessageSticker, core.MessageContact, core.MessageContacts, core.MessageLocation, core.MessageLiveLocation, core.MessageButtons, core.MessageButtonReply, core.MessageList, core.MessageListReply, core.MessageTemplate, core.MessageTemplateReply, core.MessageInteractive, core.MessageInteractiveReply, core.MessageGroupInvite, core.MessageProduct, core.MessageOrder, core.MessagePayment, core.MessageRequestPhone, core.MessageKeepInChat, core.MessageStructured:
		h.bot.Messages.SaveStored(stored)
	}
}

func (h *EventHandler) enrichStoredMessage(stored *core.StoredMessage, msg *core.NormalizedMessage) {
	if stored == nil || msg == nil || msg.Message == nil {
		return
	}

	switch msg.Kind {
	case core.MessageImage, core.MessageVideo, core.MessageAudio, core.MessageVoice, core.MessageDocument, core.MessageSticker:
		stored.MIME = serialize.GetMediaMIME(msg.Message)
		stored.Filename = serialize.GetMediaFilename(msg.Message)
		if stored.Caption == "" {
			stored.Caption = serialize.GetMediaCaption(msg.Message)
		}
		data, err := serialize.DownloadMedia(h.bot.Client, msg.Message)
		if err != nil {
			h.bot.Log.Warnf("failed to cache media for anti-delete %s: %v", msg.Info.ID, err)
			return
		}
		stored.MediaData = data
	}
}

func (h *EventHandler) shouldProcessCommand(ptz *core.Ptz) bool {
	if h.bot.BotConfig.GetSelfMode() && !ptz.IsOwner() {
		return false
	}

	if h.bot.BotConfig.GetPrivateOnly() && ptz.IsGroup {
		return false
	}

	if h.bot.BotConfig.GetGroupOnly() && !ptz.IsGroup {
		return false
	}

	if h.shouldApplyAntispam(ptz) && !h.bot.Antispam.Check(ptz.Sender.User) {
		h.bot.Log.Warnf("anti-spam blocked sender %s in chat %s", ptz.Sender.User, ptz.Chat.String())
		return false
	}

	if ptz.Command == "" {
		return false
	}

	return true
}

func (h *EventHandler) shouldApplyAntispam(ptz *core.Ptz) bool {
	if !ptz.IsGroup {
		return true
	}

	settings := h.bot.Settings.GetGroupSettings(ptz.Chat)
	return settings.AntispamEnabled
}

func (h *EventHandler) executeCommand(ptz *core.Ptz) {
	cmd, ok := h.registry.Get(ptz.Command)
	if !ok {
		return
	}

	started := time.Now()
	defer func() {
		if r := recover(); r != nil {
			h.bot.Log.Errorf("Recovered from panic in command %s: %v", cmd.Usage[0], r)
			if err := ptz.ReplyText("❌ Terjadi error internal saat menjalankan perintah."); err != nil {
				h.bot.Log.Errorf("panic fallback reply failed on %s: %v", cmd.Usage[0], err)
			}
		}
		h.bot.Log.Debugf("command %s completed in %s", cmd.Usage[0], time.Since(started))
	}()

	if err := cmd.Execute(ptz); err != nil {
		h.bot.Log.Errorf("Command %s error: %v", cmd.Usage[0], err)
		if replyErr := ptz.ReplyText("❌ Perintah gagal dijalankan. Coba lagi sebentar."); replyErr != nil {
			h.bot.Log.Errorf("command error reply failed on %s: %v", cmd.Usage[0], replyErr)
		}
	}
}
