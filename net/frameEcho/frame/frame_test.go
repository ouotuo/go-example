package frame

import (
	"testing"
	"bytes"
)

var frame4Byte=[]byte{FRAME_MAGIC_BYTE,0,0,0,4,9,9,8,3}

var frameParser=NewFrameParser()


func Test_oneFrame(t *testing.T) {
	frameParser.Reset()
	frames,err:=frameParser.Parse(frame4Byte)
	if err!=nil{
		t.Errorf("get error,%v",err)
		return
	}

	if len(frames)!=1{
		t.Errorf("frame number wrong,%d",len(frames))
		return
	}
	if bytes.Compare(frames[0],frame4Byte[FRAME_BODY_INDEX:])!=0{
		t.Errorf("frame byte not same")
		return
	}
}

func Test_oneFrameHalf(t *testing.T) {
	var datas =make([]byte,0,len(frame4Byte)*2)
	datas=append(datas,frame4Byte...)
	datas=append(datas,frame4Byte...)
	var halfIndex=len(frame4Byte)+len(frame4Byte)/2

	frameParser.Reset()
	frames,err:=frameParser.Parse(datas[0:halfIndex])
	if err!=nil{
		t.Errorf("get error,%v",err)
		return
	}

	if len(frames)!=1{
		t.Errorf("frame number wrong,%d",len(frames))
		return
	}
	if bytes.Compare(frames[0],frame4Byte[FRAME_BODY_INDEX:])!=0{
		t.Errorf("frame byte not same")
		return
	}

	frames,err=frameParser.Parse(datas[halfIndex:])
	if err!=nil{
		t.Errorf("get error,%v",err)
		return
	}

	if len(frames)!=1{
		t.Errorf("frame number wrong,%d",len(frames))
		return
	}
	if bytes.Compare(frames[0],frame4Byte[FRAME_BODY_INDEX:])!=0{
		t.Errorf("frame byte not same")
		return
	}
}


func strDataFrameTest(strs []string,t *testing.T){
	var datas=make([]byte,0,1000)
	for _,str:=range strs{
		datas=append(datas,FRAME_MAGIC_BYTE)
		var strBytes=[]byte(str)
		datas=append(datas,GetFrameHeadBytes(len(strBytes))...)
		datas=append(datas,strBytes...)
	}
	frameParser.Reset()
	frames,err:=frameParser.Parse(datas)
	if err!=nil{
		t.Errorf("parse error,%v",err)
		return
	}
	if len(frames)!=len(strs){
		t.Errorf("frame len not same,should %d,but get %d",len(strs),len(frames))
		return
	}
	//比较字符串
	for i,str:=range strs{
		var getStr=string(frames[i])
		if str!=getStr{
			t.Errorf("frame data not same,should %s,but get %s",str,getStr)
		}
	}
}

func Test_str(t *testing.T) {
	var strs=[]string{"abc"}
	strDataFrameTest(strs,t)

	strs=[]string{"9772343","88234234","9876234324","ffabcwe"}
	strDataFrameTest(strs,t)

	strs=[]string{"9sdfsbb1234772343","jjj88234234","9876234fsfdsf324","ffabcw2e","fsjdjfjsdj","dsfsfsdfdsf"}
	strDataFrameTest(strs,t)

	strs=[]string{"9sdfsbb1234772343","","","ffabcw2e","fsjdjfjsdj","dsfsfsdfdsf"}
	strDataFrameTest(strs,t)
}

