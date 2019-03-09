package main

import (
	"bytes"
)

type email interface {
	recipient() string             //收件人
	Bytes() (*bytes.Buffer, error) //信件内容
}
