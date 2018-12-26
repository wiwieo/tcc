// Package gotask Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:.
package gotask

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Tasks 任务列表
type Tasks []Tasker

func (s Tasks) Len() int      { return len(s) }
func (s Tasks) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Tasks) Less(i, j int) bool {
	return s[i].ExecuteTime().Before(s[j].ExecuteTime())
}

type taskList struct {
	// 所有任务列表
	taskers Tasks
}

type intervalChange struct {
	task     *Task
	interval time.Duration
}

var (
	tasks *taskList
	editC = make(chan interface{})
	stopC = make(chan string)
	wg    = &sync.WaitGroup{}
)

func init() {
	tasks = &taskList{}
	go doAllTask()
}

// AddToTaskList add the task to the execution list
func AddToTaskList(ts ...Tasker) {
	for _, t := range ts {
		if t == nil {
			continue
		}
		wg.Add(1)
		editC <- t
	}
}

func (tl *taskList) addToTaskList(t Tasker) {
	tl.taskers = append(tl.taskers, t)
	wg.Done()
}

// Stop stop corresponding tasks through the id of task
func Stop(id string) {
	stopC <- id
}

func (tl *taskList) stop(id string) {
	for k, v := range tl.taskers {
		if v.ID() == id {
			tl.taskers = append(tl.taskers[:k], tl.taskers[k+1:]...)
		}
	}
}

// ChangeInterval changes the interval between the tasks specified by the ID,
// Apply only to polling tasks.
func ChangeInterval(id string, interval time.Duration) error {
	tsk := tasks.get(id)
	if tsk == nil {
		wg.Wait()
		tsk = tasks.get(id)
		if tsk == nil {
			return fmt.Errorf("Task does not exist")
		}
	}
	var task *Task
	var ok bool
	if task, ok = tsk.(*Task); !ok {
		return fmt.Errorf("This type does not support modifying the execution interval")
	}
	editC <- &intervalChange{
		task:     task,
		interval: interval,
	}

	return nil
}

func changeInterval(i *intervalChange) {
	i.task.SetInterval(i.interval)
}

func (tl *taskList) get(id string) Tasker {
	for _, v := range tl.taskers {
		if v.ID() == id {
			return v
		}
	}
	return nil
}

func doAllTask() {
	var timer *time.Timer

	var now time.Time
	for {
		sort.Sort(tasks.taskers)

		now = time.Now()

		if len(tasks.taskers) == 0 {
			timer = time.NewTimer(time.Hour * 100000)
		} else {
			sub := tasks.taskers[0].ExecuteTime().Sub(now)
			if sub < 0 {
				sub = 0
			}
			timer = time.NewTimer(sub)
		}

		for {
			select {
			case now = <-timer.C:
				doNestedTask()
			case edit := <-editC:
				now = time.Now()
				timer.Stop()
				if t, ok := edit.(Tasker); ok {
					tasks.addToTaskList(t)
				} else if ic, ok := edit.(*intervalChange); ok {
					changeInterval(ic)
				}
			case id := <-stopC:
				tasks.stop(id)
			}
			break
		}
	}
}

func doNestedTask() {
	for _, v := range tasks.taskers {
		if v.ExecuteTime().Before(time.Now()) {
			go v.Do()()
			v.RefreshExecuteTime()
		} else {
			return
		}
	}
}
