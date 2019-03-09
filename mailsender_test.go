package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestMailSender(t *testing.T) {
	sender := newMailSender("", "", "", &mailConfig{})
	mail := newBasicsMail(&mailConfig{})
	mail.from = "xxx@yyy.com"
	if err := sender.Send(mail); err == nil {
		t.Error("no expect value")
	}
}

func TestNewCoon(t *testing.T) {
	ln := newLocalListener(t)
	defer ln.Close()

	var done = make(chan struct{})
	go func() {
		defer close(done)
		c, err := ln.Accept()
		if err != nil {
			t.Error(err)
		}
		go connServerHandle(c, t)
	}()

	conn := newConn(ln.Addr().String())
	if conn == nil {
		t.Error("new conn failed")
	}
	if _, err := conn.Write([]byte("something wrong?")); err != nil {
		t.Error(err)
	}

	fmt.Println("read begin")

	b, err := extractReader(conn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("read done")
	<-done
	fmt.Println("Done")

	fmt.Println(string(b))
	if fmt.Sprintf(string(b)) != "hello" {
		t.Fatal("not expect value")
	}
}

func connServerHandle(conn net.Conn, t *testing.T) {
	var err error
	// buf := make([]byte, 1024)
	// reqLen, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// 	t.Error(err)
	// }
	// for i := 0; i < reqLen; i++ {
	// 	if buf[i] > 126 {
	// 		conn.Close()
	// 	}
	// }
	// fmt.Println("recived:\n", string(buf[:reqLen]))
	if _, err = conn.Write([]byte("hello")); err != nil {
		t.Error("connServerHandle error: " + err.Error())
	}
	conn.Close()
}

// type Addr interface {
//         Network() string // name of the network (for example, "tcp", "udp")
//         String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
// }

type testAddr struct {
}

func (a *testAddr) Network() string {
	return ""
}

func (a *testAddr) String() string {
	return "fakehost"
}

type faker struct {
	io.ReadWriter
}

func (f faker) Close() error        { return nil }
func (f faker) LocalAddr() net.Addr { return nil }
func (f faker) RemoteAddr() net.Addr {
	return &testAddr{}
}
func (f faker) SetDeadline(time.Time) error      { return nil }
func (f faker) SetReadDeadline(time.Time) error  { return nil }
func (f faker) SetWriteDeadline(time.Time) error { return nil }

var newClientServer = `220 hello world
250-mx.google.com at your service
250-SIZE 35651584
250-AUTH LOGIN PLAIN
250 8BITMIME
235 Accepted
250 Ok
221 OK
`

var newClientClient = `EHLO fakehost
AUTH PLAIN AAA=
MAIL FROM:<> BODY=8BITMIME
QUIT
`

func TestNewClient(t *testing.T) {
	server := strings.Join(strings.Split(newClientServer, "\n"), "\r\n")
	client := strings.Join(strings.Split(newClientClient, "\n"), "\r\n")

	var cmdbuf bytes.Buffer
	bcmdbuf := bufio.NewWriter(&cmdbuf)
	var fake faker
	fake.ReadWriter = bufio.NewReadWriter(bufio.NewReader(strings.NewReader(server)), bcmdbuf)

	out := func() string {
		bcmdbuf.Flush()
		return cmdbuf.String()
	}

	c, err := newClient("", "", fake)
	if err != nil {
		t.Fatalf("newClient: %v\n(after %v)", err, out())
	}
	if err := c.Quit(); err != nil {
		t.Fatalf("QUIT failed: %s", err)
	}

	actualcmds := out()
	if client != actualcmds {
		t.Fatalf("Got:\n%s\nExpected:\n%s", actualcmds, client)
	}

}

func newLocalListener(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		ln, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// var sendMailServer = `220 hello world
// 250-mx.google.com at your service
// 250-SIZE 35651584
// 250-AUTH LOGIN PLAIN
// 250 8BITMIME
// 235 Accepted
// 250 Ok
// 221 OK
// `

// func TestSend(t *testing.T) {
// 	ln := newLocalListener(t)

// 	var cmdbuff bytes.Buffer
// 	bcmdbuff := bufio.NewWriter(&cmdbuff)

// 	var done = make(chan struct{})
// 	go func(data []string) {
// 		defer close(done)
// 		conn, err := ln.Accept()
// 		defer conn.Close()
// 		fmt.Println("Accepted")
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		tc := textproto.NewConn(conn)
// 		for i := 0; i < len(data) && data[i] != ""; i++ {
// 			fmt.Println(data[i])
// 			tc.PrintfLine(data[i])
// 			// tc.PrintfLine("221 Goodbye")
// 			for len(data[i]) >= 4 && data[i][3] == '-' {
// 				i++
// 				tc.PrintfLine(data[i])
// 			}
// 			read := false
// 			if !read || data[i] == "354 Go ahead" {
// 				msg, err := tc.ReadLine()
// 				bcmdbuff.Write([]byte(msg + "\r\n"))
// 				read = true
// 				if err != nil {
// 					t.Errorf("Read error: %v", err)
// 					return
// 				}
// 				if data[i] == "354 Go ahead" && msg == "." {
// 					break
// 				}
// 			}
// 		}
// 	}(strings.Split(sendMailServer, "\n"))

// 	ms := newMailSender("", "", ln.Addr().String(), &mailConfig{})
// 	bm := newBasicsMail(&mailConfig{})
// 	bm.to = "xxx@yy.com"
// 	if err := ms.Send(bm); err != nil {
// 		t.Error(err)
// 	}

// 	<-done
// 	bcmdbuff.Flush()
// 	fmt.Println("cmd:\n", cmdbuff.String())
// }

var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBjjCCATigAwIBAgIQMon9v0s3pDFXvAMnPgelpzANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJB
AM0u/mNXKkhAzNsFkwKZPSpC4lZZaePQ55IyaJv3ovMM2smvthnlqaUfVKVmz7FF
wLP9csX6vGtvkZg1uWAtvfkCAwEAAaNoMGYwDgYDVR0PAQH/BAQDAgKkMBMGA1Ud
JQQMMAoGCCsGAQUFBwMBMA8GA1UdEwEB/wQFMAMBAf8wLgYDVR0RBCcwJYILZXhh
bXBsZS5jb22HBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAEwDQYJKoZIhvcNAQELBQAD
QQBOZsFVC7IwX+qibmSbt2IPHkUgXhfbq0a9MYhD6tHcj4gbDcTXh4kZCbgHCz22
gfSj2/G2wxzopoISVDucuncj
-----END CERTIFICATE-----`)

// localhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAM0u/mNXKkhAzNsFkwKZPSpC4lZZaePQ55IyaJv3ovMM2smvthnl
qaUfVKVmz7FFwLP9csX6vGtvkZg1uWAtvfkCAwEAAQJART2qkxODLUbQ2siSx7m2
rmBLyR/7X+nLe8aPDrMOxj3heDNl4YlaAYLexbcY8d7VDfCRBKYoAOP0UCP1Vhuf
UQIhAO6PEI55K3SpNIdc2k5f0xz+9rodJCYzu51EwWX7r8ufAiEA3C9EkLiU2NuK
3L3DHCN5IlUSN1Nr/lw8NIt50Yorj2cCIQCDw1VbvCV6bDLtSSXzAA51B4ZzScE7
sHtB5EYF9Dwm9QIhAJuCquuH4mDzVjUntXjXOQPdj7sRqVGCNWdrJwOukat7AiAy
LXLEwb77DIPoI5ZuaXQC+MnyyJj1ExC9RFcGz+bexA==
-----END RSA PRIVATE KEY-----`)
