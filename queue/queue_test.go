package queue

import (
	"errors"
	"testing"
)

func quecmdBuilder(name string, fail, panics bool) (func() error, *bool) {
	notrun := false
	hasrun := &notrun

	return func() error {
		*hasrun = true
		if fail {
			return errors.New(name + " failed")
		}
		if panics {
			panic("a panic")
		}
		return nil
	}, hasrun

}

func TestQueueError(t *testing.T) {
	af, ahasrun := quecmdBuilder("a", true, false)
	bf, bhasrun := quecmdBuilder("b", false, false)

	q := NewQueue()
	q.Add(af, bf)
	err := q.Run()
	if !*ahasrun {
		t.Fatalf("af should run")
	}
	if *bhasrun {
		t.Fatalf("bf should not run")
	}
	if err == nil {
		t.Fatalf("should return an error")
	}
}

func TestQueuePanic(t *testing.T) {
	af, ahasrun := quecmdBuilder("a", false, true)
	bf, bhasrun := quecmdBuilder("b", false, true)

	q := NewQueue()
	q.Add(af, bf)
	err := q.Run()
	if !*ahasrun {
		t.Fatalf("af should run")
	}
	if *bhasrun {
		t.Fatalf("bf should not run")
	}
	if err == nil {
		t.Fatalf("should return an error")
	}
}
