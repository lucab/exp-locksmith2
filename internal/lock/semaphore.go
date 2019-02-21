package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var (
	// ErrNilSemaphore is returned on nil semaphore.
	ErrNilSemaphore = errors.New("nil Semaphore")
)

// Semaphore is a struct representation of the information held by the semaphore
type Semaphore struct {
	TotalSlots uint64   `json:"total_slots"`
	Holders    []string `json:"holders"`
}

// NewSemaphore returns a new empty semaphore.
func NewSemaphore(slots uint64) (sem *Semaphore) {
	return &Semaphore{slots, []string{}}
}

// SetTotalSlots sets the number of holders slots for the semaphore
func (s *Semaphore) SetTotalSlots(slots uint64) error {
	if s == nil {
		return ErrNilSemaphore
	}

	if int(slots) < len(s.Holders) {
		return fmt.Errorf("failed to set max to %d, %d current holders", slots, len(s.Holders))
	}

	s.TotalSlots = slots

	return nil
}

// String returns a json representation of the semaphore.
func (s *Semaphore) String() (string, error) {
	if s == nil {
		return "", ErrNilSemaphore
	}

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// addHolder adds a holder with id h to the list of holders in the semaphore
func (s *Semaphore) addHolder(h string) error {
	if s == nil {
		return ErrNilSemaphore
	}
	if len(s.Holders) >= int(s.TotalSlots) {
		return fmt.Errorf("all %d semaphore slots currently locked", s.TotalSlots)
	}

	loc := sort.SearchStrings(s.Holders, h)
	switch {
	case loc == len(s.Holders):
		s.Holders = append(s.Holders, h)
	default:
		s.Holders = append(s.Holders[:loc], append([]string{h}, s.Holders[loc:]...)...)
	}

	return nil
}

// removeHolder removes a holder with id h from the list of holders in the
// semaphore. It returns whether the holder was present in the list.
func (s *Semaphore) removeHolderIfPresent(h string) (bool, error) {
	if s == nil {
		return false, ErrNilSemaphore
	}

	loc := sort.SearchStrings(s.Holders, h)
	if loc < len(s.Holders) && s.Holders[loc] == h {
		s.Holders = append(s.Holders[:loc], s.Holders[loc+1:]...)
		return true, nil
	}

	return false, nil
}

// RecursiveLock adds a holder with id h to the semaphore,
// or returns an error if the semaphore is already a maximum
// capacity.
func (s *Semaphore) RecursiveLock(id string) (bool, error) {
	if s == nil {
		return false, ErrNilSemaphore
	}

	// Check if id is already holding a lock.
	loc := sort.SearchStrings(s.Holders, id)
	if loc < len(s.Holders) && s.Holders[loc] == id {
		return true, nil
	}

	if err := s.addHolder(id); err != nil {
		return false, err
	}

	return false, nil
}

// UnlockIfHeld removes a holder with id h from the semaphore, if present.
func (s *Semaphore) UnlockIfHeld(h string) error {
	if s == nil {
		return ErrNilSemaphore
	}

	_, err := s.removeHolderIfPresent(h)
	if err != nil {
		return err
	}

	return nil
}
