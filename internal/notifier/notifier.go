package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	kafkago "github.com/segmentio/kafka-go"
	"log/slog"
	"raspyxTelegramNotifier/internal/dto"
	"raspyxTelegramNotifier/internal/kafka"
	"raspyxTelegramNotifier/internal/usecase"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Timestamp time.Time   `json:"timestamp"`
	Message   interface{} `json:"message"`
}

type Notifier struct {
	log      *slog.Logger
	userUC   *usecase.UserUseCase
	consumer *kafka.Consumer
	bot      *tgbotapi.BotAPI
	token    string
	debug    bool
}

func NewNotifier(log *slog.Logger, userUC *usecase.UserUseCase, consumer *kafka.Consumer, token string, debug bool) *Notifier {
	return &Notifier{
		log:      log,
		userUC:   userUC,
		token:    token,
		debug:    debug,
		consumer: consumer,
	}
}

func (n *Notifier) Run(ctx context.Context) {
	bot, err := tgbotapi.NewBotAPI(n.token)
	if err != nil {
		n.log.Error("error creating bot API", slog.String("error", err.Error()))
		return
	}

	if bot == nil {
		n.log.Error("bot is nil")
		return
	}

	n.bot = bot
	n.bot.Debug = n.debug
	defer n.bot.StopReceivingUpdates()

	n.log.Info(fmt.Sprintf("authorized on account %s", n.bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := n.bot.GetUpdatesChan(u)

	msgChan := make(chan kafkago.Message, 100)

	go func() {
		defer close(msgChan)
		err = n.consumer.Consume(ctx, msgChan)
		if err != nil {
			n.log.Error(fmt.Sprintf("error kafka consumer: %v", err))
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			n.log.Info("context canceled, notifier Run() is exiting")
			return
		case update, ok := <-updates:
			if !ok {
				n.log.Info("updates channel closed, notifier Run() is exiting")
				return
			}

			// Message received
			if update.Message != nil {
				if update.Message.Command() == "start" {
					n.commandStart(ctx, &update)
				} else if update.Message.Command() == "delete" {
					n.commandDelete(ctx, &update)
				} else {
					n.respondUnknownMessage(&update)
				}
			}
		case msg, ok := <-msgChan:
			if !ok {
				n.log.Info("msgChan channel closed, notifier Run() is exiting")
				return
			}

			var data Event
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				fmt.Println(err)
			}

			n.log.Info("received kafka message", slog.Any("data", data))

			n.mailing(ctx, string(msg.Value))
		}
	}
}

func (n *Notifier) respondUnknownMessage(update *tgbotapi.Update) {
	n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "I don't answer messages, please don't text me")

	n.log.Info(
		"received message",
		slog.Any(
			"from",
			map[string]string{
				"username": update.Message.From.UserName,
				"tid":      strconv.FormatInt(update.Message.From.ID, 10),
				"is_bot":   strconv.FormatBool(update.Message.From.IsBot),
			},
		),
		slog.String("message", update.Message.Text),
	)
}

func (n *Notifier) replyToMessage(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID

	_, err := n.bot.Send(msg)
	if err != nil {
		n.log.Error("error sending message", slog.String("error", err.Error()))
	}
}

func (n *Notifier) commandStart(ctx context.Context, update *tgbotapi.Update) {
	err := n.userUC.Create(ctx, &dto.CreateUser{TelegramID: update.Message.From.ID})
	if err != nil {
		if strings.Contains(err.Error(), "exist") {
			n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "⚠️ You are already subscribed")
			return
		}
		n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "internal error")
		n.log.Error("error adding user to db", slog.String("error", err.Error()))
		return
	}

	n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "✅ You have been successfully subscribed")
	n.log.Info("user successfully added to db", slog.Any(
		"user",
		map[string]string{
			"username": update.Message.From.UserName,
			"tid":      strconv.FormatInt(update.Message.From.ID, 10),
			"is_bot":   strconv.FormatBool(update.Message.From.IsBot),
		},
	))
}

func (n *Notifier) commandDelete(ctx context.Context, update *tgbotapi.Update) {
	err := n.userUC.Delete(ctx, update.Message.From.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "⚠️ You have already unsubscribed")
			return
		}
		n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "internal error")
		n.log.Error("error deleting user from db", slog.String("error", err.Error()))
		return
	}

	n.replyToMessage(update.Message.Chat.ID, update.Message.MessageID, "✅ You have successfully unsubscribed")
	n.log.Info("user successfully deleted from db", slog.Any(
		"user",
		map[string]string{
			"username": update.Message.From.UserName,
			"tid":      strconv.FormatInt(update.Message.From.ID, 10),
			"is_bot":   strconv.FormatBool(update.Message.From.IsBot),
		},
	))
}

func prepareMessage(text string) (string, error) {
	var data Event
	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("error preparing message: %v", err))
	}

	return fmt.Sprintf("<blockquote><b>🕛 Timestamp:</b> %v\n"+
		"<b>🗨️ Message:</b> %v</blockquote>", data.Timestamp.Format("15:04:05.999 02.01.2006 MST"), data.Message), nil
}

func (n *Notifier) sendMessage(chatID int64, text string) {
	text, err := prepareMessage(text)
	if err != nil {
		n.log.Error("error", slog.String("error", err.Error()))
		return
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err = n.bot.Send(msg)
	if err != nil {
		n.log.Error("error sending message", slog.String("error", err.Error()))
		return
	}

	n.log.Info("sent message", slog.Any("data", map[string]string{
		"tid": strconv.FormatInt(chatID, 10),
		"msg": text,
	}))
}

func (n *Notifier) mailing(ctx context.Context, text string) {
	users, err := n.userUC.Get(ctx)
	if err != nil {
		n.log.Error("error sending messages", slog.String("error", err.Error()))
		return
	}

	for _, user := range users {
		select {
		case <-ctx.Done():
			return
		default:
			n.sendMessage(user.TelegramID, text)
		}
	}
}
