package command_tests

import (
	"github.com/disgoorg/snowflake/v2"
)

type fakeMessenger struct {
	sentMessage string
}

func NewFakeMessenger() fakeMessenger {
	return fakeMessenger{
		sentMessage: "",
	}
}

func (msgr fakeMessenger) SetChannel(channel snowflake.ID) {

}

func (msgr *fakeMessenger) SendSimpleMessage(msg string) {
	msgr.sentMessage = msg
}
