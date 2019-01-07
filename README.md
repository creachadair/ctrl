# ctrl: Non-local exit handling for Go main packages

In case of error, Go `main` packages typically call [`os.Exit`][osexit] or
[`log.Fatal`][logfatal]. This causes the process to terminate immediately, and
deferred calls are not invoked. Calling `log.Panic`allows deferred calls to
run, but leaves a noisy log trace.

the [ctrl][ctrl] package provides a `Run` function that performs the main
action of a program. Within its dynamic extent, calls to `ctrl.Exit` and
`ctrl.Exitf` will panic back to `Run`, which will handle logging and exiting
from the process as specified.

[osexit]: https://godoc.org/os#Exit
[logfatal]: https://godoc.org/log#Fatal
[ctrl]: https://godoc.org/bitbucket.org/creachadair/ctrl

## Usage example

The following code outlines the use of `ctrl.Run`:

```go
import "bitbucket.org/creachadair/ctrl"

// A stub main to set up the control handlers.
func main() {
  // The default flagset hard-exits on error. You could set it to panic
  // on error instead, if you want to parse inside realMain.
  flag.Parse()

  ctrl.Run(realMain)

  // if realMain returns nil, control returns here.
  // if realMain returns a non-nil error or panics, this is not reached.
}

// The real program logic goes into this function.
func realMain() error {
  defer cleanup()  // some deferred cleanup task

  if err := doSomething(); err != nil {
     return err                          // failure, exit 1, no log
  } else if !stateIsvalid() {
     ctrl.Exitf(2, "State is invalid")   // failure, exit 2 with log
  }
  return nil                             // success
}
```
