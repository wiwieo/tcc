// Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:.
package gotask

import (
	"github.com/satori/go.uuid"
	"time"
)

// Task Polling tasks
type Task struct {
	id string

	executeTime time.Time

	interval time.Duration

	do func()
}

// NewTask create a new polling task
func NewTask(t time.Duration, do func()) Tasker {
	uid := uuid.NewV4()
	return &Task{
		id:          uid.String(),
		do:          do,
		interval:    t,
		executeTime: time.Now().Add(t),
	}
}

// ExecuteTime gets the next execution time
func (t *Task) ExecuteTime() time.Time {
	return t.executeTime
}

// SetInterval modify execution interval
func (t *Task) SetInterval(td time.Duration) {
	t.interval = td
	t.changeExecuteTime(td)
}

func (t *Task) changeExecuteTime(td time.Duration) {
	t.executeTime = time.Now().Add(td)
}

// RefreshExecuteTime refresh execution interval
func (t *Task) RefreshExecuteTime() {
	t.executeTime = t.executeTime.Add(t.interval)
}

// ID return taskID
func (t *Task) ID() string {
	return t.id
}

// Do return Task Function
func (t *Task) Do() func() {
	return t.do
}

// DayTask 日任务
type DayTask struct {
	id string

	executeTime time.Time

	do func()
}

// NewDayTask create a new daily task
func NewDayTask(tm string, do func()) (Tasker, error) {
	uid := uuid.NewV4()
	pt := newTimeParser(dayParseType)
	begin, err := pt.Parse(tm)
	if err != nil {
		return nil, err
	}
	if begin.Before(time.Now()) {
		begin = begin.Add(time.Hour * 24)
	}
	return &DayTask{
		id:          uid.String(),
		do:          do,
		executeTime: begin,
	}, nil
}

// NewDayTasks create new daily tasks
func NewDayTasks(tms []string, do func()) ([]Tasker, error) {
	var ts []Tasker
	for _, tm := range tms {
		dt, err := NewDayTask(tm, do)
		if err != nil {
			return nil, err
		}
		ts = append(ts, dt)
	}
	return ts, nil
}

func (d *DayTask) ID() string {
	return d.id
}

func (d *DayTask) ExecuteTime() time.Time {
	return d.executeTime
}

func (d *DayTask) RefreshExecuteTime() {
	d.executeTime = d.executeTime.Add(time.Hour * 24)
}

func (d *DayTask) Do() func() {
	return d.do
}

// MonthTask create monthly task
type MonthTask struct {
	id string

	executeTime time.Time

	do func()
}

// NewMonthTask initialize a function that executes each month
func NewMonthTask(tm string, do func()) (Tasker, error) {
	uid := uuid.NewV4()
	pt := newTimeParser(monthParseType)
	begin, err := pt.Parse(tm)
	if err != nil {
		return nil, err
	}
	if begin.Before(time.Now()) {
		begin = begin.AddDate(0, 1, 0)
	}
	return &MonthTask{
		id:          uid.String(),
		do:          do,
		executeTime: begin,
	}, nil
}

// NewMonthTasks initialize a function that executes each month
func NewMonthTasks(tms []string, do func()) ([]Tasker, error) {
	var ts []Tasker
	for _, tm := range tms {
		mt, err := NewMonthTask(tm, do)
		if err != nil {
			return nil, err
		}
		ts = append(ts, mt)
	}
	return ts, nil
}

func (m *MonthTask) ID() string {
	return m.id
}

func (m *MonthTask) ExecuteTime() time.Time {
	return m.executeTime
}

func (m *MonthTask) RefreshExecuteTime() {
	m.executeTime = m.executeTime.AddDate(0, 1, 0)
}

func (m *MonthTask) Do() func() {
	return m.do
}
