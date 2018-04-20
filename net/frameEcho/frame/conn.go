package frame

import (
	"time"
	"fmt"
	"net"
)

type FrameConn interface{
	Open()(error)
	Close()(error)
	WriteFrame([]byte)(error)
	ReadFrame()([]byte,error)
	WriteReadFrame([]byte)([]byte,error)
	IsOpen()(bool)
}

const(
	FRAME_CONN_STATE_DISCONNECT=0
	FRAME_CONN_STATE_CONNECTING=1
	FRAME_CONN_STATE_CONNECTED=2

	FRAME_CONN_DEFAULT_DIALTIMEOUT=time.Second*5
	FRAME_CONN_DEFAULT_WRITETIMEOUT=time.Second*15
	FRAME_CONN_DEFAULT_READTIMEOUT=time.Second*20

	FRAME_CONN_BUF_LEN=1024*2

	SOCKET_FRAME_CONN_STYPE_CLIENT=0
	SOCKET_FRAME_CONN_STYPE_SERVER=1
)

type SocketFrameConn struct{
	stype int  //socket类型
	state int
	conn net.Conn
	dialTimeout time.Duration
	readTimeout time.Duration
	writeTimeout time.Duration
	addr string
	network string
	frameParser *FrameParser
	buf []byte
}

func NewClientSocketFrameConn(addr string)(*SocketFrameConn){
	conn:=&SocketFrameConn{
		stype:SOCKET_FRAME_CONN_STYPE_CLIENT,
		state:FRAME_CONN_STATE_DISCONNECT,
		dialTimeout:FRAME_CONN_DEFAULT_DIALTIMEOUT,
		readTimeout:FRAME_CONN_DEFAULT_READTIMEOUT,
		writeTimeout:FRAME_CONN_DEFAULT_WRITETIMEOUT,
		network:"tcp",
		frameParser:NewFrameParser(),
		buf:make([]byte,FRAME_CONN_BUF_LEN),
		addr:addr,
	}
	return conn
}

func NewServerSocketFrameConn(conn net.Conn)(*SocketFrameConn){
	fconn:=&SocketFrameConn{
		stype:SOCKET_FRAME_CONN_STYPE_SERVER,
		state:FRAME_CONN_STATE_CONNECTED,
		dialTimeout:FRAME_CONN_DEFAULT_DIALTIMEOUT,
		readTimeout:FRAME_CONN_DEFAULT_READTIMEOUT,
		writeTimeout:FRAME_CONN_DEFAULT_WRITETIMEOUT,
		frameParser:NewFrameParser(),
		buf:make([]byte,FRAME_CONN_BUF_LEN),
		addr:"",
		network:"tcp",
		conn:conn,
	}
	return fconn
}

func (self *SocketFrameConn)SetNetwork(network string)(*SocketFrameConn){
	self.network=network
	return self
}
func (self *SocketFrameConn)SetDialTimeout(timeout time.Duration)(*SocketFrameConn){
	self.dialTimeout=timeout
	return self
}
func (self *SocketFrameConn)SetWriteTimeout(timeout time.Duration)(*SocketFrameConn){
	self.writeTimeout=timeout
	return self
}
func (self *SocketFrameConn)SetReadTimeout(timeout time.Duration)(*SocketFrameConn){
	self.readTimeout=timeout
	return self
}

func (self *SocketFrameConn)IsOpen()(isOpen bool){
	isOpen=self.state==FRAME_CONN_STATE_CONNECTED
	return
}

func (self *SocketFrameConn)Open()(err error){
	if self.state==FRAME_CONN_STATE_CONNECTED{
		return
	}
	if self.stype!=SOCKET_FRAME_CONN_STYPE_CLIENT{
		err=fmt.Errorf("only client socket can open")
		return
	}
	if self.addr==""{
		err=fmt.Errorf("addr is empty")
		return
	}
	if self.state==FRAME_CONN_STATE_CONNECTING{
		err=fmt.Errorf("client is connecting")
		return
	}

	self.state=FRAME_CONN_STATE_CONNECTING

	self.conn, err = net.DialTimeout(self.network,self.addr,self.dialTimeout)
	if err==nil{
		self.state=FRAME_CONN_STATE_CONNECTED
		self.frameParser.Reset()
	}else{
		self.state=FRAME_CONN_STATE_DISCONNECT
	}
	return
}

func (self *SocketFrameConn)Close()(err error){
	if self.state==FRAME_CONN_STATE_CONNECTED{
		return
	}
	if self.conn!=nil{
		err=self.conn.Close()
		self.conn=nil
	}
	return
}

func (self *SocketFrameConn)WriteFrame(data []byte)(err error){
	if self.state!=FRAME_CONN_STATE_CONNECTED{
		err=fmt.Errorf("state is not connected,state=%d",self.state)
		return
	}

	//set timeout
	err=self.conn.SetWriteDeadline(time.Now().Add(self.writeTimeout));if err!=nil{
		//err=fmt.Errorf("setWriteDeadline error,%v",err)
		return
	}

	//send
	_,err=self.conn.Write(GetFrameMagicHeadBytes(len(data)));if err!=nil{
		//err=fmt.Errorf("write frameMagicHead error,%v",err)
		return
	}
	_,err=self.conn.Write(data);if err!=nil{
		//err=fmt.Errorf("write data error,%v",err)
		return
	}

	return
}

func (self *SocketFrameConn)ReadFrame()(data []byte,err error){
	if self.state!=FRAME_CONN_STATE_CONNECTED{
		err=fmt.Errorf("state is not connected,state=%d",self.state)
		return
	}

	//set timeout
	err=self.conn.SetReadDeadline(time.Now().Add(self.readTimeout));if err!=nil{
		//err=fmt.Errorf("setReadDeadline error,%v",err)
		return
	}

	//recv
	var n int=0
	var frames [][]byte
	for {
		n,err=self.conn.Read(self.buf)
		if err!=nil{
			//err=fmt.Errorf("read error,%v",err)
			return
		}

		frames,err=self.frameParser.Parse(self.buf[:n])
		if err!=nil{
			err=fmt.Errorf("frameParser parse error,%v",err)
			return
		}
		if len(frames)>0{
			data=frames[0]
			return
		}
	}

	return
}

func (self *SocketFrameConn)WriteReadFrame(toWrite []byte)(toRead []byte,err error){
	err=self.WriteFrame(toWrite);if err!=nil{
		return
	}
	toRead,err=self.ReadFrame()
	return
}


