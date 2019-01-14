package leveldb

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync"
	"tcc_transaction/constant"
	"tcc_transaction/store/data"
	"time"
)

const (
	KeyRequestInfo          = "request-info/%d"
	KeyExceptionRequestInfo = "request-info/exception/%d"
	KeySuccessStep          = "success-step/%d/%d"
)

type LevelDBClient struct {
	db *leveldb.DB
}

func NewLevelDB(path string) (*LevelDBClient, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBClient{db: db}, nil
}

// 将请求信息存入数据库
func (c *LevelDBClient) InsertRequestInfo(ri *data.RequestInfo) error {
	// key 该如何设计，以方便后续的修改及查找
	if ri.Id == 0 {
		ri.Id = generateID()
	}
	key := fmt.Sprintf(KeyRequestInfo, ri.Id)
	data, err := json.Marshal(ri)
	if err != nil {
		return err
	}
	return c.db.Put([]byte(key), data, nil)
}

// 将请求信息存入数据库
func (c *LevelDBClient) insertExceptionRequestInfo(ri *data.RequestInfo) error {
	// key 该如何设计，以方便后续的修改及查找
	if ri.Id == 0 {
		return fmt.Errorf("id not exist, please check it")
	}
	key := fmt.Sprintf(KeyExceptionRequestInfo, ri.Id)
	data, err := json.Marshal(ri)
	if err != nil {
		return err
	}

	return c.db.Put([]byte(key), data, nil)
}

func (c *LevelDBClient) delExceptionRequestInfo(id int64) error {
	key := fmt.Sprintf(KeyExceptionRequestInfo, id)
	return c.db.Delete([]byte(key), nil)
}

func (c *LevelDBClient) getRequestInfo(id int64) (*data.RequestInfo, error) {
	v, err := c.db.Get([]byte(fmt.Sprintf(KeyRequestInfo, id)), nil)
	if err != nil {
		return nil, err
	}
	var ri *data.RequestInfo
	err = json.Unmarshal(v, &ri)
	if err != nil {
		return nil, err
	}
	return ri, nil
}

// 修改请求信息--状态
// 因为每条记录都是单线程操作，无需加锁
// TODO 事务一致怎么做？
func (c *LevelDBClient) UpdateRequestInfoStatus(status int, id int64) error {
	ri, err := c.getRequestInfo(id)
	if err != nil {
		return err
	}
	ri.Status = status
	err = c.InsertRequestInfo(ri)
	if err != nil {
		return err
	}
	switch status {
	case constant.RequestInfoStatus2, constant.RequestInfoStatus4:
		return c.insertExceptionRequestInfo(ri)
	case constant.RequestInfoStatus1, constant.RequestInfoStatus3:
		c.delExceptionRequestInfo(id)
	}
	return nil
}

// 修改请求信息--请求次数
func (c *LevelDBClient) UpdateRequestInfoTimes(id int64) error {
	ri, err := c.getRequestInfo(id)
	if err != nil {
		return err
	}
	ri.Times += 1
	c.InsertRequestInfo(ri)
	return c.insertExceptionRequestInfo(ri)
}

// 修改请求信息--是否发送成功过邮件
func (c *LevelDBClient) UpdateRequestInfoSend(id int64) error {
	ri, err := c.getRequestInfo(id)
	if err != nil {
		return err
	}
	ri.IsSend = 1
	ri.Status = constant.RequestInfoStatus5
	err = c.InsertRequestInfo(ri)
	if err != nil {
		return err
	}
	// 当一个完整流程走完后，删除辅助数据
	go c.delExceptionRequestInfo(id)
	return nil
}

// 查找所有异常数据（状态为：2(提交失败)和4(回滚失败)）
func (c *LevelDBClient) ListExceptionalRequestInfo() ([]*data.RequestInfo, error) {
	var ris []*data.RequestInfo
	it := c.db.NewIterator(&util.Range{Start: []byte(fmt.Sprintf(KeyExceptionRequestInfo, 1)), Limit: []byte(fmt.Sprintf(KeyExceptionRequestInfo, time.Now().Unix()*10000))}, nil)
	for it.Next() {
		v := it.Value()
		var ri *data.RequestInfo
		err := json.Unmarshal(v, &ri)
		if err != nil {
			return nil, err
		}
		it2 := c.db.NewIterator(&util.Range{Start: []byte(fmt.Sprintf(KeySuccessStep, ri.Id, 0)), Limit: []byte(fmt.Sprintf(KeySuccessStep, ri.Id, time.Now().Unix()*10000))}, nil)
		var ss []*data.SuccessStep
		for it2.Next() {
			v = it2.Value()
			var s *data.SuccessStep
			err := json.Unmarshal(v, &s)
			if err != nil {
				return nil, err
			}
			ss = append(ss, s)
		}
		ri.SuccessSteps = ss
		ris = append(ris, ri)
	}
	return ris, nil
}

// 将成功Try的信息存入数据库
func (c *LevelDBClient) InsertSuccessStep(s *data.SuccessStep) error {
	if s.Id == 0 {
		s.Id = generateID()
	}
	key := fmt.Sprintf(KeySuccessStep, s.RequestId, s.Id)
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return c.db.Put([]byte(key), data, nil)
}

func (c *LevelDBClient) BatchInsertSuccessStep(ss []*data.SuccessStep) error {
	for _, s := range ss {
		if s.Id == 0 {
			s.Id = generateID()
		}
		key := fmt.Sprintf(KeySuccessStep, s.RequestId, s.Id)
		data, err := json.Marshal(s)
		if err != nil {
			return err
		}

		err = c.db.Put([]byte(key), data, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *LevelDBClient) getSuccessStep(rid, sid int64) (*data.SuccessStep, error) {
	v, err := c.db.Get([]byte(fmt.Sprintf(KeySuccessStep, rid, sid)), nil)
	if err != nil {
		return nil, err
	}
	var ri *data.SuccessStep
	err = json.Unmarshal(v, &ri)
	if err != nil {
		return nil, err
	}
	return ri, nil
}

// 更新成功Try的状态
func (c *LevelDBClient) UpdateSuccessStepStatus(rid, sid int64, status int) error {
	s, err := c.getSuccessStep(rid, sid)
	if err != nil {
		return err
	}
	s.Status = status
	err = c.InsertSuccessStep(s)
	if err != nil {
		return err
	}
	return nil
}

// 全部提交成功后，修改对应的状态（RequestInfo状态为：提交成功，SuccessStep状态为：提交成功）
func (c *LevelDBClient) Confirm(id int64) error {
	return c.UpdateRequestInfoStatus(constant.RequestInfoStatus1, id)
}

var counter = struct {
	now int64
	idx int // 10000, 可以允许一秒需要同时生成10000条id
	sync.Mutex
}{}

func generateID() int64 {
	now := time.Now().Unix()
	if counter.now == now {
		counter.Lock()
		counter.idx += 1
		counter.Unlock()
		return now*10000 + int64(counter.idx)
	}
	counter.now = now
	counter.idx = 0
	return now * 10000
}
