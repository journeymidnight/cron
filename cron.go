package cron

import (
	"container/heap"
	"github.com/journeymidnight/util"
	"sync"
	"time"
)

type wakeup struct {
	jobName string
	at      time.Time
}

type wakeupMinHeap []wakeup

func (h wakeupMinHeap) Len() int {
	return len(h)
}

func (h wakeupMinHeap) Less(i, j int) bool {
	return h[i].at.Before(h[j].at)
}

func (h wakeupMinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *wakeupMinHeap) Push(x interface{}) {
	*h = append(*h, x.(wakeup))
}

func (h *wakeupMinHeap) Pop() interface{} {
	old := *h
	// go's Heap would put the item to be pop at the end
	// when users call the external `Pop`
	n := len(old)
	if n < 1 {
		return nil
	}
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h wakeupMinHeap) Peek() *wakeup {
	if len(h) == 0 {
		return nil
	}
	return &h[0]
}

var runner *Runner

type Runner struct {
	nexts wakeupMinHeap
	jobs  map[string]Job
	lock  *sync.Mutex // locks both `nexts` and `jobs`
}

func run() {
	ticker := time.NewTicker(time.Second)
	for {
		t := <-ticker.C

		util.WithLock(runner.lock, func() {
			for {
				next := runner.nexts.Peek()
				if next == nil {
					return
				}
				if t.Unix() < next.at.Unix() {
					return
				}
				invoke := heap.Pop(&runner.nexts).(wakeup)
				job := runner.jobs[invoke.jobName]
				if job.Trigger.OneAfterAnother {
					go func() {
						job.Callback()
						util.WithLock(runner.lock, func() {
							heap.Push(&runner.nexts, wakeup{
								jobName: job.Name,
								at:      job.Trigger.nextWakeup(time.Now()),
							})
						})
					}()
				} else {
					go job.Callback()
					heap.Push(&runner.nexts, wakeup{
						jobName: job.Name,
						at:      job.Trigger.nextWakeup(time.Now()),
					})
				}
			}
		})
	}
}

func init() {
	if runner == nil {
		runner = &Runner{
			nexts: make([]wakeup, 0, 32),
			jobs:  make(map[string]Job),
			lock:  new(sync.Mutex),
		}
		go run()
	}
}

func Register(job Job) {
	runner.lock.Lock()
	defer runner.lock.Unlock()

	unregister(job.Name) // in case duplicate

	runner.jobs[job.Name] = job
	heap.Push(&runner.nexts, wakeup{
		jobName: job.Name,
		at:      job.Trigger.nextWakeup(time.Now()),
	})
}

func unregister(jobName string) {
	if _, ok := runner.jobs[jobName]; !ok {
		return
	}
	delete(runner.jobs, jobName)
	for i := range runner.nexts {
		if runner.nexts[i].jobName == jobName {
			heap.Remove(&runner.nexts, i)
			return
		}
	}
}

func Unregister(jobName string) {
	runner.lock.Lock()
	defer runner.lock.Unlock()

	unregister(jobName)
}
