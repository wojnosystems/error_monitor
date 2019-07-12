// Copyright 2019 Chris Wojno
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
// to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above
// copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR Scaling PARTICULAR PURPOSE AND NON-INFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
// AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package error_monitor

import (
	"errors"
	"testing"
)

func TestErrorMonitor_WaitUntilNotified_ErrorIsPassed(t *testing.T) {
	boomErr := errors.New("boom")
	cases := map[string]struct {
		notifiedError error
		expectedError error
	} {
		"nil": {
			notifiedError: nil,
			expectedError: nil,
		},
		"not-nil": {
			notifiedError: boomErr,
			expectedError: boomErr,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			n := New()
			go func() {
				n.Notify(c.notifiedError)
			}()
			actualError := n.WaitUntilNotified()
			errorsEqual(t, c.expectedError, actualError)
		})
	}
}

func TestErrorMonitor_WaitUntilNotified_MultipleRoutines(t *testing.T) {
	routineCount := 10
	expectedError := errors.New("boom")
	actualErrorChan := make(chan error, routineCount)
	n := New()
	// create lots of routines
	for i := 0; i < routineCount; i++ {
		go func() {
			// Read all of the errors
			actualErrorChan <- n.WaitUntilNotified()
		}()
	}
	// signal the error
	n.Notify(expectedError)

	// Ensure all routines got the SAME error
	for i := 0; i < routineCount; i++ {
		errorsEqual(t, expectedError, <- actualErrorChan )
	}
}

// TestErrorMonitor_WaitUntilNotified_NoReentrantDeadlock calling WaitUntilNotified after Notify should NOT deadlock
func TestErrorMonitor_WaitUntilNotified_NoReEntrantDeadlock(t *testing.T) {
	expectedError := errors.New("boom")
	n := New()
	// signal the error
	n.Notify(expectedError)
	// Get the error a few times, ensure it does not get stuck here
	_ = n.WaitUntilNotified()
	_ = n.WaitUntilNotified()
	_ = n.WaitUntilNotified()
	_ = n.WaitUntilNotified()
}

func TestErrorMonitor_Notify_PanicIfCalledTwice(t *testing.T) {
	n := New()
	n.Notify(nil)
	func() {
		defer func() {
			if r := recover(); r != nil {
				return
			}
			t.Error("expected panic, but didn't get one")
		}()
		n.Notify(nil)
	}()
}

func errorsEqual( t *testing.T, expectedErr, actualError error ) {
	if expectedErr != actualError {
		if expectedErr == nil {
			t.Errorf(`expected nil error, but got "%s"`, actualError.Error())
		} else {
			if actualError == nil {
				t.Errorf(`expected error: "%s" but got nil`, expectedErr.Error())
			} else {
				t.Errorf(`expected error: "%s" but got "%s"`, expectedErr.Error(), actualError.Error())
			}
		}
	}
}