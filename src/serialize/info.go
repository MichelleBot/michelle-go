package serialize

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func IsOnWhatsApp(client *whatsmeow.Client, phones []string) ([]types.IsOnWhatsAppResponse, error) {
	return client.IsOnWhatsApp(context.Background(), phones)
}

func CheckNumber(client *whatsmeow.Client, phone string) (bool, error) {
	results, err := client.IsOnWhatsApp(context.Background(), []string{"+" + phone})
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, nil
	}
	return results[0].IsIn, nil
}

func GetUserInfo(client *whatsmeow.Client, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return client.GetUserInfo(context.Background(), jids)
}

func GetBusinessProfile(client *whatsmeow.Client, jid types.JID) (*types.BusinessProfile, error) {
	return client.GetBusinessProfile(context.Background(), jid)
}

func GetProfilePicture(client *whatsmeow.Client, jid types.JID) (*types.ProfilePictureInfo, error) {
	return client.GetProfilePictureInfo(context.Background(), jid, nil)
}

func GetProfilePicturePreview(client *whatsmeow.Client, jid types.JID) (*types.ProfilePictureInfo, error) {
	return client.GetProfilePictureInfo(context.Background(), jid, &whatsmeow.GetProfilePictureParams{Preview: true})
}

func GetStatusMessage(client *whatsmeow.Client, jid types.JID) (string, error) {
	userInfo, err := client.GetUserInfo(context.Background(), []types.JID{jid})
	if err != nil {
		return "", err
	}
	if info, ok := userInfo[jid]; ok {
		return info.Status, nil
	}
	return "", nil
}

func ResolveBusinessLink(client *whatsmeow.Client, code string) (*types.BusinessMessageLinkTarget, error) {
	return client.ResolveBusinessMessageLink(context.Background(), code)
}

func ResolveContactQRLink(client *whatsmeow.Client, code string) (*types.ContactQRLinkTarget, error) {
	return client.ResolveContactQRLink(context.Background(), code)
}

func GetContactQRLink(client *whatsmeow.Client, revoke bool) (string, error) {
	return client.GetContactQRLink(context.Background(), revoke)
}

func SetStatusMessage(client *whatsmeow.Client, msg string) error {
	return client.SetStatusMessage(context.Background(), msg)
}

func GetBotList(client *whatsmeow.Client) ([]types.BotListInfo, error) {
	return client.GetBotListV2(context.Background())
}

func GetBotProfiles(client *whatsmeow.Client, botInfo []types.BotListInfo) ([]types.BotProfileInfo, error) {
	return client.GetBotProfiles(context.Background(), botInfo)
}

func GetUserDevices(client *whatsmeow.Client, jids []types.JID) ([]types.JID, error) {
	return client.GetUserDevices(context.Background(), jids)
}

func GetUserDevicesSingle(client *whatsmeow.Client, jid types.JID) ([]types.JID, error) {
	return client.GetUserDevices(context.Background(), []types.JID{jid})
}

func GetUserInfoBatch(client *whatsmeow.Client, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return client.GetUserInfo(context.Background(), jids)
}

func GetUserStatus(client *whatsmeow.Client, jid types.JID) (string, error) {
	info, err := client.GetUserInfo(context.Background(), []types.JID{jid})
	if err != nil {
		return "", err
	}
	if v, ok := info[jid]; ok {
		return v.Status, nil
	}
	return "", nil
}

func GetUserLID(client *whatsmeow.Client, jid types.JID) (types.JID, error) {
	info, err := client.GetUserInfo(context.Background(), []types.JID{jid})
	if err != nil {
		return types.EmptyJID, err
	}
	if v, ok := info[jid]; ok {
		return v.LID, nil
	}
	return types.EmptyJID, nil
}

func GetProfilePictureFull(client *whatsmeow.Client, jid types.JID) (*types.ProfilePictureInfo, error) {
	return client.GetProfilePictureInfo(context.Background(), jid, &whatsmeow.GetProfilePictureParams{
		Preview: false,
	})
}

func GetProfilePictureThumb(client *whatsmeow.Client, jid types.JID) (*types.ProfilePictureInfo, error) {
	return client.GetProfilePictureInfo(context.Background(), jid, &whatsmeow.GetProfilePictureParams{
		Preview: true,
	})
}

func SetProfilePhoto(client *whatsmeow.Client, jpegData []byte) (string, error) {
	jid := client.Store.GetJID().ToNonAD()
	return client.SetGroupPhoto(context.Background(), jid, jpegData)
}

func RemoveProfilePhoto(client *whatsmeow.Client) (string, error) {
	jid := client.Store.GetJID().ToNonAD()
	return client.SetGroupPhoto(context.Background(), jid, nil)
}

func GetOwnJID(client *whatsmeow.Client) types.JID {
	return client.Store.GetJID()
}

func GetOwnLID(client *whatsmeow.Client) types.JID {
	return client.Store.GetLID()
}

func IsOnWhatsAppBatch(client *whatsmeow.Client, phones []string) ([]types.IsOnWhatsAppResponse, error) {
	return client.IsOnWhatsApp(context.Background(), phones)
}

func CheckPhoneOnWhatsApp(client *whatsmeow.Client, phone string) (bool, types.JID, error) {
	results, err := client.IsOnWhatsApp(context.Background(), []string{"+" + phone})
	if err != nil {
		return false, types.EmptyJID, err
	}
	if len(results) == 0 {
		return false, types.EmptyJID, nil
	}
	return results[0].IsIn, results[0].JID, nil
}
