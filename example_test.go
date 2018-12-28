// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package ctrl_test

import (
	"fmt"

	"bitbucket.org/creachadair/ctrl"
)

func ExampleRun() {
	ctrl.Run(func() error {
		fmt.Println("This is main")
		return nil
	})
	fmt.Println("That was main")
	// Output:
	// This is main
	// That was main
}

// Hack: By having ExampleExit last, the fact that it terminates the example
// program with a zero status makes the example appear to pass. It means other
// examples beyond this point will not run, which is generally icky.

func ExampleExit() {
	ctrl.Run(func() error {
		fmt.Println("Hello")
		return ctrl.Exit(0)
	})
	fmt.Println("You don't see this")
	// Output:
	// Hello
}
