package notifier

import (
	"fmt"
	"net/http"
	"net/url"
)

type Telegram struct {
	Token  string
	ChatId string
}

func NewTelegramNotifier(token string, chatId string) *Telegram {
	return &Telegram{
		Token:  token,
		ChatId: chatId,
	}
}

func (n *Telegram) Notify(message string) {
	_, _ = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?text=%s&chat_id=%s", n.Token, url.QueryEscape(message), n.ChatId))
}
