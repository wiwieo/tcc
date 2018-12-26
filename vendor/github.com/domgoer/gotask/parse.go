// Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:.
package gotask

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	errParseTime = errors.New("parse time error")
)

var defaultValue = time.Time{}

type dayParse struct {
}

func newDayParse() Parser {
	return &dayParse{}
}

// Parse 接收格式 "hh:mm:ss",返回begintime
func (p *dayParse) Parse(s string) (time.Time, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 3 {
		return defaultValue, errParseTime
	}
	var hour, minute, second int
	var err error
	if hour, err = strconv.Atoi(ss[0]); err != nil {
		return defaultValue, errParseTime
	}
	if minute, err = strconv.Atoi(ss[1]); err != nil {
		return defaultValue, errParseTime
	}
	if second, err = strconv.Atoi(ss[2]); err != nil {
		return defaultValue, errParseTime
	}
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())
	return t, nil
}

type monthParse struct {
}

func newMonthParse() Parser {
	return &monthParse{}
}

// Parse 接收格式 dd mm:hh:ss   dd为每月几号，如果需要每月最后一天 dd=-1
func (p *monthParse) Parse(s string) (time.Time, error) {
	s2 := strings.Split(s, " ")
	if len(s2) != 2 {
		return defaultValue, errParseTime
	}
	ss := strings.Split(s2[1], ":")
	if len(ss) != 3 {
		return defaultValue, errParseTime
	}
	var day int
	var hour, minute, second int
	var err error
	if day, err = strconv.Atoi(s2[0]); err != nil {
		return defaultValue, errParseTime
	}
	if hour, err = strconv.Atoi(ss[0]); err != nil {
		return defaultValue, errParseTime
	}
	if minute, err = strconv.Atoi(ss[1]); err != nil {
		return defaultValue, errParseTime
	}
	if second, err = strconv.Atoi(ss[2]); err != nil {
		return defaultValue, errParseTime
	}
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), day, hour, minute, second, 0, now.Location())
	// 排除部分月没有31号
	if t.Month() != now.Month() {
		t.AddDate(0, 1, 0)
	}
	return t, nil
}
