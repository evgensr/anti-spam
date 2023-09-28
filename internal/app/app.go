package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/evgensr/anti-spam/internal/config"
	"github.com/evgensr/anti-spam/internal/logging"
	"github.com/evgensr/anti-spam/internal/model"
	"github.com/evgensr/anti-spam/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	l      *logging.Logger
	t      *telegram.Telegram
	filter *model.Filter
}

func NewApp(c *config.Config, l *logging.Logger, t *telegram.Telegram) *App {
	return &App{
		cfg:    c,
		l:      l,
		t:      t,
		filter: &model.Filter{},
	}
}

func (a *App) Run() {

	go a.refresh()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.t.Bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message != nil {

			if update.Message.Chat.ID == a.cfg.Telegram.GroupID {

				if update.Message.NewChatMembers != nil && len(update.Message.NewChatMembers) > 0 && a.cfg.App.DeleteNotificationNewMember {
					req := tgbotapi.NewDeleteMessage(a.cfg.Telegram.GroupID, update.Message.MessageID)
					_, err := a.t.Bot.Request(req)
					if err != nil {
						a.l.Error("error delete message", zap.Error(err))
					}
				}

				// search word spam
				for _, word := range a.filter.Word {
					lowerCaseWord := strings.ToLower(word)
					lowerCaseWordMessage := strings.ToLower(update.Message.Text)

					if strings.Contains(lowerCaseWordMessage, lowerCaseWord) {

						if err := a.t.BanMember(a.cfg.Telegram.GroupID, update.Message.From.ID, -1); err != nil {
							a.l.Error("error ban member", zap.Error(err))
						}

						if err := a.t.DeleteMessage(a.cfg.Telegram.GroupID, update.Message.MessageID); err != nil {
							a.l.Error("error delete message", zap.Error(err))
						}

						badWord := fmt.Sprintf("Found word <b>%s</b> in the message:\n", word)

						msg := fmt.Sprintf("Id: %d\nUserName: @%s\nLastName: %s\nFirstName: %s\nBot: %t\nID: %d\nLanguageCode: %s\nChat: %d\nBad: %sMsg: %s",
							update.UpdateID,
							update.Message.From.UserName,
							update.Message.From.LastName,
							update.Message.From.FirstName,
							update.Message.From.IsBot,
							update.Message.From.ID,
							update.Message.From.LanguageCode,
							update.Message.Chat.ID,
							badWord,
							update.Message.Text,
						)

						for _, admin := range a.cfg.Telegram.AdminID {
							if err := a.t.SendNotificationAdmin(
								admin,
								msg,
								update.Message.From.ID,
								a.cfg.App.EnableSoundNotification,
							); err != nil {
								a.l.Error("error send notification admin", zap.Error(err))
							}
						}
						break
					}
				}
				continue
			}

			if !a.isChatIDInAdminID(update.Message.Chat.ID, a.cfg.Telegram.AdminID) {
				log.Println("not admin")
				continue
			}

			if update.Message.Text == "/start" {
				if err := a.t.SendMs(update.Message.From.ID, "Hello, admin"); err != nil {
					a.l.Error("error send message", zap.Error(err))
				}
				if err := a.t.SendStart(update.Message.Chat.ID,
					"Hello, admin", update.Message.From.ID,
					a.cfg.App.EnableSoundNotification,
				); err != nil {
					a.l.Error("error send message", zap.Error(err))
				}
			}

		}

		if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			if _, err := a.t.Bot.Request(callback); err != nil {
				a.l.Info("error respond to the callback query", zap.Error(err))
			}

			if strings.HasPrefix(callback.Text, "unBan/") {
				userId, err := strconv.Atoi(strings.Split(callback.Text, "/")[1])
				if err != nil {
					a.l.Error("error unban member", zap.Error(err))
				}
				if err = a.t.UnBanMember(a.cfg.Telegram.GroupID, int64(userId)); err != nil {
					a.l.Error("error unban member", zap.Error(err))
				}
			}

			if strings.HasPrefix(callback.Text, "updateWord") {
				err := a.loadStopWords()
				text := fmt.Sprintf("Update word:\n%s", a.format(a.filter.Word))
				if err != nil {
					a.l.Error("error update word", zap.Error(err))
					text = fmt.Sprintf("Update word err: %v", err)
				}
				if err = a.t.SendMs(update.CallbackQuery.From.ID, text); err != nil {
					a.l.Error("error send message", zap.Error(err))
				}
			}

			if strings.HasPrefix(callback.Text, "showURL") {
				text := fmt.Sprintf(
					"URL: %s\nCSV: %s",
					a.cfg.App.GoogleDoc,
					a.cfg.App.GoogleCSV,
				)
				if err := a.t.SendMs(update.CallbackQuery.From.ID, text); err != nil {
					a.l.Error("error send message", zap.Error(err))
				}
			}

			if strings.HasPrefix(callback.Text, "showWord") {
				text := fmt.Sprintf(
					"Word:\n%s",
					a.format(a.filter.Word),
				)
				if err := a.t.SendMs(update.CallbackQuery.From.ID, text); err != nil {
					a.l.Error("error send message", zap.Error(err))
				}
			}

		}
	}
}

func (a *App) isChatIDInAdminID(chatID int64, adminID []int64) bool {
	for _, id := range adminID {
		if id == chatID {
			return true
		}
	}
	return false
}

func (a *App) format(s []string) string {
	var builder strings.Builder
	for _, v := range s {
		builder.WriteString(v + "\n")
	}
	return builder.String()
}
