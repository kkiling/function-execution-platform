package service

import "github.com/pkg/errors"

type ErrorServiceType int

const (
	InvalidInputDataErrorType ErrorServiceType = 0
	RuntimeErrorType          ErrorServiceType = 1
	NotFoundErrorType         ErrorServiceType = 2
	ConflictErrorType         ErrorServiceType = 3
)

type ErrorService struct {
	Type ErrorServiceType
	Err  error
}

func (r *ErrorService) Error() string {
	if r == nil || r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

func MakeInvalidInputDataErr(err error) error {
	return &ErrorService{Err: err, Type: InvalidInputDataErrorType}
}

func MakeInvalidInputDataWrapErr(err error, msg string) error {
	return &ErrorService{Err: errors.Wrap(err, msg), Type: InvalidInputDataErrorType}
}

func MakeRuntimeErr(err error) error {
	return &ErrorService{Err: err, Type: RuntimeErrorType}
}

func MakeRuntimeWrapErr(err error, msg string) error {
	return &ErrorService{Err: errors.Wrap(err, msg), Type: RuntimeErrorType}
}

func MakeNotFoundErr(err error) error {
	return &ErrorService{Err: err, Type: NotFoundErrorType}
}

func MakeNotFoundErrWrapErr(err error, msg string) error {
	return &ErrorService{Err: errors.Wrap(err, msg), Type: NotFoundErrorType}
}

func MakeConflictErr(err error) error {
	return &ErrorService{Err: err, Type: ConflictErrorType}
}
