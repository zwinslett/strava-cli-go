package telegram

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Client struct {
	httpClient *http.Client
	botToken   string
	chatID     string
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		botToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID:     os.Getenv("TELEGRAM_CHAT_ID"),
		baseURL:    "https://api.telegram.org/bot",
	}
}

func (c *Client) SendMessage(ctx context.Context, message string) error {
	messageURL := c.baseURL + c.botToken + "/sendMessage"
	data := url.Values{}
	data.Set("chat_id", c.chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		messageURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("failed %d body=%s", resp.StatusCode, body.String())
	}

	return nil
}
