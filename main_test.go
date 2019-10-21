package main

import (
	"testing"
	"time"
)

type slowtask struct {
	count int
}

func (t *slowtask) Do() {
	t.count++
	time.Sleep(time.Millisecond)
}

type fasttask struct {
	count int
}

func (t *fasttask) Do() {
	t.count++
}

func BenchmarkMQList(b *testing.B) {
	x := &fasttask{}
	m := NewMQTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}

func BenchmarkLevelDB(b *testing.B) {
	x := &fasttask{}
	m := NewLevelDBTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}

func BenchmarkGoLevelDB(b *testing.B) {
	x := &fasttask{}
	m := NewGoLevelDBTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}

func BenchmarkMQListSlow(b *testing.B) {
	x := &slowtask{}
	m := NewMQTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}

func BenchmarkLevelDBSlow(b *testing.B) {
	x := &slowtask{}
	m := NewLevelDBTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}

func BenchmarkGoLevelDBSlow(b *testing.B) {
	x := &slowtask{}
	m := NewGoLevelDBTask(x)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Do()
		}
	})
	b.StopTimer()
	if x.count != b.N {
		b.Fatalf("x got(%v), expect(%v)", x.count, b.N)
	}
}
