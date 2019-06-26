package api

import (
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
)

type ISlackBot interface {
	GetUsers() ([]string, error)
	SendTextMessage(channelIds []string, text string) error
	SendPollMessage(channelIds []string, block string) error
	Callback(w http.ResponseWriter, r *http.Request)
	Send(w http.ResponseWriter, r *http.Request)
}

type bot struct {
	Client     *slack.Client
	CallbackID string
	Question   string
}

const actionName = "answer"

func NewSlackBot(token string, callbackId string, question string) ISlackBot {
	client := slack.New(token)
	return &bot{Client: client, CallbackID: callbackId, Question: question}
}

func (c *bot) GetUsers() ([]string, error) {
	data := []string{}
	users, err := c.Client.GetUsers()
	if err != nil {
		return data, err
	}

	for _, user := range users {
		if user.Deleted || user.IsBot {
			continue
		}

		data = append(data, user.ID)
		fmt.Println(user.ID, user.RealName)
	}

	return data, nil
}

func (c *bot) SendTextMessage(channelIds []string, text string) error {
	for _, id := range channelIds {
		_, _, err := c.Client.PostMessage(id, slack.MsgOptionText(text, false))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *bot) SendPollMessage(channelIds []string, text string) error {
	message := slack.Message{}
	message.Attachments = make([]slack.Attachment, 1)
	message.Attachments[0].Text = text
	message.Attachments[0].CallbackID = c.CallbackID
	message.Attachments[0].Actions = []slack.AttachmentAction{
		{
			Name:  actionName,
			Text:  ":rage:",
			Type:  "button",
			Value: "1",
		},
		{
			Name:  actionName,
			Text:  ":white_frowning_face:",
			Type:  "button",
			Value: "2",
		},
		{
			Name:  actionName,
			Text:  ":neutral_face:",
			Type:  "button",
			Value: "3",
		},
		{
			Name:  actionName,
			Text:  ":blush:",
			Type:  "button",
			Value: "4",
		},
		{
			Name:  actionName,
			Text:  ":sunglasses:",
			Type:  "button",
			Value: "5",
		},
	}

	for _, id := range channelIds {
		slack.NewMessageItem(id, &message)
		_, _, err := c.Client.PostMessage(id, slack.MsgOptionAttachments(message.Attachments...))
		if err != nil {
			return err
		}
	}

	return nil
}
