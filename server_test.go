package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	sender := newMockSender(&mailConfig{})
	server := newServer("/rest/v1", "127.0.0.1:8080", sender, &mailConfig{})
	server.initRouter()
	go server.run()
	bs := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(bs)

	from := `Robot<aNxFi37X@outlook.com>`
	fromname := `来自测试发送`
	subject := `测试`
	html := `<p>这是一封测试邮件</p>`
	to := `aNxFi37X@outlook.com`
	cc := `aNxFi37X@outlook.com`
	bcc := `aNxFi37X@outlook.com`
	replyto := `aNxFi37X@outlook.com`
	params := map[string]string{
		"from":     from,
		"fromname": fromname,
		"subject":  subject,
		"html":     html,
		"to":       to,
		"cc":       cc,
		"bcc":      bcc,
		"replyto":  replyto,
	}

	for k, v := range params {
		if err := mw.WriteField(k, v); err != nil {
			t.Error("wirte ", k, " err: ", err)
		}
	}

	w, err := mw.CreateFormFile("README.md", "README.md")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadFile("README.md")
	if err != nil {
		t.Error(err)
	}
	_, err = w.Write(file)
	if err != nil {
		t.Error(err)
	}

	w, err = mw.CreateFormFile("CHANGELOG.md", "CHANGELOG.md")
	if err != nil {
		t.Error(err)
	}
	file, err = ioutil.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Error(err)
	}
	_, err = w.Write(file)
	if err != nil {
		t.Error(err)
	}

	mw.Close()

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/rest/v1/mail/send", bs)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{}
	reps, err := client.Do(req)
	defer reps.Body.Close()

	if err != nil {
		t.Error(err)
	}
	v, err := extractReader(reps.Body)
	if err != nil {
		t.Error(err)
	}
	result := &sendCloudV1{}
	err = json.Unmarshal(v, result)
	if err != nil {
		t.Error(err)
	}
	if result.Message != "success" {
		t.Error("wrong result")
	}
}

func TestGetMethod(t *testing.T) {
	sender := newMockSender(&mailConfig{})
	server := newServer("/rest/v1", "127.0.0.1:8080", sender, &mailConfig{})
	server.initRouter()
	go server.run()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/rest/v1/mail/send", nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	reps, err := client.Do(req)
	defer reps.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if reps.StatusCode != 200 {
		t.Error("bad reponse")
	}
}
