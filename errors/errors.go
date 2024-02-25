package errors

import (
	"errors"
	"fmt"
	gxe "golang.org/x/xerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"strings"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

func (e *Error) Error() string {
	var meta string
	if len(e.Metadata) > 0 {
		meta = fmt.Sprintf("\tmetadata = %v\t", e.Metadata)
	}
	var reason string
	if e.Reason != "" {
		reason = fmt.Sprintf("\treason[%s]\t", e.Reason)
	}
	if e.Xerror != nil {
		return strings.Replace(fmt.Sprintf("code[%d]%s message[%s]%s stack[%s] ", e.Code, reason, e.Message, meta, e.Xerror.Error()), "\n", "", -1)
	}

	return fmt.Sprintf("error:[%d]%s message[%s]%s stack[%s]", e.Code, reason, e.Message, meta, e.Stack)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

func (e *Error) WithMessage(msg string) *Error {
	err := proto.Clone(e).(*Error)
	err.Message = msg
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}

func (e *Error) Format(ptn string, param ...interface{}) *Error {
	e.Message += " :" + fmt.Sprintf(ptn, param...)
	return e
}

func (e *Error) SetBasicErr(err error) *Error {
	ret := NewWarp(int(e.Code), e.Reason, e.Message, gxe.Errorf("WARP : %w", err))
	return ret
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	return newError(code, reason, message, 0)
}

func newError(code int, reason, message string, skip int) *Error {
	err := &Error{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
	err.Stack = stack(skip)
	return err
}

func NewError(code int, msg string) *Error {
	err := newError(code, msg, msg, 0)
	return err
}

func NewWarp(code int, reason, message string, err error) *Error {
	e := newError(code, reason, message, 0)
	e.Xerror = err
	return e
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...interface{}) *Error {
	return newError(code, reason, fmt.Sprintf(format, a...), 0)
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...interface{}) error {
	return newError(code, reason, fmt.Sprintf(format, a...), 0)
}

func NewErrorSkip(code int, reason, message string, skip int) *Error {
	return newError(code, reason, message, skip)
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Message returns the message for a particular error.
// It supports wrapped errors.
func Message(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Message
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(
			httpstatus.FromGRPCCode(gs.Code()),
			UnknownReason,
			gs.Message(),
		)
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Reason = d.Reason
				return ret.WithMessage(gs.Message()).WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(UnknownCode, UnknownReason, err.Error())
}
