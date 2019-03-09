package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestAttachmentExtractRequest(t *testing.T) {
	sa := &structAttachments{}

	bs := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(bs)
	w, err := mw.CreateFormFile("README.md", "README.md")
	if err != nil {
		t.Error(err)
	}
	file := []byte("README.md")
	_, err = w.Write(file)
	if err != nil {
		t.Error(err)
	}

	w, err = mw.CreateFormFile("CHANGELOG.md", "CHANGELOG.md")
	if err != nil {
		t.Error(err)
	}
	file = []byte("CHANGELOG.md")
	_, err = w.Write(file)
	if err != nil {
		t.Error(err)
	}

	mw.Close()

	req, err := http.NewRequest(http.MethodPost, "", bs)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	err = sa.attachmentExtractRequest(req)
	if err != nil {
		t.Error(err)
	}
	if len(sa.attachments) != 2 {
		t.Log("lenght attachments: ", len(sa.attachments))
		t.Error("wrong count of attachments")
	}
	if sa.attachments[0].fileName != "README.md" {
		t.Error("wrong fileName")
	}
	if len(sa.attachments[1].value) != len(file) {
		t.Error("wrong file value")
	}
}
