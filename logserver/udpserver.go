package main

import (
	"fmt"
	"log"
	"net"

	"github.com/zhoufuture/golite/logger"
	"github.com/zhoufuture/golite/util"
)

func StartUdpListen(udpport int) {
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", udpport))
	util.CheckError(err)
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	util.CheckError(err)
	defer udpConn.Close()
	logger.Info("listener for server. address: %s", udpConn.LocalAddr().String())

	for {
		buf := make([]byte, 512)
		n, addr, err := udpConn.ReadFromUDP(buf[0:])
		if err != nil {
			logger.Warn("%s %s %s.", addr.Network(), addr.String(), err.Error())
			continue
		}
		UdpFileWrite.Write(buf[0:n])
	}
}
