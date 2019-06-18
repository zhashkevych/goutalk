package chat

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrMissingAccessRights = errors.New("user isn't owner of the room")
	ErrWrongPassword       = errors.New("wrong password")
)

type ErrorNotFound struct {
	entity string
	field  string
	value  string
}

func NewErrorNotFound(entity, field, value string) *ErrorNotFound {
	return &ErrorNotFound{entity: entity, field: field, value: value}
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("%s with %s %s not found", e.entity, e.field, e.value)
}
