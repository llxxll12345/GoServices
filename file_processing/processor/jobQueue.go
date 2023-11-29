package processor

import (
	"sync"
)

var (
	jobQueue       chan *ProcessJob
	processedFiles = make(map[string]string) // Map to store processed file paths by token
	mutex          sync.Mutex
	wg             sync.WaitGroup
)

type ProcessJob struct {
	FileToken string
	FilePath  string
}

func StartJobs() {
	jobQueue = make(chan *ProcessJob, 1000) // Adjust the buffer size as needed
	for i := 0; i < 10; i++ {
		go worker()
	}
}

func GetProcessedFilePath(token string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	path, ok := processedFiles[token]
	return path, ok
}

func AddJob(job *ProcessJob) {
	wg.Add(1)
	jobQueue <- job
}

func worker() {
	for {
		select {
		case job := <-jobQueue:
			processFile(job)
			wg.Done()
		}
	}
}
