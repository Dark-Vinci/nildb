package diskscheduler

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
)

func (w *DiskWorker) Read(pageNumber base.PageNumber, page faces.PageHandle) chan faces.DiskResult {
	resultChan := make(chan faces.DiskResult, 1)

	w.queue <- faces.DiskRequest{
		Type:       base.ReadOp,
		PageNumber: pageNumber,
		Page:       page,
		ResultChan: resultChan,
	}

	return resultChan
}

func (w *DiskWorker) Write(pageNumber base.PageNumber, page faces.PageHandle) chan faces.DiskResult {
	resultChan := make(chan faces.DiskResult, 1)

	w.queue <- faces.DiskRequest{
		Type:       base.WriteOp,
		PageNumber: pageNumber,
		Page:       page,
		ResultChan: resultChan,
	}

	return resultChan
}

func (w *DiskWorker) Stop() {
	close(w.stopChan)

	w.waitGroup.Wait()
}
