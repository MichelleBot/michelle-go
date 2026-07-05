package handler

import (
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/commands/games"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func resolveSender(info types.MessageInfo) (phoneJID types.JID, displayName string) {
	sender := info.Sender
	alt := info.SenderAlt

	if sender.Server == types.HiddenUserServer {
		if !alt.IsEmpty() && alt.Server == types.DefaultUserServer {
			phoneJID = types.NewJID(alt.User, types.DefaultUserServer)
		} else {
			phoneJID = types.NewJID(sender.User, types.DefaultUserServer)
		}
	} else {
		phoneJID = types.NewJID(sender.User, types.DefaultUserServer)
	}

	if info.PushName != "" && info.PushName != "-" {
		displayName = info.PushName
	} else {
		displayName = phoneJID.User
	}
	return
}

func isCommand(text string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(text, p) {
			return true
		}
	}
	return false
}

func getReplyInfo(msg *waE2E.Message) (quotedID, userAnswer string) {
	if msg == nil {
		return "", ""
	}

	switch {
	case msg.ExtendedTextMessage != nil:
		ctx := msg.ExtendedTextMessage.GetContextInfo()
		if ctx.GetStanzaID() == "" {
			return "", ""
		}
		return ctx.GetStanzaID(), strings.TrimSpace(msg.ExtendedTextMessage.GetText())

	case msg.ImageMessage != nil:
		ctx := msg.ImageMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.ImageMessage.GetCaption())

	case msg.VideoMessage != nil:
		ctx := msg.VideoMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.VideoMessage.GetCaption())

	case msg.AudioMessage != nil:
		ctx := msg.AudioMessage.GetContextInfo()
		return ctx.GetStanzaID(), ""

	case msg.DocumentMessage != nil:
		ctx := msg.DocumentMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.DocumentMessage.GetCaption())
	}

	return "", ""
}

func (h *EventHandler) handleMessageEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Message == nil || msg.Event == nil {
		return
	}

	if h.tryHandleGameAnswer(msg) {
		return
	}

	if msg.IsFromMe {
		return
	}

	h.logNormalizedMessage(msg)
	h.trackMessage(msg)

	ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)
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

func (h *EventHandler) tryHandleGameAnswer(msg *core.NormalizedMessage) bool {
	quotedID, userAnswer := getReplyInfo(msg.Message)
	if quotedID == "" || userAnswer == "" || isCommand(userAnswer, h.bot.Config.Prefixes) {
		return false
	}

	if res, sess, found := games.CheckAnswer(quotedID, userAnswer); found {
		phoneJID, _ := resolveSender(msg.Info)
		ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)

		switch res {
		case games.AnswerCorrect:
			rewardAmount := games.RewardForClueCount(sess.ClueCount)
			rewardText := ""
			if rewardAmount > 0 && h.bot != nil && h.bot.Users != nil {
				if err := h.bot.Users.AddCredit(phoneJID.User, rewardAmount); err != nil {
					h.bot.Log.Errorf("game reward add credit error: %v", err)
				} else {
					rewardText = fmt.Sprintf("\n\nReward game: +%d balance", rewardAmount)
				}
			}

			text := fmt.Sprintf(
				"✅ *Benar!* Selamat @%s!\n\n🎮 Jawaban: *%s*%s",
				phoneJID.User, sess.Answer,
				rewardText,
			)
			if err := serialize.SendTextMentionReplyToID(
				h.bot.Client,
				msg.Info.Chat,
				text,
				[]types.JID{phoneJID},
				sess.QuestionID,
				phoneJID.String(),
			); err != nil {
				h.bot.Log.Errorf("game answer mention reply error: %v", err)
			}
			games.DeleteSession(sess)
		case games.AnswerVeryClose:
			if err := ptz.ReplyText(fmt.Sprintf("Jawabanmu sudah sangat dekat. Coba periksa lagi ejaan atau kata akhirnya.\nKalau pakai clue sekarang, reward turun jadi +%d balance.\nKalau mentok, reply soal lalu ketik *%sclue*. Kalau mau berhenti, reply soal lalu ketik *%snyerah*.", games.RewardForClueCount(sess.ClueCount+1), h.bot.GetPrefix(), h.bot.GetPrefix())); err != nil {
				h.bot.Log.Errorf("game very-close reply error: %v", err)
			}
		case games.AnswerGettingClose:
			if err := ptz.ReplyText(fmt.Sprintf("Jawabanmu mulai mendekati. Coba pikirkan kata kunci yang paling cocok.\nKalau pakai clue sekarang, reward jadi +%d balance.\nKamu juga bisa reply soal lalu ketik *%sclue* untuk petunjuk.", games.RewardForClueCount(sess.ClueCount+1), h.bot.GetPrefix())); err != nil {
				h.bot.Log.Errorf("game getting-close reply error: %v", err)
			}
		case games.AnswerWrong:
			if err := ptz.ReplyText(fmt.Sprintf("Jawaban masih belum tepat. Coba jawab lagi dengan reply ke soal.\nKalau perlu bantuan, pakai *%sclue* dan reward akan mengikuti jumlah clue yang sudah dipakai.\nKalau mau berhenti, pakai *%snyerah* dengan reply soal dan reward jadi 0.", h.bot.GetPrefix(), h.bot.GetPrefix())); err != nil {
				h.bot.Log.Errorf("game wrong-answer reply error: %v", err)
			}
		}

		return true
	}

	ptzForGame := core.NewPtzFromNormalizedMessage(h.bot, msg)
	if handled, err := games.ProcessCcsdAnswer(ptzForGame, quotedID, userAnswer); handled {
		if err != nil {
			h.bot.Log.Errorf("ccsd answer error: %v", err)
		}
		return true
	}

	if handled, err := games.ProcessFamily100Answer(ptzForGame, quotedID, userAnswer); handled {
		if err != nil {
			h.bot.Log.Errorf("family100 answer error: %v", err)
		}
		return true
	}

	return false
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
			h.bot.Log.Errorf("Recovered from panic in command %s: %v", cmd.Name, r)
			if err := ptz.ReplyText("❌ Terjadi error internal saat menjalankan perintah."); err != nil {
				h.bot.Log.Errorf("panic fallback reply failed on %s: %v", cmd.Name, err)
			}
		}
		h.bot.Log.Debugf("command %s completed in %s", cmd.Name, time.Since(started))
	}()

	if err := cmd.Execute(ptz); err != nil {
		h.bot.Log.Errorf("Command %s error: %v", cmd.Name, err)
		if replyErr := ptz.ReplyText("❌ Perintah gagal dijalankan. Coba lagi sebentar."); replyErr != nil {
			h.bot.Log.Errorf("command error reply failed on %s: %v", cmd.Name, replyErr)
		}
	}
}
