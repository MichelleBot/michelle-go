package middleware

import (
	"fmt"
	"sync"
	"time"
)

type limitRecord struct {
	count      int
	resetAt    time.Time
	lastAccess time.Time
}

type CommandLimiter struct {
	mu      sync.Mutex
	records map[string]*limitRecord
}

func NewCommandLimiter() *CommandLimiter {
	l := &CommandLimiter{
		records: make(map[string]*limitRecord),
	}
	go l.cleanup()
	return l
}

func (l *CommandLimiter) Allow(commandName, userID string, max int, window time.Duration) (bool, time.Duration) {
	if max <= 0 || window <= 0 {
		return true, 0
	}

	now := time.Now()
	key := fmt.Sprintf("%s|%s", commandName, userID)

	l.mu.Lock()
	defer l.mu.Unlock()

	rec, ok := l.records[key]
	if !ok || now.After(rec.resetAt) {
		rec = &limitRecord{count: 0, resetAt: now.Add(window), lastAccess: now}
		l.records[key] = rec
	}

	rec.lastAccess = now
	if rec.count >= max {
		retryAfter := rec.resetAt.Sub(now)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	rec.count++
	return true, 0
}

func (l *CommandLimiter) cleanup() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		now := time.Now()
		l.mu.Lock()
		for key, rec := range l.records {
			if now.Sub(rec.lastAccess) > 2*time.Hour {
				delete(l.records, key)
			}
		}
		l.mu.Unlock()
	}
}
