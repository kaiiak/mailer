package main

import "flag"

func main() {
	var (
		userName   = flag.String("smtp-username", "", "user name of mail")
		passWord   = flag.String("smtp-password", "", "password of mail")
		isMock     = flag.Bool("enable-mock", false, "whether mock")
		host       = flag.String("host", "127.0.0.1", "")
		localPort  = flag.String("port", "8080", "")
		endPoint   = flag.String("endpoint", "/rest/v1", "")
		configPath = flag.String("config", "./config.json", "config path")
		hostName   = "smtp.exmail.qq.com"
		serverPort = "465"
		addr       = hostName + ":" + serverPort
	)
	flag.Parse()
	if *userName == "" || *passWord == "" {
		panic("nil username or password")
	}

	//todo
	mc, err := newMailConfig(*configPath)
	if err != nil {
		panic("config error: " + err.Error())
	}

	var sender Sender
	if *isMock {
		sender = newMockSender(mc)
	} else {
		sender = newMailSender(*userName, *passWord, addr, mc)
	}
	s := newServer(*endPoint, *host+":"+*localPort, sender, mc)
	s.initRouter()
	s.run()
}
