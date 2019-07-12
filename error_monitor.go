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

import "sync"

type errorMonitor struct {
	onReadyOrErrorCond         *sync.Cond
	onReadyOrErrorLock         sync.Mutex
	onReadyOrErrorValue        error
	onReadyOrErrorHasCompleted bool
}

func New() MonitorNotifier {
	e := &errorMonitor{
		onReadyOrErrorLock: sync.Mutex{},
	}
	e.onReadyOrErrorCond = sync.NewCond(&e.onReadyOrErrorLock)
	return e
}

func (e *errorMonitor) Notify(err error) {
	e.onReadyOrErrorLock.Lock()
	defer e.onReadyOrErrorLock.Unlock()
	if e.onReadyOrErrorHasCompleted {
		panic("notify was called again, it should only be called ONCE")
	}
	e.onReadyOrErrorValue = err
	e.onReadyOrErrorHasCompleted = true
	e.onReadyOrErrorCond.Broadcast()
}

func (e *errorMonitor) WaitUntilNotified() error {
	e.onReadyOrErrorLock.Lock()
	defer e.onReadyOrErrorLock.Unlock()
	// If we've already gotten the value, do not wait
	if e.onReadyOrErrorHasCompleted {
		return e.onReadyOrErrorValue
	}
	// no value yet, wait until it arrives
	e.onReadyOrErrorCond.Wait()
	return e.onReadyOrErrorValue
}
