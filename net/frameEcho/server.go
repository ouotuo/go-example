package main

import (
	"net"
	"log"
	"io"
	"github.com/ouotuo/go-example/net/frameEcho/frame"
	"flag"
)

func main() {
	addr:=flag.String("addr","127.0.0.1:8888","listen addr")
	log.Println("server listen",*addr)

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	var socketFrameConn=frame.NewServerSocketFrameConn(conn)
	var id =conn.RemoteAddr().String()
	log.Printf("[%s],new conn",id)
	var err error
	defer func(){
		if err==nil {
			log.Printf("[%s],close conn", id)
		}else{
			log.Printf("[%s],close conn,with error %v", id,err)
		}
	}()

	//set timeout
	var data []byte
	for {
		data,err=socketFrameConn.ReadFrame()

		if err==nil && data!=nil{
			log.Printf("[%s] echo %s", id,string(data))
			//write frame
			err=socketFrameConn.WriteFrame(data)
		}

		//close socket
		if err != nil {
			socketFrameConn.Close()
			if err==io.EOF{
				err=nil
			}
			return
		}


	}

	//should not reach
	return
}

