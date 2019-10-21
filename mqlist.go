package main

import (
	"sync"
	"unsafe"

	"github.com/lrita/atomic1"
)

type Tasker interface {
	Do()
}

type mqitem struct {
	done atomic1.AtomicBool
	// add some return value at following fields
}

type mqtask struct {
	mqflag   atomic1.AtomicBool
	mqlock   sync.Mutex
	mqnotify NotifyList
	mqidx    uint64
	mqlist   [2][]*mqitem
	x        Tasker
}

func (m *mqtask) Do() {
	var itemobj mqitem
	item := (*mqitem)(noescape(unsafe.Pointer(&itemobj)))

	m.mqlock.Lock()
	idx := m.mqidx & 1
	m.mqlist[idx] = append(m.mqlist[idx], item)
	m.mqlock.Unlock()

	for {
		ticket := m.mqnotify.Ticket()
		if !m.mqflag.Get() && m.mqflag.CAS(true) { // get first to reduce cpu bus traffic
			m.mqlock.Lock()
			idx := m.mqidx & 1
			list := m.mqlist[idx]
			// switch and lock the pending list, the goroutine who is holding
			// mqflag can access the list
			m.mqidx++
			m.mqlock.Unlock()
			for _, it := range list {
				// here can change to some real work, e.g:
				// `it.xx = do_xx()`
				m.x.Do()
				it.done.Set(true)
			}

			// optimize to memcpy/memmove, overwrite to nil to avoid memory leakage
			for i := range list {
				list[i] = nil
			}
			m.mqlock.Lock()
			m.mqlist[idx] = m.mqlist[idx][:0]
			m.mqlock.Unlock()
			m.mqflag.Set(false)
			m.mqnotify.Broadcast()
			break
		}
		m.mqnotify.Wait(ticket)
		if item.done.Get() { // merge done
			break
		}
	}
	// here can change to real work return value, e.g:
	// `return item.xx`
}

func NewMQTask(x Tasker) *mqtask {
	return &mqtask{x: x}
}

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	//lint:ignore SA4016 copy from runtime package, disable staticcheck blame.
	return unsafe.Pointer(x ^ 0)
}
