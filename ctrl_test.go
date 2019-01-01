// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package ctrl

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func init() {
	SetPanic(true)
}

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
	Run(func() error { return nil })
}

// When control returns with an error, Run should exit 1.
func TestRunError(t *testing.T) {
	exited := false
	SetHook(codeHook(t, &exited, 1))

	p := panicOK(t, func() {
		Run(func() error { return errors.New("bogus") })
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
	const msg = "unrelated"
	SetHook(logHook(t))
	got := panicOK(t, func() {
		Run(func() error { panic(msg) })
		t.Fatalf("Reached the unreachable star")
	})
	if got == nil {
		t.Error("No panic was observed")
	} else if got != msg {
		t.Errorf("Panic: got %v, want %q", got, msg)
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
			SetHook(codeHook(t, &exited, code))
			buf.Reset()
			panicOK(t, func() {
				Run(func() error { return Exit(code) })
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

// A call to Exitf returns the reported code, and logs what it got.
func TestExitf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			SetHook(codeHook(t, &exited, code))
			buf.Reset()
			panicOK(t, func() {
				Run(func() error { return Exitf(code, "msg") })
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
