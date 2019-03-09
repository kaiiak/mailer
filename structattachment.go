package main

import (
	"mime"
	"mime/multipart"
	"net/http"
)

type attachment struct {
	fileName    string
	value       []byte
	contentType string
	inline      bool
}

type structAttachments struct {
	sizeLimit   int64
	attachments []*attachment
}

func (sa *structAttachments) attachmentExtractRequest(r *http.Request) (err error) {
	var sizeLimit int64 = 32 << 20
	var f *multipart.Form

	v := r.Header.Get("Content-Type")
	if v == "" {
		return nil
	}
	d, _, err := mime.ParseMediaType(v)
	if err != nil {
		return err
	}
	if d != "multipart/form-data" {
		return nil
	}

	//
	// 如果不判断,可能会报错;
	// multipart handled by ParseMultipartForm
	if r.MultipartForm == nil {
		mr, err := r.MultipartReader()
		if err != nil {
			return err
		}
		f, err = mr.ReadForm(sizeLimit)
		if err != nil {
			return err
		}
	} else {
		f = r.MultipartForm
	}
	for _, v := range f.File {
		for _, fh := range v {
			fileSize := int64(10 << 20)
			temp := &attachment{}
			temp.fileName = fh.Filename
			temp.contentType = fh.Header.Get("Content-Type")
			file, err := fh.Open()
			if err != nil {
				return err
			}
			if sa.sizeLimit > 0 {
				fileSize = sa.sizeLimit
			}
			temp.value, err = extractReaderSizeLimit(file, fileSize)
			if err != nil {
				return err
			}
			sa.attachments = append(sa.attachments, temp)
		}
	}
	return nil
}
