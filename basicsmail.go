package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/textproto"
	"strings"
	"time"
)

type basicsMail struct {
	from                string
	fromName            string
	to                  string
	subject             string
	html                string
	cc                  string
	bcc                 string
	replyTo             string
	plain               string
	fromer              string
	attachmentSizeLimit int64
	structAttachments
}

func newBasicsMail(config *mailConfig) *basicsMail {
	bm := &basicsMail{}
	bm.structAttachments.sizeLimit = config.SizeLimit
	return bm
}

func (bm *basicsMail) extractRequest(r *http.Request) error {
	var err error
	bm.from = r.FormValue("from")
	bm.to = strings.Join(strings.Split(r.FormValue("to"), ";"), ",")
	bm.subject = r.FormValue("subject")
	bm.html = r.FormValue("html")
	bm.fromName = r.FormValue("fromname")
	bm.cc = r.FormValue("cc")
	bm.bcc = r.FormValue("bcc")
	bm.replyTo = r.FormValue("replyto")
	bm.plain = r.FormValue("plain")
	if bm.from == "" {
		return newReturnCode(40801)
	}
	var maddr *mail.Address
	if maddr, err = mail.ParseAddress(bm.from); err == nil {
		if bm.fromName == "" {
			bm.fromer = maddr.Name + "<" + maddr.Address + ">"
		} else {
			bm.fromer = bm.fromName + "<" + maddr.Address + ">"
		}
	} else {
		return newReturnCode(40802)
	}
	if bm.to == "" {
		return newReturnCode(40862)
	}
	if bm.to, err = parseMailList(bm.to); err != nil {
		return newReturnCode(40862)
	}
	if bm.subject == "" {
		return newReturnCode(40209)
	}

	if !bm.isHTML() {
		if bm.plain == "" {
			return newReturnCode(40830)
		}
	}

	if bm.cc != "" {
		if bm.cc, err = parseMailList(bm.cc); err != nil {
			log.Println("parseMailList(bm.cc) err: ", err)
			return newReturnCode(40853)
		}
	}
	if bm.bcc != "" {
		if bm.bcc, err = parseMailList(bm.bcc); err != nil {
			log.Println("parseMailList(bm.bcc) err: ", err)
			return newReturnCode(40856)
		}
	}
	if bm.replyTo != "" {
		if bm.replyTo, err = parseMailList(bm.replyTo); err != nil {
			log.Println("parseMailList(bm.replyTo) err: ", err)
			return newReturnCode(40811)
		}
	}
	if err = bm.attachmentExtractRequest(r); err != nil {
		log.Println("bm.attachmentExtractRequest(r) err: ", err)
		switch err.(type) {
		case *returnCode:
		default:
			err = newReturnCode(40867)
		}
		return err
	}

	return nil
}

func (bm *basicsMail) recipient() string {
	var recipient string
	if len(bm.bcc) > 0 {
		recipient = bm.to + "," + bm.bcc
	} else {
		recipient = bm.to
	}
	return recipient
}

func (bm *basicsMail) isHTML() bool {
	if len(bm.html) == 0 {
		return false
	}
	return true
}

func (bm *basicsMail) Bytes() (*bytes.Buffer, error) {
	bs := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(bs)
	defer mw.Close()

	bs.WriteString("From:" + bm.fromer + "\r\n")
	bs.WriteString("To:" + bm.to + "\r\n")
	bs.WriteString("Date:" + time.Now().String() + "\r\n")
	if bm.cc != "" {
		bs.WriteString("Cc:" + bm.cc + "\r\n")
	}
	bs.WriteString("Subject:" + "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(bm.subject)) + "?=\r\n")
	if bm.replyTo != "" {
		bs.WriteString("Reply-To:" + bm.replyTo + "\r\n")
	}

	bs.WriteString("Content-Type: multipart/mixed; boundary=" + mw.Boundary() + "\r\n\r\n")
	bs.WriteString("\r\nMIME-Version: 1.0\r\n")

	//写邮件内容
	bodyHeader := make(textproto.MIMEHeader)
	if bm.isHTML() {
		bodyHeader.Add("Content-type", "text/html")
	} else {
		bodyHeader.Add("Content-type", "text/plain")
	}
	bodyHeader.Add("Content-Transfer-Encoding", "base64")
	w, err := mw.CreatePart(bodyHeader)
	if err != nil {
		return bs, err
	}
	if bm.isHTML() {
		_, err = w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(bm.html))))
	} else {
		_, err = w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(bm.plain))))
	}
	if err != nil {
		return bs, err
	}

	//添加附件
	if len(bm.attachments) > 0 {
		for _, v := range bm.attachments {
			header := make(textproto.MIMEHeader)
			if v.contentType == "" {
				header.Add("Content-Type", "application/octet-stream")
			} else {
				header.Add("Content-Type", fmt.Sprintf("application/%s", v.contentType))
			}
			header.Add("Content-Transfer-Encoding", "base64")
			header.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="=?UTF-8?B?%s?="`, base64.StdEncoding.EncodeToString([]byte(v.fileName))))
			w, err := mw.CreatePart(header)
			if err != nil {
				return bs, err
			}
			_, err = w.Write([]byte(base64.StdEncoding.EncodeToString(v.value)))
			if err != nil {
				return bs, err
			}
			bs.WriteString("\r\n")
		}
	}
	return bs, nil
}
