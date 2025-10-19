package diskscheduler

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
)

func (w *DiskWorker) Read(pageNumber base.PageNumber, page interfaces.PageHandle) chan interfaces.DiskResult {
	resultChan := make(chan interfaces.DiskResult, 1)

	w.queue <- interfaces.DiskRequest{
		Type:       base.ReadOp,
		PageNumber: pageNumber,
		Page:       page,
		ResultChan: resultChan,
	}

	return resultChan
}

func (w *DiskWorker) Write(pageNumber base.PageNumber, page interfaces.PageHandle) chan interfaces.DiskResult {
	resultChan := make(chan interfaces.DiskResult, 1)

	w.queue <- interfaces.DiskRequest{
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
