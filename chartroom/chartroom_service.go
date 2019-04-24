// chartroom_service
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	LOG_DIRECTORY = "./test.log"
)

// 消息队列，带缓存的buf
var mesgque = make(chan string, 1000)

var quitChan = make(chan bool)

var logger *log.Logger
var logFile *os.File

//存储客户端链接映射
var onlinesConns = make(map[string]net.Conn)

func init() {

	logFile, err := os.OpenFile(LOG_DIRECTORY, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Println("log file create failure!")
		os.Exit(-1)
	}

	logger=log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}

}
func main() {
	defer logFile.Close()
	// 监听socket本地8080端口
	lis_socket, err := net.Listen("tcp", "127.0.0.1:8083")
	CheckErr(err)
	defer lis_socket.Close()

	fmt.Println("Service is waiting......")
	logger.Println("i am here")
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
		sendMesg := strings.Join(contents[1:], "#") //防止传输的消息体也有#
		addr = strings.Trim(addr, " ")
		//通过map查看是否有客户端，没有则不能发送
		if conn, ok := onlinesConns[addr]; ok {

			_, err := conn.Write([]byte(sendMesg))
			if err != nil {
				fmt.Println("online conns send failure!")
			}

		}
	} else {
		contents := strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST" {
			ips := ""
			for i := range onlinesConns {

				ips = ips + "|" + i

			}
			if conn, ok := onlinesConns[contents[0]]; ok {

				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("online conns send failure!")
				}

			}

		}

	}
}

// 负责接收协程

func ProcessInfo(conn net.Conn) {

	buf := make([]byte, 1024)
	defer func() {

		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlinesConns, addr)

		conn.Close()
	}()

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
