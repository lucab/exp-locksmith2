package lock

import (
	"reflect"
	"testing"
)

func TestSingleLock(t *testing.T) {
	sem := NewSemaphore(1)

	if sem.Max != 1 {
		t.Errorf("unexpected semaphore size: %d", sem.Max)
	}

	held, err := sem.RecursiveLock("a")
	if err != nil {
		t.Error(err)
	}
	if held {
		t.Error("unexpected holding lock")
	}
	if !reflect.DeepEqual(sem.Holders, []string{"a"}) {
		t.Error("lock did not add a to the holders")
	}
	if sem.Max != 1 {
		t.Errorf("unexpected semaphore size: %d", sem.Max)
	}

	if err := sem.UnlockIfHeld("a"); err != nil {
		t.Error(err)
	}
	if len(sem.Holders) != 0 {
		t.Error("lock did not remove a from the holders")
	}
	if sem.Max != 1 {
		t.Errorf("unexpected semaphore size: %d", sem.Max)
	}
}

func TestRecursivelock(t *testing.T) {
	sem := NewSemaphore(1)

	heldOne, err := sem.RecursiveLock("a")
	if err != nil {
		t.Error(err)
	}
	if heldOne {
		t.Error("unexpected holding lock")
	}

	heldTwo, err := sem.RecursiveLock("a")
	if err != nil {
		t.Error(err)
	}
	if !heldTwo {
		t.Error("unexpected not holding lock")
	}

	if err := sem.UnlockIfHeld("a"); err != nil {
		t.Error(err)
	}
}

func TestHolderOrdering(t *testing.T) {
	sem := NewSemaphore(3)

	if _, err := sem.RecursiveLock("c"); err != nil {
		t.Error(err)
	}
	if _, err := sem.RecursiveLock("a"); err != nil {
		t.Error(err)
	}
	if _, err := sem.RecursiveLock("b"); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(sem.Holders, []string{"a", "b", "c"}) {
		t.Error("unexpected ordering")
	}
	if err := sem.UnlockIfHeld("b"); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(sem.Holders, []string{"a", "c"}) {
		t.Error("unexpected ordering")
	}
}

func TestSetMax(t *testing.T) {
	oldMax := uint64(1)
	newMax := uint64(4)
	sem := NewSemaphore(oldMax)

	if err := sem.SetMax(newMax); err != nil {
		t.Error(err)
	}

	if sem.Max != newMax {
		t.Errorf("unexpected semaphore size: got %d, expectedd %d", sem.Max, newMax)

	}
}
