// chartroom_service
package main

import (
	"fmt"
	"net"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}

}
func main() {
	// 监听socket本地8080端口
	lis_socket, err := net.Listen("tcp", "127.0.0.1:8083")
	CheckErr(err)
	defer lis_socket.Close()

	fmt.Println("Service is waiting......")
	for {
		conn, err := lis_socket.Accept()
		CheckErr(err)

		//如果有客户端链接，则打开一个协程处理
		go ProcessInfo(conn)

	}

}

// 负责接收协程

func ProcessInfo(conn net.Conn) {

	buf := make([]byte, 1024)
	defer conn.Close()

	for {

		numOfByte, err := conn.Read(buf)
		if err != nil {
			break

		}

		// 消息为不为0说明有消息
		if numOfByte != 0 {
			//同时返回链接方地址
			remoteAdd := conn.RemoteAddr()
			fmt.Print(remoteAdd)
			fmt.Printf("Has received this message: %s\n", string(buf[:numOfByte]))
		}
	}

}
