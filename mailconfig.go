package main

import (
	"encoding/json"
	"io/ioutil"
)

//
//
//
// 配置每个邮件服务商的一些信息
// 比如 附件的大小
//

type mailConfig struct {
	SecondWait int64 `json:"secondwait"`
	SizeLimit  int64 `json:"sizelimit"`
}

func newMailConfig(path string) (*mailConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	nc := &mailConfig{}
	err = json.Unmarshal(file, nc)
	if err != nil {
		return nc, err
	}
	return nc, err
}
