package main

//Sender 邮件的发送者
type Sender interface {
	Send(email) error
}
