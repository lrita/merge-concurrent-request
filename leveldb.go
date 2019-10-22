package main

import (
	"sync"

	"github.com/lrita/atomic1"
	"github.com/lrita/cache"
)

type leveldbitem struct {
	done atomic1.AtomicBool
	cond *sync.Cond
	// add some return value at following fields
}

type leveldbtask struct {
	lock sync.Mutex
	list []*leveldbitem
	x    Tasker
}

var litempool = cache.Cache{New: func() interface{} { return new(leveldbitem) }, Size: 128}

func (m *leveldbtask) Do() {
	item := litempool.Get().(*leveldbitem)
	item.done.Set(false)
	item.cond = sync.NewCond(&m.lock)

	m.lock.Lock()
	m.list = append(m.list, item)
	for !item.done.Get() && len(m.list) != 0 && m.list[0] != item {
		item.cond.Wait()
	}
	if item.done.Get() {
		// here can change to real work return value, e.g:
		// `return item.xx`
		m.lock.Unlock()
		litempool.Put(item)
		return
	}

	list := m.list
	last := len(m.list)
	m.lock.Unlock()

	for _, it := range list {
		// here can change to some real work, e.g:
		// `it.xx = do_xx()`
		m.x.Do()
		it.done.Set(true)
	}

	m.lock.Lock()

	for _, it := range list[1:] {
		it.cond.Signal()
	}

	copy(m.list, m.list[last:])
	list = m.list[len(m.list)-last:]
	for i := range list {
		list[i] = nil
	}
	m.list = m.list[:len(m.list)-last]
	if len(m.list) != 0 {
		m.list[0].cond.Signal()
	}

	m.lock.Unlock()
	litempool.Put(item)
	// here can change to real work return value, e.g:
	// `return item.xx`
	return
}

func NewLevelDBTask(x Tasker) *leveldbtask {
	return &leveldbtask{x: x}
}
