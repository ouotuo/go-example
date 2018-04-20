package main

import (
	"net"
	"log"
	"flag"
	"bufio"
	"io"
)


func main() {
	addr:=flag.String("addr","127.0.0.1:8888","listen addr")
	protocol:=flag.String("p","tcp","protocol")

	flag.Parse()

	log.Println("listen ",*protocol,*addr)

	if *protocol=="udp"{
		udpLoop(*addr)
	}else{
		otherLoop(*protocol,*addr)
	}
}

func udpLoop(addr string){
	udpAddr,err:=net.ResolveUDPAddr("udp",addr)
	if err!=nil{
		log.Fatal(err)
	}

	lc,err:=net.ListenUDP("udp",udpAddr)
	if err!=nil{
		log.Fatal(err)
	}


	buf := make([]byte,10000)
	var n int
	var ad *net.UDPAddr
	for{
		n,ad,err=lc.ReadFromUDP(buf)
		if err!=nil{
			log.Fatal(err)
		}
		_,err=lc.WriteTo(buf[:n],ad)
		if err!=nil{
			log.Println(err)
		}
	}
}

func otherLoop(protocol,addr string){
	ln, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Printf("accept conn error,%v",err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	log.Println("open conn",conn.RemoteAddr())
	defer func(){
		conn.Close()
		log.Println("close conn",conn.RemoteAddr())
	}()

	//读字节流，然后返回
	var reader=bufio.NewReader(conn)

	for {
		slice, err := reader.ReadSlice('\n')
		if err != nil {
			if err!=io.EOF{
				log.Println("read error",err)
			}
			break
		} else {
			_, err = conn.Write(slice); if err != nil {
				log.Println("write error",err)
				break
			}
		}
	}
}
