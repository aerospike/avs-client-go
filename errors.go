package avs

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type Error struct {
	msg string
}

func NewAVSError(msg string, err error) *Error {
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}

	return &Error{msg: msg}
}

func NewAVSErrorFromGrpc(msg string, gErr error) *Error {
	if gErr == nil {
		return nil
	}

	gStatus, ok := status.FromError(gErr)
	if !ok {
		return NewAVSError(msg, gErr)
	}

	errStr := fmt.Sprintf("%s: server error: %s", msg, gStatus.Code().String())
	gMsg := gStatus.Message()
	details := gStatus.Details()

	if gMsg != "" {
		errStr = fmt.Sprintf("%s, msg=%s", errStr, gMsg)
	}

	if len(details) > 0 {
		errStr = fmt.Sprintf("%s, details=%v", errStr, details)
	}

	return NewAVSError(errStr, nil)
}

func (e *Error) Error() string {
	return e.msg
}
