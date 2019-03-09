package main

import (
	"log"
	"net/http"
)

type server struct {
	endPoint string
	sender   Sender
	addr     string
	s        *http.ServeMux
	config   *mailConfig
}

func newServer(endPoint, addr string, sender Sender, config *mailConfig) *server {
	return &server{
		endPoint: endPoint,
		sender:   sender,
		addr:     addr,
		s:        &http.ServeMux{},
		config:   config,
	}
}

func (s *server) initRouter() {
	s.s.HandleFunc(s.endPoint+"/"+"mail/send", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if er := recover(); er != nil {
				log.Println("painc: ", er)
				w.Write(newReturnCode(49906).Bytes())
				return
			}
		}()
		if r.Method == http.MethodPost || r.Method == http.MethodGet {
			bm := newBasicsMail(s.config)
			err := bm.extractRequest(r)
			if err != nil {
				w.Write([]byte(err.Error()+"\n"))
				return
			}

			if err = s.sender.Send(bm); err != nil {
				log.Println("send mail err: ", err)
				w.Write(newReturnCode(40901).Bytes())
				return
			}

			w.Write(newReturnCode(200).Bytes())
		} else {
			log.Println("MethodNotAllowed: ", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	})
}

//Run run server
func (s *server) run() {
	srv := http.Server{
		Addr:    s.addr,
		Handler: s.s,
	}
	srv.ListenAndServe()
}
