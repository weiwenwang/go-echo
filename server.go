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
		fmt.Println("避免for一直循环,--------")
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if (err == io.EOF) {
			fmt.Println("读取完了")
			break
		}
		fmt.Println(n)
		util.CheckErr(err)
		fmt.Println("client:", string(buffer))
		// 不是读到数据就往管道里面发送的，要判断是不是一次完成的数据
		sendChannel <- string(buffer[:n])
	}
}

func handleStdIo() {
	for {
		fmt.Println("避免for一直循环,--------")
		read := bufio.NewReader(os.Stdin)
		line, err := read.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if (line != "") {
			util.CheckErr(err)
			fmt.Println(len(all_chan))
			for _, num := range all_chan {
				*num <- line
			}
		}
	}
}

func send(conn net.Conn, sendChannel chan string) {
	for {
		fmt.Println("避免for一直循环,--------")
		c := <-sendChannel
		c = strings.Replace(c, "\n", "", -1)
		fmt.Println("server:", string(c))
		conn.Write([]byte(c))
	}
}

var all_chan [] *chan string

func main() {
	tcpaddr, err1 := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	util.CheckErr(err1)
	listconn, err2 := net.ListenTCP("tcp4", tcpaddr)
	util.CheckErr(err2)

	go handleStdIo()

	for {
		fmt.Println("避免for一直循环,--------")
		conn, err3 := listconn.Accept()
		defer conn.Close()
		util.CheckErr(err3)
		sendCha := make(chan string, 10000)
		all_chan = append(all_chan, &sendCha)
		go handleConnection(conn, sendCha)

		go send(conn, sendCha)
	}

}
