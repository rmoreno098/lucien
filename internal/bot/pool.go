package bot

import "log"

type Task interface {
	Execute() error
}

type WorkerPool struct {
	tasks   chan Task
	workers int
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		go func(id int) {
			for task := range p.tasks {
				if task == nil {
					log.Printf("Worker %v stopped", id)
					return
				}
				if err := task.Execute(); err != nil {
					log.Printf("An error occurred trying to execute task: %v", err)
				}
			}
		}(i)
	}
}

func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}
