package wasm

import "sync"

type MemPoolInfo struct {
	isBlocked bool
	sync.Mutex
}

func (mp *MemPoolInfo) SetBlockedStatus(val bool) {
	mp.isBlocked = val
}

func (mp *MemPoolInfo) SetBlockedStatusWithLock(val bool) {
	mp.Lock()
	defer mp.Unlock()
	mp.isBlocked = val
}

func (mp *MemPoolInfo) IsBlocked() bool {
	return mp.isBlocked
}
