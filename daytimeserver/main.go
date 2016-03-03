// main
package main

import (
	"net"
	"time"

	"github.com/zhoufuture/golite/logger"
	"github.com/zhoufuture/golite/util"
)

func main() {
	service := ":9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	util.CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	util.CheckError(err)
	logger.Info("listen %s , %s", tcpAddr.Network(), tcpAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Warn("accept failed, %s", err.Error())
			continue
		}
		logger.Debug("accept new client, %s", conn.RemoteAddr().String())
		daytime := time.Now().String()
		conn.Write([]byte(daytime))
		conn.Close()
	}
}
