package wasm

import (
	"fmt"
	"sync"
)

type AccountInfo struct {
	AccountNumber uint64
	Sequence      uint64
}

type SequenceTracker struct {
	accounts map[string]*AccountInfo // maps account's address to accountInfo
	*sync.Mutex
}

func NewSeqTracker(accounts map[string]*AccountInfo) *SequenceTracker {
	return &SequenceTracker{
		accounts: accounts,
		Mutex:    new(sync.Mutex),
	}
}

func (s *SequenceTracker) Set(address string, ac *AccountInfo) error {
	s.Lock()
	defer s.Unlock()
	acInfo, ok := s.accounts[address]
	if !ok {
		return fmt.Errorf("failed to set sequence: address %s not found in sequence tracker", address)
	}
	acInfo.Sequence = ac.Sequence
	acInfo.AccountNumber = ac.AccountNumber
	s.accounts[address] = acInfo
	return nil
}

func (s *SequenceTracker) GetWithLock(address string) (*AccountInfo, error) {
	s.Lock()
	defer s.Unlock()
	currAcInfo, ok := s.accounts[address]
	if !ok {
		return nil, fmt.Errorf("failed to get sequence with lock: address %s not found in sequence tracker", address)
	}
	s.accounts[address] = &AccountInfo{
		AccountNumber: currAcInfo.AccountNumber,
		Sequence:      currAcInfo.Sequence + 1,
	}
	return currAcInfo, nil
}

// Get use this method with caution. Requires explicit lock handling.
func (s *SequenceTracker) Get(address string) (*AccountInfo, error) {
	currAcInfo, ok := s.accounts[address]
	if !ok {
		return nil, fmt.Errorf("failed to get sequence: address %s not found in sequence tracker", address)
	}
	return currAcInfo, nil
}

// IncrementSequence use this method with caution. Requires explicit lock handling.
func (s *SequenceTracker) IncrementSequence(address string) error {
	currAcInfo, ok := s.accounts[address]
	if !ok {
		return fmt.Errorf("failed to increment sequence: address %s not found in sequence tracker", address)
	}
	s.accounts[address] = &AccountInfo{
		AccountNumber: currAcInfo.AccountNumber,
		Sequence:      currAcInfo.Sequence + 1,
	}
	return nil
}
