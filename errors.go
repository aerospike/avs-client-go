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
	status, ok := status.FromError(gErr)

	if gErr == nil {
		return nil
	}

	if !ok {
		// Should we instead return
		return NewAVSError(gErr.Error())
	}

	errStr := fmt.Sprintf("%s: avs server error: code=%s", msg, status.Code().String())
	gMsg := status.Message()
	details := status.Details()

	if gMsg != "" {
		errStr = fmt.Sprintf(", msg=%s", gMsg)
	}

	if len(details) > 0 {
		errStr = fmt.Sprintf(", details=%v", details)
	}

	return NewAVSError(errStr)
}

func (e *Error) Error() string {
	return e.msg
}
