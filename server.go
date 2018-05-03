package main

import (
	"echo/util"
	"net"
	"bufio"
	"os"
	"fmt"
	"io"
	"strings"
)

func handleConnection(conn net.Conn, sendChannel chan string) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if (err == io.EOF) {
			fmt.Println("读取完了")
			break
		}
		util.CheckErr(err)
		fmt.Println("client:", string(buffer))
		// 不是读到数据就往管道里面发送的，要判断是不是一次完成的数据
		sendChannel <- string(buffer[:n])
	}
}

func handleStdIo() {
	for {
		read := bufio.NewReader(os.Stdin)
		line, err := read.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1) // 去掉标准输入后面带的\n
		if (line != "") { // 如果输入有内容就处理
			util.CheckErr(err)
			for _, num := range all_chan { // 遍历每一个channel, 把数据传给所有客户端
				*num <- line
			}
		}
	}
}

func send(conn net.Conn, sendChannel chan string) {
	for {
		c := <-sendChannel // 只要管道里面有数据要发，就取出来发掉，不然就阻塞在这
		c = strings.Replace(c, "\n", "", -1)
		fmt.Println("server:", string(c))
		conn.Write([]byte(c))
	}
}

// 用一个全局变量存所有客户端的channel， 这个是为了，把服务器从标准输入传入数据传给每一个连接到客户端的channel
var all_chan [] *chan string

func main() {
	tcpaddr, err1 := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	util.CheckErr(err1)
	listconn, err2 := net.ListenTCP("tcp4", tcpaddr)
	util.CheckErr(err2)

	// 开一个goroutine处理
	go handleStdIo()

	for {
		conn, err3 := listconn.Accept() // 循环的处理新建立的连接，然后开goroutine处理具体的业务，这样就不会阻塞accept
		defer conn.Close()
		util.CheckErr(err3)
		sendCha := make(chan string, 10000)   // 为每一个连接开一个管道，用于该该客户端的接受和发送goroutine通信
		all_chan = append(all_chan, &sendCha) // 把改客户端的channel放入到全局变量里面
		go handleConnection(conn, sendCha)    // 为该客户端创建一个接受的goroutine
		go send(conn, sendCha)                // 为该客户端创建一个接受的goroutine
	}

}
