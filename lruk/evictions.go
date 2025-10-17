package lruk

import (
	"math"
	"time"
)

func (l *LRUKCache) SetEvictable(frameID int, evictable bool) {
	l.locker.Lock()
	defer l.locker.Unlock()

	l.evictable[frameID] = evictable
}

func (l *LRUKCache) Evict() (int, bool) {
	l.locker.Lock()
	defer l.locker.Unlock()

	var (
		victim, maxDist = -1, int64(-1)
		now             = time.Now().UnixNano()
	)

	for frameID, accessHistory := range l.accessHistory {
		if !l.evictable[frameID] {
			continue
		}

		var dist int64 = math.MaxInt64

		if len(accessHistory) >= l.k {
			dist = now - accessHistory[0]
		}

		if dist > maxDist {
			maxDist = dist
			victim = frameID
		}
	}

	if victim == -1 {
		return -1, false
	}

	delete(l.accessHistory, victim)
	delete(l.evictable, victim)

	return victim, true
}
