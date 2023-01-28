package command

import (
	"github.com/disgoorg/disgo/events"
)

type Command interface {
	GetName() string
	GetDescription() string
	Execute(e *events.MessageCreate, opts ...string) error
}
