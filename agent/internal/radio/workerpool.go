package radio

import "sync"

type WorkerPool struct {
	maxWorkers int
	workerChan chan struct{}
	wg         sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		workerChan: make(chan struct{}, maxWorkers),
	}
}

func (p *WorkerPool) Submit(task func()) {
	p.wg.Add(1)

	p.workerChan <- struct{}{}

	go func() {
		defer func() {
			<-p.workerChan
			p.wg.Done()
		}()

		task()
	}()
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}