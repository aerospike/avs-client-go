package avs

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type Error struct {
	msg string
}

func NewAVSError(msg string) error {
	return &Error{msg: msg}
}

func NewAVSErrorFromGrpc(msg string, gErr error) error {
	gStatus, ok := status.FromError(gErr)

	if gErr == nil {
		return nil
	}

	if !ok {
		return NewAVSError(gErr.Error())
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

	return NewAVSError(errStr)
}

func (e *Error) Error() string {
	return e.msg
}
