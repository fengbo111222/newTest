// chart_client
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}

}
func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8083")
	CheckErr(err)
	defer conn.Close()

	//开启消息发送协程
	go MessageSend(conn)

	//主协程负责接收消息
	buf := make([]byte, 1024)

	for {
		numOfByte, err := conn.Read(buf)
		CheckErr(err)
		fmt.Println("receive server message content:" + string(buf[:numOfByte]))
	}

	fmt.Println("Client program end!")

}

func MessageSend(conn net.Conn) {
	defer conn.Close()
	var input string
	for {

		//接收系统标准输入
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		//客户端输入exit表示结束
		if strings.ToUpper(input) == "EXIT" {

			break

		}
		
		_,err:=conn.Write([]byte(input))
		if err != nil {
			fmt.Println("client connect failure: " + err.Error())
			break
		}

	}
}
