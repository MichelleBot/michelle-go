package middleware

import (
	"sync"
	"time"
)

type userRecord struct {
	countPerSec  int
	countPerMin  int
	lastSecReset time.Time
	lastMinReset time.Time
	lastActivity time.Time
	bannedUntil  time.Time
}

type Antispam struct {
	mu          sync.Mutex
	records     map[string]*userRecord
	maxPerSec   int
	maxPerMin   int
	banDuration time.Duration
}

func NewAntispam(maxPerSec, maxPerMin, banSecs int) *Antispam {
	a := &Antispam{
		records:     make(map[string]*userRecord),
		maxPerSec:   maxPerSec,
		maxPerMin:   maxPerMin,
		banDuration: time.Duration(banSecs) * time.Second,
	}
	go a.cleanup()
	return a
}

func (a *Antispam) Check(userID string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	rec, ok := a.records[userID]
	if !ok {
		rec = &userRecord{lastSecReset: now, lastMinReset: now}
		a.records[userID] = rec
	}
	rec.lastActivity = now

	if now.Before(rec.bannedUntil) {
		return false
	}

	if now.Sub(rec.lastSecReset) >= time.Second {
		rec.countPerSec = 0
		rec.lastSecReset = now
	}
	if now.Sub(rec.lastMinReset) >= time.Minute {
		rec.countPerMin = 0
		rec.lastMinReset = now
	}

	rec.countPerSec++
	rec.countPerMin++

	if rec.countPerSec > a.maxPerSec || rec.countPerMin > a.maxPerMin {
		rec.bannedUntil = now.Add(a.banDuration)
		return false
	}

	return true
}

func (a *Antispam) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		a.mu.Lock()
		now := time.Now()
		for id, rec := range a.records {
			if now.Sub(rec.lastActivity) > 30*time.Minute {
				delete(a.records, id)
			}
		}
		a.mu.Unlock()
	}
}
