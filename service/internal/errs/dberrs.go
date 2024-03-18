package errs

import "github.com/lib/pq"

const ForeignKeyViolationErrorCode = pq.ErrorCode("23503")
