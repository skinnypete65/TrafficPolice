package errs

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrUserNotExists = errors.New("user not exists")

var ErrNoLastNotSolvedCase = errors.New("no last not solved case")

var ErrNoCase = errors.New("no case")
