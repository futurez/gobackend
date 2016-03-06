package main

import (
	"fmt"
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
		buf = buf[:n]
		if err != nil || n == 0 {
			logger.Warn("%s %s %s, readLen=%d", addr.Network(), addr.String(), err.Error(), n)
			continue
		}
		if buf[n-1] != byte('\n') {
			buf = append(buf, '\n')
		}
		UdpFileWrite.Write(buf)
	}
}
