package errs

import "errors"

var (
	ErrNoRows        = errors.New("no rows")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidPass   = errors.New("invalid password")
	ErrEmptyPayload  = errors.New("empty payload")

	ErrUnknownCameraID       = errors.New("unknown camera id")
	ErrUnknownCameraType     = errors.New("unknown camera type")
	ErrInvalidCameraID       = errors.New("invalid camera id")
	ErrCameraNotExists       = errors.New("camera with id not exists")
	ErrInvalidViolationID    = errors.New("invalid violation id")
	ErrInvalidRelevantParams = errors.New("invalid relevant params")

	ErrUserNotExists = errors.New("user not exists")

	ErrNoLastNotSolvedCase = errors.New("no last not solved case")
	ErrNoNotSolvedCase     = errors.New("no not solved case")

	ErrNoCase      = errors.New("no case")
	ErrNoTransport = errors.New("no transport")
	ErrNoImage     = errors.New("no image")

	ErrExpertNotExists = errors.New("expert not exists")
)
