# Overview

A way to broadcast when some service or go-routine is ready or if it failed with an error.

# Getting it

```
go get github.com/wojnosystems/errormonitor
```

# How to use

```go
import error_monitor

func main() {
	someService := error_monitor.New()
	
	for i := 0; i < 200; i++ {
		go func() {
			// spin up 200 routines, waiting for someService to be ready or get the error if it failed
			whichError := someService.WaitUntilNotified()
			if whichError != nil {
				// it failed!
				return
			}
			// do something, it worked!
		}()
	}
	
	someService.Notify(errors.New("ka-boom! oh-noes, world is fire!"))
	
	// alternatively, you can pass in nil to indicate success
	// someService.Notify(nil)
}
```

# Motivation

I needed a way to notify tests when a service that ran in another routine was ready. Channels didn't work because they operated on a single value and I wanted to be able to notify anyone waiting that the value was ready. The guts of this object lived in a class and it became clear that it needed encapsulation. So now, this can be used by anyone who needs it. If you have some main routine that's building something and lots of other routines waiting for that thing to finish, this is a simple and great structure to use.