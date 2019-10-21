package main

import "unsafe"

// NotifyList is a ticket-based notification list.
// Approximation of notifyList in runtime/sema.go. Size and alignment must
// agree.
type NotifyList struct {
	wait   uint32
	notify uint32
	lock   uintptr
	head   unsafe.Pointer
	tail   unsafe.Pointer
}

// Ticket returns a ticket of notification list. We must invoke this method
// to generate a new ticket before every Wait() invoking.
func (l *NotifyList) Ticket() uint32 { return notifyListAdd(l) }

// Wait will park this goroutine by the hold ticket, and wait for someone
// wake it up.
func (l *NotifyList) Wait(ticket uint32) { notifyListWait(l, ticket) }

// Signal wake up someone.
func (l *NotifyList) Signal() { notifyListNotifyOne(l) }

// Broadcast wake up all
func (l *NotifyList) Broadcast() { notifyListNotifyAll(l) }

//go:linkname notifyListAdd sync.runtime_notifyListAdd
func notifyListAdd(l *NotifyList) uint32

//go:linkname notifyListWait sync.runtime_notifyListWait
func notifyListWait(l *NotifyList, t uint32)

//go:linkname notifyListCheck sync.runtime_notifyListCheck
func notifyListCheck(size uintptr)

//go:linkname notifyListNotifyOne sync.runtime_notifyListNotifyOne
func notifyListNotifyOne(l *NotifyList)

//go:linkname notifyListNotifyAll sync.runtime_notifyListNotifyAll
func notifyListNotifyAll(l *NotifyList)

func init() {
	var n NotifyList
	notifyListCheck(unsafe.Sizeof(n))
}
