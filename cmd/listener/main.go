package main

import (
	"fmt"
	"http/internal/request"
	"log"
	"net"
)



func main(){
	listner,err:=net.Listen("tcp",":42067")
	if err!=nil{
		log.Fatal(err)
	}
	for{
		conn,err:=listner.Accept()
		if err!=nil{
			log.Fatal(err)
			break
		}
		r,err:=request.RequestFromReader(conn)
		if err!=nil{
			log.Fatal(err)
			break
		}
		fmt.Println("Request line:")
		fmt.Println("- Method: ", r.RequestLine.Method)
		fmt.Println("- Target: ", r.RequestLine.RequestTarget)
		fmt.Println("- Version: ", r.RequestLine.HttpVersion)
	}
}