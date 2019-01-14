package data

type DataClient interface {
	// 将请求信息存入数据库
	InsertRequestInfo(ri *RequestInfo) error
	// 修改请求信息--状态
	UpdateRequestInfoStatus(status int, id int64) error
	// 修改请求信息--请求次数
	UpdateRequestInfoTimes(id int64) error
	// 修改请求信息--是否发送成功过邮件
	UpdateRequestInfoSend(id int64) error
	// 查找所有异常数据（状态为：2(提交失败)和4(回滚失败)）
	ListExceptionalRequestInfo() ([]*RequestInfo, error)
	// 将成功Try的信息存入数据库
	InsertSuccessStep(s *SuccessStep) error
	BatchInsertSuccessStep(s []*SuccessStep) error
	// 更新成功Try的状态
	UpdateSuccessStepStatus(rid, sid int64, status int) error
	// 全部提交成功后，修改对应的状态（请求信息为：提交成功，Try信息状态为：提交成功）
	Confirm(id int64) error
}
