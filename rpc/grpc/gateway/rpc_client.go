
package main

import (
	"log"
	"time"
	"flag"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/ouotuo/go-example/rpc/grpc/gateway/proto"
)


func main() {
	addr:=flag.String("addr","127.0.0.1:9090","echo server addr")
	num:=flag.Int("n",100,"echo num")
	print:=flag.Bool("print",false,"print res or not")

	flag.Parse()

	// Set up a connection to the server.
	var t1,t2 time.Time
	t1=time.Now()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	t2=time.Now()
	log.Println("dial time",t2.Sub(t1))

	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)


	t1=time.Now()
	var name="abc"
	for i:=0;i<*num;i++{
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		r, err := c.Echo(context.Background(), &pb.StringMessage{Value: name})
		if err != nil {
			log.Fatal(err)
		}
		//cancel()
		if *print{
			log.Println("res",r.Value)
		}
	}
	t2=time.Now()
	log.Println("num",*num)
	log.Println("costTime",t2.Sub(t1))
}
