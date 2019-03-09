package main

import (
	"encoding/json"
	"testing"
)

func TestNewReturnCode(t *testing.T) {
	c := newReturnCode(200)
	if c.code() != 200 {
		t.Error("wrong code")
	}
	result := &sendCloudV1{}
	err := json.Unmarshal([]byte(c.Error()), result)
	if err != nil {
		t.Error(err)
	}
	if result.Message != "success" {
		t.Error("wrong message")
	}
}

func TestBytes(t *testing.T) {
	c := newReturnCode(200)
	if c.code() != 200 {
		t.Error("wrong code")
	}
	result := &sendCloudV1{}
	err := json.Unmarshal(c.Bytes(), result)
	if err != nil {
		t.Error(err)
	}
	if result.Message != "success" {
		t.Error("wrong message")
	}
}

func TestNoExistCode(t *testing.T) {
	c := newReturnCode(11111)
	if c.code() != 99999 {
		t.Error("wrong code")
	}
}
