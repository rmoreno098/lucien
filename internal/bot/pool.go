package bot

import "log"

type Task interface {
	Execute() error
}

type WorkerPool struct {
	Tasks   chan Task
	handler func(Task) error
	workers int
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		go func(id int) {
			for task := range p.Tasks {
				if task == nil {
					log.Printf("Worker %v stopped", id)
					return
				}
				if err := p.handler(task); err != nil {
					log.Printf("An error occurred trying to execute task: %v", err)
				}
			}
		}(i)
	}
}

func (p *WorkerPool) Execute(task Task) {
	p.Tasks <- task
}
