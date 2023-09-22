package main

import "log"

func createWorkerPool() {
	TaskQueue = make(chan Task, 100)
	ErrQueue = make(chan error, 100)
	wg.Add(NumberOfWorkers)

	for i := 0; i < NumberOfWorkers; i++ {
		go func() {
			defer wg.Done()
			for task := range TaskQueue {
				err := createTask(task)
				if err != nil {
					select {
					case ErrQueue <- err:
						log.Printf("Error creating task: %v", err)
					default:
						log.Printf("Error queue is full: %v", err)
					}
				}
			}
		}()
	}
}
