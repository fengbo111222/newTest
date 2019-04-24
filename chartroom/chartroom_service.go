// chartroom_service
package main

import (
	"fmt"
	"net"
	"strings"
)

// 消息队列，带缓存的buf
var mesgque = make(chan string, 1000)

var quitChan = make(chan bool)

//存储客户端链接映射
var onlinesConns = make(map[string]net.Conn)

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

	go ConsumeMessage()

	for {
		conn, err := lis_socket.Accept()
		CheckErr(err)

		//将conn存到online映射

		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onlinesConns[addr] = conn
		for i := range onlinesConns {

			fmt.Println(i)

		}
		go ProcessInfo(conn)

		//如果有客户端链接，则打开一个协程处理
		go ProcessInfo(conn)

	}

}

//消费者协程
func ConsumeMessage() {

	for {
		select {
		case mesg := <-mesgque:
			//对消息解析
			doProcessMessage(mesg)
		case <-quitChan:
			break
		}

	}

}

//消息解析函数
func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMesg := contents[1]
		addr = strings.Trim(addr, " ")
		//通过map查看是否有客户端，没有则不能发送
		if conn, ok := onlinesConns[addr]; ok {

			_, err := conn.Write([]byte(sendMesg))
			if err != nil {
				fmt.Println("online conns send failure!")
			}

		}
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
			mesgque <- string(buf[:numOfByte])

		}
	}

}
