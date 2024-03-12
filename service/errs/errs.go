package errs

import "errors"

var (
	ErrNoRows        = errors.New("no rows")
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user not exists")

	ErrNoLastNotSolvedCase = errors.New("no last not solved case")
	ErrNoNotSolvedCase     = errors.New("no not solved case")

	ErrNoCase = errors.New("no case")
)
