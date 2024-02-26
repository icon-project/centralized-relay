package wasm

import (
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type AccountInfo struct {
	AccountNumber uint64
	Sequence      uint64
}

type SequenceTracker struct {
	accounts map[string]*AccountInfo // maps account's address to accountInfo
	*sync.Mutex
}

func (p *Provider) NewSeqTracker(addr types.AccAddress) *SequenceTracker {
	accounts := map[string]*AccountInfo{
		addr.String(): {
			AccountNumber: p.wallet.GetAccountNumber(),
			Sequence:      p.wallet.GetSequence(),
		},
	}
	return &SequenceTracker{
		accounts: accounts,
		Mutex:    new(sync.Mutex),
	}
}

func (s *SequenceTracker) Set(address types.AccAddress, ac *AccountInfo) error {
	s.Lock()
	defer s.Unlock()
	acInfo, ok := s.accounts[address.String()]
	if !ok {
		return fmt.Errorf("failed to set sequence: address %s not found in sequence tracker", address)
	}
	acInfo.Sequence = ac.Sequence
	acInfo.AccountNumber = ac.AccountNumber
	s.accounts[address.String()] = acInfo
	return nil
}

func (s *SequenceTracker) GetWithLock(address sdkTypes.AccAddress) (*AccountInfo, error) {
	s.Lock()
	defer s.Unlock()
	currAcInfo, ok := s.accounts[address.String()]
	if !ok {
		return nil, fmt.Errorf("failed to get sequence with lock: address %s not found in sequence tracker", address)
	}
	s.accounts[address.String()] = &AccountInfo{
		AccountNumber: currAcInfo.AccountNumber,
		Sequence:      currAcInfo.Sequence + 1,
	}
	return currAcInfo, nil
}

// Get use this method with caution. Requires explicit lock handling.
func (s *SequenceTracker) Get(address sdkTypes.AccAddress) (*AccountInfo, error) {
	currAcInfo, ok := s.accounts[address.String()]
	if !ok {
		return nil, fmt.Errorf("failed to get sequence: address %s not found in sequence tracker", address)
	}
	return currAcInfo, nil
}

// IncrementSequence use this method with caution. Requires explicit lock handling.
func (s *SequenceTracker) IncrementSequence(address types.AccAddress) error {
	currAcInfo, ok := s.accounts[address.String()]
	if !ok {
		return fmt.Errorf("failed to increment sequence: address %s not found in sequence tracker", address)
	}
	s.accounts[address.String()] = &AccountInfo{
		AccountNumber: currAcInfo.AccountNumber,
		Sequence:      currAcInfo.Sequence + 1,
	}
	return nil
}
