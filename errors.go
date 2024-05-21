package avs

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type Error struct {
	msg string
}

func NewAVSError(msg string) *Error {
	return &Error{msg: msg}
}

func NewAVSErrorFromGrpc(gErr error) *Error {
	status := status.Convert(gErr)

	if status.Err() == nil {
		return nil
	}

	errStr := fmt.Sprintf("avs server error: code=%s", status.Code().String())
	msg := status.Message()
	details := status.Details()

	if msg != "" {
		errStr = fmt.Sprintf(", msg=%s", msg)
	}

	if len(details) > 0 {
		errStr = fmt.Sprintf(", details=%v", details)
	}

	return NewAVSError(errStr)
}

func (e *Error) Error() string {
	return fmt.Sprintf("avs error: %v", e.msg)
}
