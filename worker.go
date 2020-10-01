package worker

import "fmt"

type JobResponse struct {
	err  error
	name string
	res  int
	url  string
}

type Job interface {
	Name() string
	Callback() JobResponse
}

func Do(jobs []Job, maxWorkers int) {
    jobsPool := make(chan Job, len(jobs))
	resultsPool := make(chan JobResponse, len(jobs))

	for i := 0; i < maxWorkers; i++ {
		go worker(jobsPool, resultsPool)
	}

	makeJobs(jobsPool, jobs)
	getResults(resultsPool, jobs)
}

func worker(jobs <-chan Job, response chan<- JobResponse) {
	for n := range jobs {
		response <- n.Callback()
	}
}

func makeJobs(jobs chan<- Job, queue []Job) {
	for _, t := range queue {
		jobs <- t
	}
}

func getResults(response <-chan JobResponse, queue []Job) {
	for range queue {
		job := <-response
		status := fmt.Sprintf("[result] '%s' to '%s' was fetched with status '%d'\n", job.name, job.url, job.res)
		if job.err != nil {
			status = fmt.Sprintf(job.err.Error())
		}
		fmt.Printf(status)
	}
}
