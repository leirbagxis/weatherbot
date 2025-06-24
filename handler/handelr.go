package handler

import (
	"fmt"
	"log"
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/leirbagxis/weatherbot/client/openweather"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owclient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owclient,
	}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.HandlerUpdate(update)
	}
}

func (h *Handler) HandlerUpdate(update tgbotapi.Update) {
	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		coordinates, err := h.owClient.Coordinates(update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao pegar as coordenadas")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao pegar a temperatura")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("A temperatura em %s é de %d °C.",
				update.Message.Text, int(math.Round(weather.Temp))),
		)
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)
	}
}
