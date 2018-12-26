# 使用MMAP写的日志
MMAP：[MMAP基本概念](https://www.cnblogs.com/huxiao-tee/p/4660352.html)

```
win10（64bit 8G） linux子系统 
普通日志100W次写入需要耗时：18.01s
MMAP日志100W次定稿需要耗时：1.679s
```

# 问题
因为mmap的特性，导致服务一旦crash，则日志文件会遗留大量的占位符
目前不支持windows系统，使用一般的文件写入方式

# 注意事项
* 一、直接使用写日志时，默认使用终端打印日志
* 二、如果需要使用mmap来写日志，则需要调用文件`exported.go`中的`SetPath()`方法，指定日志文件
* 三、日志打印级别默认为`info`，如果需要指定其他级别，则需要调用文件`exported.go`中的`SetLevel()`方法，重新设定日志级别
* 四、具体可调用方法，参照：[日志方法](./exported.go)