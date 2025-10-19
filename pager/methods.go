package pager

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/frame"
	"github.com/dark-vinci/nildb/interfaces"
)

func (p *Pager) GetNewPage(pin bool) (*interfaces.RepPage, base.PageNumber, error) {
	pn := p.AllocatePage()
	page, err := p.GetPage(pn, pin)
	if err != nil {
		return nil, 0, err
	}
	p.cache.MarkDirty(pn)
	// Initialize the page if needed
	return page, pn, nil
}

func (p *Pager) AllocatePage() base.PageNumber {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.freePages.Len() > 0 {
		return heap.Pop(&p.freePages).(base.PageNumber)
	}

	pn := p.nextPageNumber
	p.nextPageNumber++

	return pn
}

func (p *Pager) ReleasePage(pn base.PageNumber) {
	p.cache.Unpin(pn)
}

func (p *Pager) Stop() {
	p.worker.Stop()
}

// GetPage retrieves a page from cache or disk
func (p *Pager) GetPage(pn base.PageNumber, pin bool) (*frame.Frame, error) {
	frameID := p.cache.GetFrameID(pn)

	// PAGE IS IN CACHE
	if frameID != nil {
		if pin {
			p.cache.Pin(pn)
		}

		page := p.cache.GetFrame(*frameID).(*frame.Frame)

		return page, nil
	}

	// Cache miss, need to evict if full and load
	//if p.cache.Size() >= p.cache.MaxSize() {
	//	victimID := p.cache.findVictim()
	//	victimFrame := p.cache.GetFrame(victimID).(frame.Frame)
	//
	//	if victimFrame.IsSet(constants.DirtyFlag) {
	//		resultChan := p.worker.Write(victimFrame.PageNumber, victimFrame.Page)
	//
	//		result := <-resultChan
	//		if result.Error != nil {
	//			return nil, result.Error
	//		}
	//
	//		p.cache.MarkClean(victimFrame.PageNumber)
	//	}
	//}

	// EVICT IF NEEDED
	frameID2 := p.cache.Map(pn)

	// Load from disk
	frame := p.cache.GetFrame(frameID2)

	resultChan := p.worker.Read(pn, frame.Page)
	result := <-resultChan

	if result.Error != nil {
		if os.IsNotExist(result.Error) || fmt.Sprint(result.Error).Contains("not found") {
			// Initialize empty page
			// Assume interfaces.RepPage has a New method or something
			frame.Page.FromBytes(make([]byte, p.cache.PageSize)) // or initialize
		} else {
			return nil, result.Error
		}
	}

	if pin {
		p.cache.Pin(pn)
	}

	return &frame.Page, nil
}
