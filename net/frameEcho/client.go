package main

import (
	"github.com/ouotuo/go-example/net/frameEcho/frame"
	"flag"
	"time"
	"log"
)



func main() {
	addr:=flag.String("addr","127.0.0.1:8888","echo server addr")
	num:=flag.Int("n",100,"echo num")
	print:=flag.Bool("print",false,"print res or not")
	flag.Parse()

	log.Println("connect to",*addr)

	var client frame.FrameConn
	var err error
	client=frame.NewClientSocketFrameConn(*addr)

	err=client.Open()
	if err != nil {
		panic(err)
	}

	var word="hello"
	var t1,t2 time.Time

	t1=time.Now()
	var toRecv []byte

	for i:=0;i<*num;i++{
		toRecv, err = client.WriteReadFrame([]byte(word))
		if err != nil {
			panic(err)
		}

		if *print{
			log.Println(len(toRecv),string(toRecv))
		}
	}
	t2=time.Now()
	log.Println("num",*num)
	log.Println("costTime",t2.Sub(t1))

	client.Close()
}

