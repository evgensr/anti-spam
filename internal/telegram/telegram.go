package telegram

import (
	"fmt"
	"time"

	"github.com/evgensr/anti-spam/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot *tgbotapi.BotAPI
	cfg *config.Config
}

func NewTelegram(c *config.Config) (*Telegram, error) {

	bot, err := tgbotapi.NewBotAPI(c.Telegram.Token)
	bot.Debug = true
	if err != nil {
		return nil, err
	}

	return &Telegram{
		Bot: bot,
		cfg: c,
	}, nil
}

func (t *Telegram) BanMember(gid int64, uid int64, sec int64) error {
	if sec <= 0 {
		sec = 9999999999999
	}
	chatuserconfig := tgbotapi.ChatMemberConfig{ChatID: gid, UserID: uid}

	Permissions := tgbotapi.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
	}

	restricconfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        time.Now().Unix() + sec,
		Permissions:      &Permissions,
	}
	_, err := t.Bot.Request(restricconfig)
	return err
}

func (t *Telegram) UnBanMember(gid int64, uid int64) error {

	chatuserconfig := tgbotapi.ChatMemberConfig{ChatID: gid, UserID: uid}

	Permissions := tgbotapi.ChatPermissions{
		CanSendMessages:       true,
		CanSendMediaMessages:  true,
		CanSendOtherMessages:  true,
		CanAddWebPagePreviews: true,
	}

	restricconfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        9999999999999,
		Permissions:      &Permissions,
	}
	_, err := t.Bot.Request(restricconfig)

	return err

}

func (t *Telegram) DeleteMessage(chatID int64, messageID int) error {
	dMessage := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}
	_, err := t.Bot.Request(dMessage)

	return err
}

func (t *Telegram) SendNotificationAdmin(chatID int64, text string, userID int64, disableNotification bool) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = disableNotification

	button := tgbotapi.NewInlineKeyboardMarkup(

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Unban", "unBan/"+fmt.Sprint(userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Refresh Word List", "updateWord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Remind the link", "showURL"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Show the current list of words", "showWord"),
		),
	)

	msg.Text = text
	msg.ReplyMarkup = button

	_, err := t.Bot.Send(msg)
	return err
}

func (t *Telegram) SendStart(chatID int64, text string, userID int64, disableNotification bool) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = disableNotification

	button := tgbotapi.NewInlineKeyboardMarkup(

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Refresh Word List", "updateWord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Remind the link", "showURL"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Show the current list of words", "showWord"),
		),
	)

	msg.Text = text
	msg.ReplyMarkup = button

	_, err := t.Bot.Send(msg)
	return err
}

func (t *Telegram) Send(msgSend tgbotapi.MessageConfig) error {
	msgSend.ParseMode = "HTML"
	_, err := t.Bot.Send(msgSend)
	return err
}
func (t *Telegram) SendMs(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "HTML"
	_, err := t.Bot.Send(msg)
	return err
}
