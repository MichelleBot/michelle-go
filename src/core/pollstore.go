package core

import (
	"fmt"
	"sync"
	"time"
)

type PollState struct {
	ID          string
	Chat        string
	Name        string
	OptionCount int
	UpdateCount int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PollStore struct {
	mu   sync.RWMutex
	data map[string]*PollState
}

func NewPollStore() *PollStore {
	return &PollStore{data: make(map[string]*PollState)}
}

func (ps *PollStore) SaveCreation(chat string, poll *NormalizedPoll, ts time.Time) {
	if ps == nil || poll == nil || poll.ID == "" {
		return
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.data[ps.key(chat, poll.ID)] = &PollState{
		ID:          poll.ID,
		Chat:        chat,
		Name:        poll.Name,
		OptionCount: poll.OptionCount,
		CreatedAt:   ts,
		UpdatedAt:   ts,
	}
}

func (ps *PollStore) RegisterUpdate(chat string, poll *NormalizedPoll, ts time.Time) *PollState {
	if ps == nil || poll == nil || poll.TargetID == "" {
		return nil
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	state, ok := ps.data[ps.key(chat, poll.TargetID)]
	if !ok {
		state = &PollState{ID: poll.TargetID, Chat: chat, CreatedAt: ts, UpdatedAt: ts}
		ps.data[ps.key(chat, poll.TargetID)] = state
	}
	state.UpdateCount += poll.UpdateCount
	state.UpdatedAt = ts
	return state
}

func (ps *PollStore) Get(chat, id string) (*PollState, bool) {
	if ps == nil {
		return nil, false
	}
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	state, ok := ps.data[ps.key(chat, id)]
	return state, ok
}

func (ps *PollStore) key(chat, id string) string {
	return fmt.Sprintf("%s|%s", chat, id)
}
