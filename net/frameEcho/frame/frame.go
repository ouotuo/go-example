package frame

import (
)
import (
	"encoding/binary"
	"fmt"
)

//结构magic|head|body
const(
	FRAME_MAGIC_INDEX=0
	FRAME_MAGIC_LEN=1
	FRAME_MAGIC_BYTE=99
	FRAME_HEAD_INDEX=FRAME_MAGIC_LEN
	FRAME_HEAD_LEN=4

	FRAME_BODY_INDEX=FRAME_MAGIC_LEN+FRAME_HEAD_LEN

	STATE_FRAME_PARSER_MAGIC=0    //读magic
	STATE_FRAME_PARSER_HEAD=1    //读head
	STATE_FRAME_PARSER_BODY=2  //读body
)

func GetFrameHeadBytes(len int)[]byte{
	bs:=make([]byte,FRAME_HEAD_LEN,FRAME_HEAD_LEN)
	binary.BigEndian.PutUint32(bs,uint32(len))
	return bs
}

func GetFrameMagicHeadBytes(len int)[]byte{
	bs:=make([]byte,FRAME_MAGIC_LEN+FRAME_HEAD_LEN,FRAME_MAGIC_LEN+FRAME_HEAD_LEN)
	bs[FRAME_MAGIC_INDEX]=FRAME_MAGIC_BYTE
	binary.BigEndian.PutUint32(bs[FRAME_HEAD_INDEX:],uint32(len))
	return bs
}

type FrameParser struct{
	state int   //状态
	headBytes []byte  //长度字节
	headBytesSize int  //长度字节数量

	bodyLen int
	bodyBytes []byte
	bodyBytesSize int //数据长度
}

func NewFrameParser()*FrameParser{
	parser:=&FrameParser{state:STATE_FRAME_PARSER_MAGIC,headBytes:make([]byte,FRAME_HEAD_LEN,FRAME_HEAD_LEN),headBytesSize:0}
	return parser
}

//消费数据流，产生frame
func (self *FrameParser)Parse(buf []byte)(frames [][]byte,err error){
	frames=make([][]byte,0,1)
	var copySize int

	for ;len(buf)>0 && err==nil; {
		switch self.state{
		case STATE_FRAME_PARSER_MAGIC:
			if buf[0]==FRAME_MAGIC_BYTE {
				//ok进入下一个状态
				self.state=STATE_FRAME_PARSER_HEAD
				buf=buf[FRAME_MAGIC_LEN:]
				self.headBytesSize=0
			}else{
				err=fmt.Errorf("frame magic byte not ok,should %d,get %d",FRAME_MAGIC_BYTE,buf[0])
				return
			}
		case STATE_FRAME_PARSER_HEAD:
			//状态为解释长度
			copySize = copy(self.headBytes[self.headBytesSize:], buf)
			self.headBytesSize += copySize
			buf = buf[copySize:]
			if self.headBytesSize == FRAME_HEAD_LEN {
				//长度字节已经全部读到，转为数字
				self.bodyLen = int(binary.BigEndian.Uint32(self.headBytes))
				//test frameLen

				self.state = STATE_FRAME_PARSER_BODY
				self.bodyBytes = make([]byte, self.bodyLen, self.bodyLen)
				self.bodyBytesSize = 0
			}

		case STATE_FRAME_PARSER_BODY:
			//状态为读数据
			copySize = copy(self.bodyBytes[self.bodyBytesSize:], buf)
			self.bodyBytesSize += copySize
			buf = buf[copySize:]

			if self.bodyBytesSize == self.bodyLen {
				//成功解释新frame
				frames = append(frames, self.bodyBytes)
				self.state = STATE_FRAME_PARSER_MAGIC
			}
		default:
			err=fmt.Errorf("unknown state %d",self.state)
		}
	}
	return
}

//重置
func (self *FrameParser)Reset(){
	self.state=STATE_FRAME_PARSER_MAGIC
}

