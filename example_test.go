// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package ctrl_test

import (
	"errors"
	"fmt"

	"github.com/creachadair/ctrl"
)

func catchPanic(f func()) (val any) {
	// For purposes of the examples, convert exits into panics.
	ctrl.SetPanic(true)

	// Log reported errors to stdout.
	ctrl.SetHook(func(code int, err error) {
		if code != 0 || err != nil {
			fmt.Printf("[exit] code=%d err=%v\n", code, err)
		}
	})

	defer func() { val = recover() }()
	f()
	return
}

func ExampleRun_success() {
	ctrl.Run(func() error {
		fmt.Println("This is main")
		return nil
	})
	fmt.Println("That was main")
	// Output:
	// This is main
	// That was main
}

func ExampleRun_failure() {
	// N.B. catchPanic prevents ctrl.Run from terminating the example
	// program. You do not need this in production.
	catchPanic(func() {
		ctrl.Run(func() error {
			fmt.Println("Hello")
			return errors.New("goodbye")
		})
		fmt.Println("You do not see this")
	})
	// Output:
	// Hello
	// [exit] code=1 err=goodbye
}

func ExampleRun_panic() {
	// N.B. catchPanic prevents ctrl.Run from terminating the example
	// program. You do not need this in production.
	v := catchPanic(func() {
		ctrl.Run(func() error {
			fmt.Println("Hello")
			panic("omgwtfbbq")
		})
		fmt.Println("You do not see this")
	})
	fmt.Println("panic:", v)
	// Output:
	// Hello
	// panic: omgwtfbbq
}

func ExampleExit() {
	// N.B. catchPanic prevents ctrl.Run from terminating the example
	// program. You do not need this in production.
	catchPanic(func() {
		ctrl.Run(func() error {
			fmt.Println("Hello")
			return ctrl.Exit(0)
		})
		fmt.Println("You don't see this")
	})
	// Output:
	// Hello
}

func ExampleExitf() {
	// N.B. catchPanic prevents ctrl.Run from terminating the example
	// program. You do not need this in production.
	catchPanic(func() {
		ctrl.Run(func() error {
			fmt.Println("Hello")
			return ctrl.Exitf(5, "everything is bad")
		})
		fmt.Println("You don't see this")
	})
	// Output:
	// Hello
	// [exit] code=5 err=everything is bad
}

func ExampleFatalf() {
	// N.B. catchPanic prevents ctrl.Run from terminating the example
	// program. You do not need this in production.
	catchPanic(func() {
		ctrl.Run(func() error {
			fmt.Println("Hello")
			return ctrl.Fatalf("badness: %d", 25)
		})
		fmt.Println("You don't see this")
	})
	// Output:
	// Hello
	// [exit] code=1 err=badness: 25
}
