// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package ctrl_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/creachadair/ctrl"
)

func init() { ctrl.SetPanic(true) }

// panicOK runs a function f that is expected to panic, and returns the value
// recovered from that panic (if it does) or nil.
func panicOK(t *testing.T, f func()) (val interface{}) {
	t.Helper()

	defer func() {
		val = recover()
		if val != nil {
			t.Logf("Panic captured with value: %v", val)
		}
	}()
	f()
	return
}

func codeHook(t *testing.T, ex *bool, want int) func(int, error) {
	return func(code int, err error) {
		t.Helper()
		*ex = true
		if code != want {
			t.Errorf("Unexpected exit code: got %d, want %d (err=%v", code, want, err)
		}
	}
}

func logHook(t *testing.T) func(int, error) {
	return func(code int, err error) {
		t.Helper()
		t.Logf("Exit code %d, error %v", code, err)
	}
}

// When control returns normally from main, Run should return.
func TestRunOK(t *testing.T) {
	ctrl.Run(func() error { return nil })
}

// When control returns with an error, Run should exit 1.
func TestRunError(t *testing.T) {
	exited := false
	ctrl.SetHook(codeHook(t, &exited, 1))

	p := panicOK(t, func() {
		ctrl.Run(func() error { return errors.New("bogus") })
	})
	if p == nil {
		t.Error("No panic was observed")
	}
	if !exited {
		t.Fatal("The exit handler did not get invoked")
	}
}

// Panics from other sources get propagated.
func TestRunPanic(t *testing.T) {
	testErr := errors.New("this is not an error")
	ctrl.SetHook(logHook(t))
	tests := []struct {
		fn   func() error
		want interface{}
	}{
		// A panic that produces a non-error.
		{func() error { panic("unwanted") }, "unwanted"},

		// A panic that produces an error value.
		{func() error { panic(testErr) }, testErr},
	}
	for _, test := range tests {
		got := panicOK(t, func() {
			ctrl.Run(test.fn)
			t.Fatalf("Reached the unreachable star")
		})
		if got == nil {
			t.Error("No panic was observed")
		} else if got != test.want {
			t.Errorf("Panic: got %#v, want %#v", got, test.want)
		}
	}
}

// A call to Exit returns the reported code, but does not log.
func TestExit(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			ctrl.SetHook(codeHook(t, &exited, code))
			buf.Reset()
			panicOK(t, func() {
				ctrl.Run(func() error { return ctrl.Exit(code) })
				t.Error("No panic was observed")
			})
			if code != 0 && !exited {
				t.Errorf("Code %d: exit handler not invoked", code)
			}
			if buf.Len() != 0 {
				t.Errorf("Code %d: unexpected log output: %v", code, buf.String())
			}
		})
	}
}

// A call to ctrl.Exitf returns the reported code, and logs what it got.
func TestExitf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			ctrl.SetHook(codeHook(t, &exited, code))
			buf.Reset()
			panicOK(t, func() {
				ctrl.Run(func() error { return ctrl.Exitf(code, "msg") })
				t.Error("No panic was observed")
			})
			if code != 0 && !exited {
				t.Errorf("Code %d: exit handler not invoked", code)
			}
			if s := buf.String(); !strings.HasSuffix(s, " msg\n") {
				t.Errorf("Code %d: unexpected log output: %q", code, s)
			}
		})
	}
}
