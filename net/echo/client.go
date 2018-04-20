package main

import (
	"net"
	"log"
	"time"
	"bufio"
	"flag"
	"io"
)


func main() {
	addr:=flag.String("addr","127.0.0.1:8888","echo server addr")
	num:=flag.Int("n",100,"echo num")
	protocol:=flag.String("p","tcp","protocol")
	print:=flag.Bool("print",false,"print res or not")
	flag.Parse()

	log.Println("connect to ",*protocol,*addr)

	var e echo

	if *protocol=="udp"{
		e=newUdpEcho(*addr)
	}else{
		e=newOtherEcho(*protocol,*addr)
	}

	var word="hello"
	var err error
	var t1,t2 time.Time


	t1=time.Now()
	var res string
	for i:=0;i<*num;i++{
		res,err=e.echo(word)
		if err!=nil{
			log.Fatal(err)
		}
		if *print{
			log.Println(len(res),res)
		}
	}
	t2=time.Now()
	log.Println("num",*num)
	log.Println("costTime",t2.Sub(t1))

	e.close()
}


type echo interface {
	echo(string)(string,error)
	close()
}

type otherEcho struct{
	conn net.Conn
	reader *bufio.Reader
}

func newOtherEcho(protocol string,addr string)(e echo){
	bt:=time.Now()
	conn, err := net.DialTimeout(protocol,addr,time.Second*time.Duration(5))
	if err != nil {
		log.Fatal(err)
	}
	et:=time.Now()
	log.Println("echo dial costTime",et.Sub(bt))

	e=&otherEcho{
		conn:conn,
		reader:bufio.NewReader(conn),
	}
	return
}

func(e *otherEcho)echo(word string)(res string,err error){
	_,err=e.conn.Write([]byte(word));if err!=nil{
		log.Fatal(err)
	}
	_,err=e.conn.Write([]byte{'\n'});if err!=nil{
		log.Fatal(err)
	}

	bs,err:=e.reader.ReadSlice('\n')
	res=string(bs[0:len(bs)-1])

	if err!=nil{
		if err!=io.EOF{
			log.Fatal(err)
		}
	}else if res==""{
		log.Fatal("res is empty")
	}
	return
}

func(e *otherEcho)close(){
	e.conn.Close()
}

type udpEcho struct{
	conn *net.UDPConn
	buf []byte
}

func newUdpEcho(addr string)(e echo){
	bt:=time.Now()
	udpAddr,err:=net.ResolveUDPAddr("udp",addr)
	if err!=nil{
		log.Fatal(err)
	}

	conn,err:=net.DialUDP("udp",nil,udpAddr)
	if err!=nil{
		log.Fatal(err)
	}

	et:=time.Now()
	log.Println("echo dial costTime",et.Sub(bt))

	e=&udpEcho{
		conn:conn,
		buf:make([]byte,1500,1500),
	}
	return
}

func(e *udpEcho)echo(word string)(res string,err error){
	_,err=e.conn.Write([]byte(word))
	if err!=nil{
		log.Fatal(err)
	}

	n,err:=e.conn.Read(e.buf)
	if err!=nil{
		log.Fatal(err)
	}

	res=string(e.buf[:n])

	return
}

func(e *udpEcho)close(){
	e.conn.Close()
}