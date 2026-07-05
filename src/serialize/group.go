package serialize

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func matchParticipant(p types.GroupParticipant, jid types.JID) bool {
	if jid.Server == types.HiddenUserServer {
		return p.LID.User == jid.User
	}
	return p.PhoneNumber.User == jid.User || p.JID.User == jid.User
}

func GetGroupAdmins(info *types.GroupInfo) []types.JID {
	admins := []types.JID{}
	for _, p := range info.Participants {
		if p.IsAdmin || p.IsSuperAdmin {
			admins = append(admins, p.JID)
		}
	}
	return admins
}

func IsParticipantAdmin(info *types.GroupInfo, jid types.JID) bool {
	for _, p := range info.Participants {
		if matchParticipant(p, jid) {
			return p.IsAdmin || p.IsSuperAdmin
		}
	}
	return false
}

func IsParticipantSuperAdmin(info *types.GroupInfo, jid types.JID) bool {
	for _, p := range info.Participants {
		if matchParticipant(p, jid) {
			return p.IsSuperAdmin
		}
	}
	return false
}

func IsParticipantInGroup(info *types.GroupInfo, jid types.JID) bool {
	for _, p := range info.Participants {
		if matchParticipant(p, jid) {
			return true
		}
	}
	return false
}

func KickParticipant(client *whatsmeow.Client, group types.JID, targets []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupParticipants(context.Background(), group, targets, whatsmeow.ParticipantChangeRemove)
}

func AddParticipant(client *whatsmeow.Client, group types.JID, targets []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupParticipants(context.Background(), group, targets, whatsmeow.ParticipantChangeAdd)
}

func PromoteParticipant(client *whatsmeow.Client, group types.JID, targets []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupParticipants(context.Background(), group, targets, whatsmeow.ParticipantChangePromote)
}

func DemoteParticipant(client *whatsmeow.Client, group types.JID, targets []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupParticipants(context.Background(), group, targets, whatsmeow.ParticipantChangeDemote)
}

func SetGroupEphemeral(client *whatsmeow.Client, group types.JID, duration time.Duration) error {
	return client.SetDisappearingTimer(context.Background(), group, duration, time.Now())
}

func SetGroupAnnounce(client *whatsmeow.Client, group types.JID, announce bool) error {
	return client.SetGroupAnnounce(context.Background(), group, announce)
}

func SetGroupLocked(client *whatsmeow.Client, group types.JID, locked bool) error {
	return client.SetGroupLocked(context.Background(), group, locked)
}

func SetGroupName(client *whatsmeow.Client, group types.JID, name string) error {
	return client.SetGroupName(context.Background(), group, name)
}

func SetGroupTopic(client *whatsmeow.Client, group types.JID, topic string) error {
	return client.SetGroupTopic(context.Background(), group, "", "", topic)
}

func GetInviteLink(client *whatsmeow.Client, group types.JID, reset bool) (string, error) {
	return client.GetGroupInviteLink(context.Background(), group, reset)
}

func JoinGroupWithLink(client *whatsmeow.Client, code string) (types.JID, error) {
	return client.JoinGroupWithLink(context.Background(), code)
}

func GetGroupInfoFromLink(client *whatsmeow.Client, code string) (*types.GroupInfo, error) {
	return client.GetGroupInfoFromLink(context.Background(), code)
}

func LeaveGroup(client *whatsmeow.Client, group types.JID) error {
	return client.LeaveGroup(context.Background(), group)
}

func GetGroupInfo(client *whatsmeow.Client, group types.JID) (*types.GroupInfo, error) {
	return client.GetGroupInfo(context.Background(), group)
}

func GetJoinedGroups(client *whatsmeow.Client) ([]*types.GroupInfo, error) {
	return client.GetJoinedGroups(context.Background())
}

func CreateGroup(client *whatsmeow.Client, name string, participants []types.JID) (*types.GroupInfo, error) {
	return client.CreateGroup(context.Background(), whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
	})
}

func GetGroupParticipants(client *whatsmeow.Client, group types.JID) ([]types.GroupParticipant, error) {
	info, err := client.GetGroupInfo(context.Background(), group)
	if err != nil {
		return nil, err
	}
	return info.Participants, nil
}

func GetGroupMembers(client *whatsmeow.Client, group types.JID) ([]types.GroupParticipant, error) {
	info, err := client.GetGroupInfo(context.Background(), group)
	if err != nil {
		return nil, err
	}
	var members []types.GroupParticipant
	for _, p := range info.Participants {
		if !p.IsAdmin && !p.IsSuperAdmin {
			members = append(members, p)
		}
	}
	return members, nil
}

func GetGroupOwner(info *types.GroupInfo) types.JID {
	return info.OwnerJID
}

func IsGroupOwner(info *types.GroupInfo, jid types.JID) bool {
	return info.OwnerJID.User == jid.User || info.OwnerPN.User == jid.User
}

func GetGroupMemberCount(info *types.GroupInfo) int {
	return len(info.Participants)
}

func GetGroupAdminCount(info *types.GroupInfo) int {
	count := 0
	for _, p := range info.Participants {
		if p.IsAdmin || p.IsSuperAdmin {
			count++
		}
	}
	return count
}

func GetGroupParticipantByPhone(info *types.GroupInfo, phone string) *types.GroupParticipant {
	phone = NormalizePhone(phone)
	for i, p := range info.Participants {
		if p.PhoneNumber.User == phone || p.JID.User == phone {
			return &info.Participants[i]
		}
	}
	return nil
}

func GetGroupAllJIDs(info *types.GroupInfo) []types.JID {
	jids := make([]types.JID, 0, len(info.Participants))
	for _, p := range info.Participants {
		jids = append(jids, p.JID)
	}
	return jids
}

func GetGroupMemberJIDs(info *types.GroupInfo) []types.JID {
	var jids []types.JID
	for _, p := range info.Participants {
		if !p.IsAdmin && !p.IsSuperAdmin {
			jids = append(jids, p.JID)
		}
	}
	return jids
}

func GetGroupAdminJIDs(info *types.GroupInfo) []types.JID {
	var jids []types.JID
	for _, p := range info.Participants {
		if p.IsAdmin || p.IsSuperAdmin {
			jids = append(jids, p.JID)
		}
	}
	return jids
}

func IsGroupEphemeral(info *types.GroupInfo) bool {
	return info.IsEphemeral
}

func IsGroupLocked(info *types.GroupInfo) bool {
	return info.IsLocked
}

func IsGroupAnnounce(info *types.GroupInfo) bool {
	return info.IsAnnounce
}

func IsGroupIncognito(info *types.GroupInfo) bool {
	return info.IsIncognito
}

func IsGroupCommunity(info *types.GroupInfo) bool {
	return info.IsParent
}

func SetGroupDescription(client *whatsmeow.Client, group types.JID, description string) error {
	return client.SetGroupDescription(context.Background(), group, description)
}

func SetGroupJoinApprovalMode(client *whatsmeow.Client, group types.JID, requireApproval bool) error {
	return client.SetGroupJoinApprovalMode(context.Background(), group, requireApproval)
}

func SetGroupMemberAddMode(client *whatsmeow.Client, group types.JID, adminOnly bool) error {
	mode := types.GroupMemberAddModeAllMember
	if adminOnly {
		mode = types.GroupMemberAddModeAdmin
	}
	return client.SetGroupMemberAddMode(context.Background(), group, mode)
}

func SetGroupPhoto(client *whatsmeow.Client, group types.JID, jpegData []byte) (string, error) {
	return client.SetGroupPhoto(context.Background(), group, jpegData)
}

func RemoveGroupPhoto(client *whatsmeow.Client, group types.JID) (string, error) {
	return client.SetGroupPhoto(context.Background(), group, nil)
}

func LinkGroup(client *whatsmeow.Client, parent, child types.JID) error {
	return client.LinkGroup(context.Background(), parent, child)
}

func UnlinkGroup(client *whatsmeow.Client, parent, child types.JID) error {
	return client.UnlinkGroup(context.Background(), parent, child)
}

func GetGroupRequestParticipants(client *whatsmeow.Client, group types.JID) ([]types.GroupParticipantRequest, error) {
	return client.GetGroupRequestParticipants(context.Background(), group)
}

func ApproveJoinRequests(client *whatsmeow.Client, group types.JID, jids []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupRequestParticipants(context.Background(), group, jids, whatsmeow.ParticipantChangeApprove)
}

func RejectJoinRequests(client *whatsmeow.Client, group types.JID, jids []types.JID) ([]types.GroupParticipant, error) {
	return client.UpdateGroupRequestParticipants(context.Background(), group, jids, whatsmeow.ParticipantChangeReject)
}

func GetSubGroups(client *whatsmeow.Client, community types.JID) ([]*types.GroupLinkTarget, error) {
	return client.GetSubGroups(context.Background(), community)
}

func GetLinkedGroupsParticipants(client *whatsmeow.Client, community types.JID) ([]types.JID, error) {
	return client.GetLinkedGroupsParticipants(context.Background(), community)
}

func CreateCommunity(client *whatsmeow.Client, name string, participants []types.JID) (*types.GroupInfo, error) {
	return client.CreateGroup(context.Background(), whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
		GroupParent:  types.GroupParent{IsParent: true},
	})
}

func CreateGroupInCommunity(client *whatsmeow.Client, name string, communityJID types.JID, participants []types.JID) (*types.GroupInfo, error) {
	return client.CreateGroup(context.Background(), whatsmeow.ReqCreateGroup{
		Name:              name,
		Participants:      participants,
		GroupLinkedParent: types.GroupLinkedParent{LinkedParentJID: communityJID},
	})
}
