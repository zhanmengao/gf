package errors

import (
	"fmt"
	"golang.org/x/xerrors"
	"testing"
)

type mockErr struct{}

func (*mockErr) Error() string {
	return "mock error"
}

func TestWarp(t *testing.T) {
	ErrBase := New(1, "ErrBase", "error1 msg")
	err := NewWarp(2, "error1 warp", "error1 msg", xerrors.Errorf("raiseError: %w", ErrBase))
	err2 := NewWarp(3, "error2 warp", "error2 msg", xerrors.Errorf("wrap#01: %w", err))
	err3 := NewWarp(4, "error3 warp", "error3 msg", xerrors.Errorf("wrap#02: %w", err2))
	//fmt.Printf("%+v\n", err3)
	fmt.Println(err3)
}
