package types

import "errors"

var ErrNoMatchingParser = errors.New("message did not match any known pattern for parsing")

var ErrIncompleteMatch = errors.New("message did not conform to the expected format for parsing")

var ErrNoMatchingPlatform = errors.New("platform not supported")
