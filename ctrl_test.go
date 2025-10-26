// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package ctrl_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/creachadair/ctrl"
	"github.com/creachadair/mds/mtest"
)

func init() { ctrl.SetPanic(true) }

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

func setLogger(t *testing.T, w io.Writer) {
	old := log.Writer()
	log.SetOutput(w)
	t.Cleanup(func() { log.SetOutput(old) })
}

// When control returns normally from main, Run should return.
func TestRunOK(t *testing.T) {
	ctrl.Run(func() error { return nil })
}

// When control returns with an error, Run should exit 1.
func TestRunError(t *testing.T) {
	exited := false
	ctrl.SetHook(codeHook(t, &exited, 1))

	mtest.MustPanic(t, func() {
		ctrl.Run(func() error { return errors.New("bogus") })
	})
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
		want any
	}{
		// A panic that produces a non-error.
		{func() error { panic("unwanted") }, "unwanted"},

		// A panic that produces an error value.
		{func() error { panic(testErr) }, testErr},
	}
	for _, test := range tests {
		got := mtest.MustPanic(t, func() {
			ctrl.Run(test.fn)
		})
		if got != test.want {
			t.Errorf("Panic: got %#v, want %#v", got, test.want)
		}
	}
}

// A call to Exit returns the reported code, but does not log.
func TestExit(t *testing.T) {
	var buf bytes.Buffer
	setLogger(t, &buf)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			ctrl.SetHook(codeHook(t, &exited, code))
			buf.Reset()
			mtest.MustPanic(t, func() {
				ctrl.Run(func() error { return ctrl.Exit(code) })
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
	setLogger(t, &buf)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			ctrl.SetHook(codeHook(t, &exited, code))
			buf.Reset()
			mtest.MustPanic(t, func() {
				ctrl.Run(func() error { return ctrl.Exitf(code, "msg") })
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
