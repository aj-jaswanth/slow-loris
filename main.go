package main

import (
	"crypto/tls"
	"log"
	"net"
	"io"
	"sync"
	"time"
)

var address = "localhost:8888"

func main() {
	waitGroup := &sync.WaitGroup{}
	for i := 0; i < 400; i++ {
		provider := provider{pool: []string{"GET ", "/blahblahbla ", "HTTP/1.1", "\r\n", "Allow"}, delay: 0, pos: 0}
		worker := worker{waitGroup: waitGroup, address: address, provider: provider}
		waitGroup.Add(1)
		go worker.work()
	}
	waitGroup.Wait()
}

type worker struct {
	waitGroup *sync.WaitGroup
	address   string
	writer    io.Writer
	provider  provider
}

type provider struct {
	pool  []string
	delay time.Duration
	pos   int
}

func (p *provider) next() []byte {
	defer func() {
		p.delay = p.delay + 2
		p.pos = p.pos + 1
	}()
	time.Sleep(p.delay * time.Second)
	if p.pos >= len(p.pool) {
		return nil
	}
	return []byte(p.pool[p.pos])
}

func (worker *worker) work() {
	defer worker.waitGroup.Done()
	worker.writer = getHttpConnection(address)
	for {
		bytes := worker.provider.next()
		if bytes == nil {
			return
		}
		_, err := worker.writer.Write(bytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func getHttpConnection(address string) net.Conn {
	conn, e := net.Dial("tcp", address)
	if e != nil {
		log.Println("no conn", e)
	}
	return conn
}

func getHttpsConnection(address string) *tls.Conn {
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
