package command_tests

import (
	"github.com/disgoorg/snowflake/v2"
)

type fakeMessenger struct {
	previousMessage string
	sentMessage     string
}

func NewFakeMessenger() fakeMessenger {
	return fakeMessenger{
		previousMessage: "",
		sentMessage:     "",
	}
}

func (msgr fakeMessenger) SetChannel(channel snowflake.ID) {

}

func (msgr *fakeMessenger) SendSimpleMessage(msg string) {
	msgr.previousMessage = msgr.sentMessage
	msgr.sentMessage = msg
}
