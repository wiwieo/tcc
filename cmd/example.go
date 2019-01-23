package main

import (
	"fmt"
	"net/http"
	"tcc_transaction/model"
	"tcc_transaction/store/data"
)

type ExampleTcc struct {
}

func NewExampleTcc() tcc {
	return &ExampleTcc{}
}

func (d *ExampleTcc) Try(r *http.Request, api *model.RuntimeApi) ([]*data.SuccessStep, error) {
	return nil, fmt.Errorf("example is not support")
}

func (d *ExampleTcc) Confirm(r *http.Request, api *model.RuntimeApi) error {
	return fmt.Errorf("example is support")
}

func (d *ExampleTcc) Cancel(r *http.Request, api *model.RuntimeApi, nodes []*model.RuntimeTCC) ([]int64, error) {
	return nil, fmt.Errorf("example is not support")
}
