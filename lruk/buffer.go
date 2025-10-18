package lruk

import "github.com/dark-vinci/nildb/interfaces"

func (l *LRUKCache) Index(frameID int) *interfaces.RepPage {
	if frameID >= len(l.buffer) {
		return nil
	}

	return l.buffer[frameID].GetPage()
}
