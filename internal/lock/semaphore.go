package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var (
	// ErrExist is the error returned if a holder with the specified id
	// is already holding the semaphore
	ErrExist = errors.New("holder exists")
	// ErrNotExist is the error returned if there is no holder with the
	// specified id holding the semaphore
	ErrNotExist = errors.New("holder does not exist")
	// ErrNilSemaphore is returned on nil semaphore.
	ErrNilSemaphore = errors.New("nil Semaphore")
)

// Semaphore is a struct representation of the information held by the semaphore
type Semaphore struct {
	Semaphore uint64   `json:"semaphore"`
	Max       uint64   `json:"max"`
	Holders   []string `json:"holders"`
}

// NewSemaphore returns a new empty semaphore.
func NewSemaphore() (sem *Semaphore) {
	return &Semaphore{1, 1, []string{}}
}

// SetMax sets the maximum number of holders of the semaphore
func (s *Semaphore) SetMax(max uint64) error {
	if s == nil {
		return ErrNilSemaphore
	}

	diff := s.Max - max

	s.Semaphore = s.Semaphore - diff
	s.Max = s.Max - diff

	return nil
}

// String returns a json representation of the semaphore
// if there is an error when marshalling the json, it is ignored and the empty
// string is returned.
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
// it returns ErrExist if the given id is in the list
func (s *Semaphore) addHolder(h string) error {
	if s == nil {
		return ErrNilSemaphore
	}

	loc := sort.SearchStrings(s.Holders, h)
	switch {
	case loc == len(s.Holders):
		s.Holders = append(s.Holders, h)
	default:
		s.Holders = append(s.Holders[:loc], append([]string{h}, s.Holders[loc:]...)...)
	}

	s.Semaphore = s.Semaphore - 1
	return nil
}

// removeHolder removes a holder with id h from the list of holders in the
// semaphore. It returns whether the holder was present in the list.
func (s *Semaphore) removeHolder(h string) (bool, error) {
	if s == nil {
		return false, ErrNilSemaphore
	}

	loc := sort.SearchStrings(s.Holders, h)
	if loc < len(s.Holders) && s.Holders[loc] == h {
		s.Holders = append(s.Holders[:loc], s.Holders[loc+1:]...)
		s.Semaphore = s.Semaphore + 1
		return true, nil
	}

	return false, nil
}

// RecursiveLock adds a holder with id h to the semaphore
// It adds the id h to the list of holders, returning ErrExist the id already
// exists, then it subtracts one from the semaphore. If the semaphore is already
// held by the maximum number of people it returns an error.
func (s *Semaphore) RecursiveLock(id string) (bool, error) {
	if s == nil {
		return false, ErrNilSemaphore
	}

	// Check if id is already holding a lock.
	loc := sort.SearchStrings(s.Holders, id)
	if loc < len(s.Holders) && s.Holders[loc] == id {
		return true, nil
	}

	if s.Semaphore <= 0 {
		return false, fmt.Errorf("semaphore is at %v", s.Semaphore)
	}

	if err := s.addHolder(id); err != nil {
		return false, err
	}

	return false, nil
}

// UnlockIfHeld removes a holder with id h from the semaphore
// It removes the id h from the list of holders, returning ErrNotExist if the id
// does not exist in the list, then adds one to the semaphore.
func (s *Semaphore) UnlockIfHeld(h string) error {
	if s == nil {
		return ErrNilSemaphore
	}

	_, err := s.removeHolder(h)
	if err != nil {
		return err
	}

	return nil
}
