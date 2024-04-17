package slack

import (
	"strings"
	"errors"

	"github.com/slack-go/slack"
)

type Config struct {
	Token		string
	ChannelId	string
	IsUsed		bool
}

var config Config

func Init(cfg Config) {
	config = cfg
}

func GetConfig() Config {
	return config
}

// @param msg, token, channelId
func SendMessage(args ... string) error {
	var (
		sendMsg 		string
		sendToken		string = config.Token
		sendChannelId	string = config.ChannelId
	)

	if len(args) == 0 {
		return errors.New("not found params")
	}
	
	if strings.Compare(args[0], "") == 0 {
		sendMsg = "no message"
	} else {
		sendMsg = args[0]
	}

	switch len(args) {
	case 1:
		sendMsg = args[0]
	case 2:
		sendMsg = args[0]
		sendToken = args[1]
	case 3:
		sendMsg = args[0]
		sendToken = args[1]
		sendChannelId = args[2]
	}

	api := slack.New(sendToken)

	_, _, err := api.PostMessage(
		sendChannelId,
		slack.MsgOptionText(sendMsg, false),
	)

	if err != nil {
		return err
	}

	return nil
}