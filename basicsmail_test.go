package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNewBasicsmailWithForm(t *testing.T) {
	bs := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(bs)

	from := `Robot<aNxFi37X@outlook.com>`
	fromname := `来自测试发送`
	subject := `测试`
	html := `<p>这是一封测试邮件</p>`
	to := `aNxFi37X@outlook.com`
	params := map[string]string{
		"from":     from,
		"fromname": fromname,
		"subject":  subject,
		"html":     html,
		"to":       to,
	}

	for k, v := range params {
		if err := mw.WriteField(k, v); err != nil {
			t.Error("wirte ", k, " err: ", err)
		}
	}

	mw.Close()

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/rest/v1/mail/send", bs)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	bm := newBasicsMail(&mailConfig{})
	if err = bm.extractRequest(req); err != nil {
		t.Error(err)
	}
	if bm.from != from {
		t.Error("wrong from")
	}
	if bm.fromName != fromname {
		t.Error("wrong fromname")
	}
	if bm.subject != subject {
		t.Error("wrong subject")
	}
	if bm.html != html {
		t.Error("wrong html")
	}
	if bm.to != to {
		t.Error("wrong to")
	}
	if bm.isHTML() == false {
		t.Error("wrong type")
	}
	bf, err := bm.Bytes()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(bf.String(), "attachments") {
		t.Error("wrong mail Value")
	}
}

func TestNewBasicsmailWithQuery(t *testing.T) {
	from := `Robot<aNxFi37X@outlook.com>`
	fromname := `来自测试发送`
	subject := `测试`
	plain := `这是一封测试邮件`
	to := `aNxFi37X@outlook.com`
	mail := url.Values{
		"from":     {from},
		"fromname": {fromname},
		"subject":  {subject},
		"plain":    {plain},
		"to":       {to},
	}
	req, err := http.NewRequest(http.MethodPost, "?"+mail.Encode(), nil)
	if err != nil {
		t.Error(err)
	}
	bm := newBasicsMail(&mailConfig{})
	if err = bm.extractRequest(req); err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	if bm.from != from {
		t.Error("wrong from")
	}
	if bm.fromName != fromname {
		t.Error("wrong fromname")
	}
	if bm.subject != subject {
		t.Error("wrong subject")
	}
	if bm.plain != plain {
		t.Error("wrong plain")
	}
	if bm.to != to {
		t.Error("wrong to")
	}
	if bm.isHTML() == true {
		t.Error("wrong type")
	}

	bf, err := bm.Bytes()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(bf.String(), "attachments") {
		t.Error("wrong mail Value")
	}
}
