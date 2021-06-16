package cron

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestWakeupMinHeap(t *testing.T) {
	h := wakeupMinHeap(make([]wakeup, 0, 32))
	heap.Push(&h, wakeup{at: time.Unix(1, 0)})
	heap.Push(&h, wakeup{at: time.Unix(3, 0)})
	heap.Push(&h, wakeup{at: time.Unix(5, 0)})
	heap.Push(&h, wakeup{at: time.Unix(7, 0)})
	heap.Push(&h, wakeup{at: time.Unix(2, 0)})
	heap.Push(&h, wakeup{at: time.Unix(4, 0)})
	heap.Push(&h, wakeup{at: time.Unix(6, 0)})
	heap.Push(&h, wakeup{at: time.Unix(8, 0)})

	last := 0
	for {
		peek := h.Peek()
		if peek == nil {
			break
		}
		pop := heap.Pop(&h).(wakeup)
		fmt.Println(pop.at.Second())
		if pop.at.Second() < last {
			t.Error("bad heap", pop.at.Second())
		}
		last = pop.at.Second()
	}
}

func TestRunner(t *testing.T) {
	job1Count := 0
	Register(Job{
		Name: "job1",
		Trigger: Trigger{
			Every: Period{
				Every: 1,
				Unit:  Second,
			},
		},
		Callback: func() {
			fmt.Println("job1", time.Now())
			job1Count += 1
		},
	})
	job2Count := 0
	Register(Job{
		Name: "job2",
		Trigger: Trigger{
			Every: Period{
				Every: 2,
				Unit:  Second,
			},
		},
		Callback: func() {
			fmt.Println("job2", time.Now())
			job2Count += 1
		},
	})
	job3Count := 0
	Register(Job{
		Name: "job3",
		Trigger: Trigger{
			Every: Period{
				Every: 3,
				Unit:  Second,
			},
		},
		Callback: func() {
			fmt.Println("job3", time.Now())
			job3Count += 1
		},
	})
	time.Sleep(10 * time.Second)
	if job1Count < 9 || job1Count > 11 ||
		job2Count < 4 || job2Count > 6 ||
		job3Count < 2 || job3Count > 4 {
		t.Error("bad runner", job1Count, job2Count, job3Count)
	}
	fmt.Println("unregister...")
	Unregister("job1")
	Unregister("job2")
	time.Sleep(6 * time.Second)
	if job1Count < 9 || job1Count > 11 ||
		job2Count < 4 || job2Count > 6 ||
		job3Count < 4 || job3Count > 6 {
		t.Error("bad runner", job1Count, job2Count, job3Count)
	}
	fmt.Println(job1Count, job2Count, job3Count)

	// clean up
	Unregister("job3")
	fmt.Println(runner.nexts)
}

func TestCronOneAfterAnother(t *testing.T) {
	jobCount := 0
	Register(Job{
		Name: "job1",
		Trigger: Trigger{
			Every: Period{
				Every: 1,
				Unit:  Second,
			},
			OneAfterAnother: true,
		},
		Callback: func() {
			time.Sleep(3 * time.Second)
			fmt.Println("job1", time.Now())
			jobCount += 1
		},
	})
	time.Sleep(10 * time.Second)
	if jobCount > 3 {
		t.Error("triggered too many times")
	}
}
