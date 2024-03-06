package errs

import "errors"

var ErrUserExists = errors.New("user already exists")

var ErrNoLastNotSolvedCase = errors.New("no last not solved case")
