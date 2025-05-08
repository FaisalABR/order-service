package error

import "errors"

var (
	ErrInternalServer    = errors.New("internal server error")
	ErrSqlQuery          = errors.New("database server failed to execute query")
	ErrTooManyRequest    = errors.New("too many request")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrNotFound          = errors.New("data not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidUploadFile = errors.New("invalid upload file")
	ErrSizeTooBig        = errors.New("file size too big")
	ErrForbiden          = errors.New("forbiden")
)

var GeneralErrors = []error{
	ErrInternalServer,
	ErrSqlQuery,
	ErrTooManyRequest,
	ErrNotFound,
	ErrInvalidToken,
	ErrForbiden,
}
