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

func setExit(f func(int)) func() {
	old := osExit
	osExit = f
	return func() { osExit = old }
}

// When control returns normally from main, Run should return.
func TestRunOK(t *testing.T) {
	defer setExit(func(code int) {
		t.Fatalf("Unexpected exit: code=%d", code)
	})()

	Run(func() error { return nil })
}

// When control returns with an error, Run should exit 1.
func TestRunError(t *testing.T) {
	exited := false
	defer setExit(func(code int) {
		exited = true
		if code != 1 {
			t.Errorf("Unexpected exit code: %v", code)
		}
	})()

	Run(func() error { return errors.New("bogus") })
	if !exited {
		t.Fatal("The exit handler did not get invoked")
	}
}

// Panics from other sources get propagated.
func TestRunPanic(t *testing.T) {
	const msg = "unrelated"
	defer setExit(func(code int) { t.Logf("osExit code %d", code) })()
	defer func() {
		v := recover()
		if v == nil {
			t.Error("No panic was observed")
		} else if v != msg {
			t.Errorf("Panic: got %v, want %q", v, msg)
		}
	}()
	Run(func() error { panic(msg) })
	t.Fatal("Reached the unreachable star")
}

// A call to Exit returns the reported code, but does not log.
func TestExit(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	for code := 0; code <= 5; code++ {
		t.Run(fmt.Sprintf("Code-%d", code), func(t *testing.T) {
			exited := false
			defer setExit(func(got int) {
				exited = true
				if got != code {
					t.Errorf("Wrong code: got %d, want %d", got, code)
				}
			})()
			buf.Reset()
			Run(func() error { return Exit(code) })
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
			defer setExit(func(got int) {
				exited = true
				if got != code {
					t.Errorf("Wrong code: got %d, want %d", got, code)
				}
			})()
			buf.Reset()
			Run(func() error { return Exitf(code, "msg") })
			if code != 0 && !exited {
				t.Errorf("Code %d: exit handler not invoked", code)
			}
			if s := buf.String(); !strings.HasSuffix(s, " msg\n") {
				t.Errorf("Code %d: unexpected log output: %q", code, s)
			}
		})
	}
}
