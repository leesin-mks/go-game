package socket

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

func Start()  {
	address := "192.168.0.196"
	port := 8300
	tcpAddr,err :=net.ResolveTCPAddr("tcp4",fmt.Sprintf("%s:%d",address,port))
	if err != nil {
		log.Fatalln("tcp resole error: " , err)
	}
	conn, err := net.DialTCP("tcp4",nil,tcpAddr)
	if err != nil {
		log.Fatalln("tcp error: " , err)
	}
	msg := "hello ds"
	buffer := bytes.NewBufferString(msg)
	// byteMsg := make([]byte, len(temp))

	log.Println("length: " ,len(buffer.Bytes()))
	conn.Write(buffer.Bytes())
	conn.Close()
	time.Sleep(1000)
}

func reverseBytes(i int) int {
	var k = uint32(i)
	return int((k >> 24) | ((k >> 8) & 0xFF00) | ((k << 8) & 0xFF0000) | (k << 24))
}

func writeInt(){

}

