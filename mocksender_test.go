package main

import (
	"testing"
)

func TestMockSender(t *testing.T) {
	sender := newMockSender(&mailConfig{})
	mail := newBasicsMail(&mailConfig{})
	mail.from = "xxxx@yyy.com"
	mail.to = "22222@11.com"
	if err := sender.Send(mail); err != nil {
		t.Error("mock mail failed")
	}
}
