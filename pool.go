//
// dnsdiff
// 2014 @oseemann
//

// Simple Work Queue implementation
// with a fixed number of workers

package main

type Job struct {
	run      func(string)
	arg      string
	running  bool
	finished bool
}

type WorkQueue struct {
	joblist []Job
	size    int
	jobs    chan int
	result  chan int
}

func NewWorkQueue(size int) *WorkQueue {
	wq := new(WorkQueue)
	wq.size = size
	wq.jobs = make(chan int, size)
	wq.result = make(chan int, 100)
	wq.joblist = make([]Job, 0, 100)
	return wq
}

func (wq *WorkQueue) Add(job Job) {
	wq.joblist = append(wq.joblist, job)
}

func (wq *WorkQueue) Run() {
	// start workers
	for w := 1; w <= wq.size; w++ {
		go worker(wq.joblist, wq.jobs, wq.result)
	}

	// kick off initial jobs
	for i := range wq.joblist {
		wq.jobs <- i
	}
}

func (wq *WorkQueue) WaitAll() {

	for {
		allfinished := true
		for _, job := range wq.joblist {
			if job.finished == false {
				allfinished = false
			}

			if job.running {
				// at least one job is running
				// wait to get a result from the chan
				i := <-wq.result
				// i does not necessarily refer to the same job
				wq.joblist[i].finished = true
				wq.joblist[i].running = false
			}
		}

		if allfinished {
			break
		}
	}
}

func worker(joblist []Job, jobs chan int, result chan int) {

	for {
		// wait for incoming new job ids
		job_id := <-jobs
		joblist[job_id].running = true
		joblist[job_id].run(joblist[job_id].arg)
		result <- job_id
	}
}

// vim: set filetype=go ts=4 sw=4 sts=4 noexpandtab:
