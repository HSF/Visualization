// Copyright 2017 The vizmon-demo Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	var (
		addr = flag.String("addr", ":8080", "[ip]:port for web server")
		freq = flag.Duration("freq", 100*time.Millisecond, "data polling interval")
	)

	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("vizmon-demo: ")

	log.Printf("starting up web-server on: %v\n", *addr)
	srv, err := newServer(*addr, *freq)
	if err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}

	http.Handle("/", srv)
	http.Handle("/data", websocket.Handler(srv.dataHandler))
	http.HandleFunc("/echo", srv.wrap(srv.echoHandler))

	err = http.ListenAndServe(srv.addr, nil)
	if err != nil {
		log.Fatalf("error running server: %v\n", err)
	}
}

type server struct {
	addr string
	freq time.Duration

	tmpl    *template.Template
	dataReg registry // clients interested in monitoring data

	data  chan MonData
	plots chan Plots
	echo  chan MonData
}

func newServer(addr string, freq time.Duration) (*server, error) {
	if addr == "" {
		addr = getHostIP() + ":80"
	}

	srv := &server{
		addr:    addr,
		freq:    freq,
		dataReg: newRegistry(),
		tmpl:    template.Must(template.New("vizmon").Parse(indexTmpl)),
		data:    make(chan MonData),
		plots:   make(chan Plots),
		echo:    make(chan MonData),
	}
	go srv.run()

	return srv, nil
}

func (srv *server) Freq() float64 {
	return 1 / srv.freq.Seconds()
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.wrap(srv.rootHandler)(w, r)
}

func (srv *server) rootHandler(w http.ResponseWriter, r *http.Request) error {
	return srv.tmpl.Execute(w, srv)
}

func (srv *server) wrap(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (srv *server) echoHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid HTTP request (got=%v, want=%v)", r.Method, http.MethodGet)
	}
	timeout := time.NewTimer(2 * srv.freq)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		return fmt.Errorf("timeout retrieving data from board")
	case data := <-srv.echo:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (srv *server) dataHandler(ws *websocket.Conn) {
	log.Printf("new client...")
	c := &client{
		srv:   srv,
		reg:   &srv.dataReg,
		datac: make(chan []byte, 256),
		ws:    ws,
	}
	c.reg.register <- c
	defer c.Release()

	c.run()
}

func (srv *server) run() {
	go srv.daq()
	go srv.mon()
	for {
		select {
		case c := <-srv.dataReg.register:
			log.Printf("client registering [%v]...", c.ws.LocalAddr())
			srv.dataReg.clients[c] = true

		case c := <-srv.dataReg.unregister:
			if _, ok := srv.dataReg.clients[c]; ok {
				delete(srv.dataReg.clients, c)
				close(c.datac)
				log.Printf(
					"client disconnected [%v]\n",
					c.ws.LocalAddr(),
				)
			}

		case plots := <-srv.plots:
			if len(srv.dataReg.clients) == 0 {
				// no client connected
				continue
			}
			dataBuf := new(bytes.Buffer)
			err := json.NewEncoder(dataBuf).Encode(&plots)
			if err != nil {
				log.Printf("error marshalling data: %v\n", err)
				continue
			}
			for c := range srv.dataReg.clients {
				select {
				case c.datac <- dataBuf.Bytes():
				default:
					close(c.datac)
					delete(srv.dataReg.clients, c)
				}
			}
		}
	}
}

func (srv *server) daq() {
	tick := time.NewTicker(srv.freq)
	defer tick.Stop()

	i := 0
	for range tick.C {
		data, err := srv.fetchData()
		if err != nil {
			log.Printf("error fetching data: %v\n", err)
			continue
		}

		i++
		if i%10 == 0 {
			log.Printf("daq: %+v\n", data)
		}
		select {
		case srv.data <- data:
		default:
			// nobody is listening
			// drop it on the floor
		}
	}
}

func (srv *server) mon() {
	table := make([]MonData, 0, 1024)
	var data MonData
	for {
		select {
		case data = <-srv.data:
			if len(table) == cap(table) {
				i := len(table)
				copy(table[:i/2], table[i/2:])
				table = table[:i/2]
			}
			table = append(table, data)
			ctrl, err := newControlPlots(table)
			if err != nil {
				log.Printf("error creating monitoring plots: %v", err)
				continue
			}
			ps := Plots{
				update: time.Now().UTC(),
				plots:  ctrl,
				data:   data,
			}
			select {
			case srv.plots <- ps:
			default:
				// nobody is listening
			}
		case srv.echo <- data:
		}
	}
}

func (srv *server) fetchData() (MonData, error) {
	var data MonData
	for i, val := range []MonValue{
		{"vtx-x", 0},
		{"vtx-y", 0},
		{"vtx-z", 0},
		{"ele-tof", 0},
	} {
		switch i {
		case 0, 2:
			val.Value = rand.ExpFloat64() * float64(i+1) * rand.NormFloat64()
		case 1:
			val.Value = rand.NormFloat64() + float64(i+1)*5
		case 3:
			switch v := rand.Float64(); v < 0.5 {
			case true:
				val.Value = rand.NormFloat64() + 17.5
			case false:
				val.Value = rand.NormFloat64() + 22.5
			}
		}
		data.Values = append(data.Values, val)
	}
	return data, nil
}

func getHostIP() string {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("could not retrieve hostname: %v\n", err)
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Fatalf("could not lookup hostname IP: %v\n", err)
	}

	for _, addr := range addrs {
		ipv4 := addr.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4.String()
	}

	log.Fatalf("could not infer host IP")
	return ""
}
