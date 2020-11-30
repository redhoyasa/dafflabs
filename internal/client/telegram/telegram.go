package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Client struct {
	botApi *tgbotapi.BotAPI
}

func NewClient(token string) (*Client, error) {
	c := new(Client)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	c.botApi = bot
	return c, nil
}

func (c *Client) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(143635997, text)
	if _, err := c.botApi.Send(msg); err != nil {
		return err
	}

	return nil
}
