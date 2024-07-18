package types

import (
	"errors"
	"fmt"
)

var ErrNoMatchingParser = errors.New("message did not match any known pattern for parsing")

var ErrIncompleteMatch = errors.New("message did not conform to the expected format for parsing")

var ErrNoMatchingPlatform = errors.New("platform not supported")

func MissingFieldErr(field string) error {
	return fmt.Errorf("request is missing field '%s'", field)
}

func InvalidTypeErr(field string) error {
	return fmt.Errorf("type of request field '%s' is invalid", field)
}
