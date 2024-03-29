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

type Notifier interface {
	// Notify routines waiting on WaitUntilNotified that a value is ready.
	// @param err the error encountered, if any. It will be returned to WaitUntilNotified. OK to use nil to indicate that no error was encountered
	Notify(err error)
}

type Monitorer interface {
	// WaitUntilNotified will block the caller until Notify is called
	// @return the err value passed to Notify
	WaitUntilNotified() error
}

type MonitorNotifier interface {
	Monitorer
	Notifier
}