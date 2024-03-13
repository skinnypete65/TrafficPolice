package errs

import "errors"

var (
	ErrNoRows        = errors.New("no rows")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidPass   = errors.New("invalid passwrod")
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user not exists")

	ErrNoLastNotSolvedCase = errors.New("no last not solved case")
	ErrNoNotSolvedCase     = errors.New("no not solved case")

	ErrNoCase = errors.New("no case")
)
