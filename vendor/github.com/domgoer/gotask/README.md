# gotask

### 轮询任务

```go
package main

import (
    "time"
    
    "github.com/domgoer/gotask"
)

func main()  {
     tk := gotask.NewTask(time.Second*20,func() {
            // do ... 
     })
     gotask.AddToTaskList(tk)
}
```

### 定时任务

```go
package main

import (
    "github.com/domgoer/gotask"
)

func main()  {
     tkDay,_ := gotask.NewDayTask("12:20:00",func() {
            // do ... 
     })
     tkMonth,_ := gotask.NewMonthTask("20 12:20:00",func() {
             // do ... 
      })
     gotask.AddToTaskList(tkDay)
     gotask.AddToTaskList(tkMonth)
}
```

> 多任务

```go
package main

import (
    "github.com/domgoer/gotask"
)

func main()  {
     tkDays,_ := gotask.NewDayTasks([]string{"12:20:00","10:10:10"},func() {
            // do ... 
     })
     tkMonths,_ := gotask.NewMonthTasks([]string{"20 12:20:00","21 10:10:10"},func() {
             // do ... 
      })
     gotask.AddToTaskList(tkDays...)
     gotask.AddToTaskList(tkMonths...)
}
```

### 停止

```go
package main

import (
    "github.com/domgoer/gotask"
)

func main()  {
     gotask.Stop("task.ID()")
}
```
