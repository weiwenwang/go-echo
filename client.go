package main

import (
	"bufio"
	"os"
	"echo/util"
	"net"
	"fmt"
	"io"
	"strings"
)

func main() {
	tcpaddr, err1 := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	util.CheckErr(err1)

	conn, err2 := net.DialTCP("tcp", nil, tcpaddr)
	util.CheckErr(err2)
	defer conn.Close()
	sendChannel := make(chan string, 10000)

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if (err == io.EOF) {
				fmt.Println("读取完了")
				break
			} else {
				fmt.Println("从服务端收到", n)

				util.CheckErr(err)
				fmt.Println("server: ", string(buffer[:n]))
			}
		}
	}()

	go func() {
		for {
			read := bufio.NewReader(os.Stdin)
			line, err := read.ReadString('\n')
			line = strings.Replace(line, "\n", "", -1)
			if (line != "") {
				fmt.Println("stdio:", line)
				util.CheckErr(err)
				sendChannel <- line
			}

		}
	}()

	for {
		c := <-sendChannel
		c = strings.Replace(c, "\n", "", -1)
		conn.Write([]byte(c))
	}
}
