// Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:.
package gotask

import (
	"time"
)

// Tasker interface class
type Tasker interface {

	// ID Returns the ID of the execution function
	ID() string

	// ExecuteTime Gets the next execution time
	ExecuteTime() time.Time

	// RefreshExecuteTime Refresh execution time
	RefreshExecuteTime()

	// Do Return execution function
	Do() func()
}
