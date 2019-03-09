package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type mailSender struct {
	userName string
	password string
	server   string
	config   *mailConfig
}

func newMailSender(user, passwd, addr string, config *mailConfig) *mailSender {
	return &mailSender{
		userName: user,
		password: passwd,
		server:   addr,
		config:   config,
	}
}

func (sender *mailSender) Send(m email) error {
	message := make(chan error)
	var er error
	go func() {
		if _, err := parseMailList(m.recipient()); err != nil {
			message <- err
			return
		}
		client, err := newClient(sender.userName, sender.password, newConn(sender.server))
		if err != nil {
			message <- err
			return
		}
		al, err := mail.ParseAddressList(m.recipient())
		if err != nil {
			message <- err
			return
		}
		for _, v := range al {
			if err = client.Rcpt(v.Address); err != nil {
				message <- err
				return
			}
		}
		wc, err := client.Data()
		if err != nil {
			message <- err
			return
		}
		defer wc.Close()
		text, err := m.Bytes()
		if err != nil {
			message <- err
			return
		}
		_, err = wc.Write(text.Bytes())
		if err != nil {
			message <- err
			return
		}
		err = client.Quit()
		if strings.HasPrefix(strings.TrimSpace(err.Error()), "250") {
			message <- client.Close()
			return
		}
		message <- err
		return
	}()
	select {
	case er = <-message:
	case <-time.After(time.Second * 10):
		er = fmt.Errorf("send mail timeout")
	}
	return er
}

//newClient 返回一个登陆后的客户端，登陆失败就会返回error
// 不使用时，手动调用 Quit()
func newClient(user, passwd string, conn net.Conn) (*smtp.Client, error) {
	if conn == nil {
		return nil, fmt.Errorf("nil connect")
	}
	if conn.RemoteAddr() == nil {
		return nil, fmt.Errorf("function RemoteAddr is nil")
	}
	addr := conn.RemoteAddr().String()
	c, err := smtp.NewClient(conn, addr)
	if err != nil {
		return c, err
	}
	err = c.Hello(addr)
	if err != nil {
		fmt.Println("hello err: ", err)
		return c, err
	}
	auth := smtp.PlainAuth("", user, passwd, addr)
	err = c.Auth(auth)
	if err != nil {
		fmt.Println("auth err: ", err)
		return c, err
	}
	if err := c.Mail(user); err != nil {
		fmt.Println("mail err: ", err)
		return c, err
	}
	return c, nil
}

func newConn(addr string) net.Conn {
	var conn net.Conn
	var err error
	conn, err = tls.Dial("tcp", addr, &tls.Config{})
	if err != nil {
		// qq邮箱的TLS端口用 net.Dial 无响应且无报错
		// fix it
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return nil
		}
	}
	return conn
}
