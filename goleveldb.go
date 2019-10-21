package main

type goleveldbitem struct {
	// add some return value at following fields
}

type goleveldbtask struct {
	writeMergeC  chan goleveldbitem
	writeMergedC chan bool
	writeLockC   chan struct{}
	writeAckC    chan error
	x            Tasker
}

func (db *goleveldbtask) Do() {
	var item goleveldbitem

	select {
	case db.writeMergeC <- item:
		if <-db.writeMergedC {
			// Write is merged.
			<-db.writeAckC
			return
		}
		// Write is not merged, the write lock is handed to us. Continue.
	case db.writeLockC <- struct{}{}:
	}
	db.x.Do()
merge:
	for {
		select {
		case incoming := <-db.writeMergeC:
			_ = incoming
			db.writeMergedC <- true
			db.x.Do()
			db.writeAckC <- nil
		default:
			break merge
		}
	}
	<-db.writeLockC
}

func NewGoLevelDBTask(x Tasker) *goleveldbtask {
	return &goleveldbtask{
		writeMergeC:  make(chan goleveldbitem),
		writeMergedC: make(chan bool),
		writeLockC:   make(chan struct{}, 1),
		writeAckC:    make(chan error),
		x:            x,
	}
}
