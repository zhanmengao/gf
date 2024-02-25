package errors

import "testing"

func TestStack(t *testing.T) {
	err := New(100, "err", "error")
	t.Log(err.Error())
	err = BadRequest("err", "msg")
	t.Log(err.Error())
}
