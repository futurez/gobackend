package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/zhoufuture/golite/logger"
	"github.com/zhoufuture/golite/util"
)

func checkPacket(buf []byte) ([]byte, bool) {
	if len(buf) < 2 {
		return nil, false
	}

	bufReader := bytes.NewReader(buf[:2])
	var bufLen int16
	if err := binary.Read(bufReader, binary.LittleEndian, &bufLen); err != nil {
		return nil, false
	}

	if len(buf) < int(bufLen+2) {
		return nil, false
	}
	return buf[2 : bufLen+2], true
}

func handleTcpConnect(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		var dataBuf bytes.Buffer
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				logger.Warn("%s connection is closed.", conn.RemoteAddr().String())
			} else {
				logger.Error("Read Error: %s", err.Error())
			}
			break
		}
		dataBuf.Write(buf[:n])
		writeBuf, b := checkPacket(dataBuf.Bytes())
		if b {
			if writeBuf[len(writeBuf)-1] != byte('\n') {
				writeBuf = append(writeBuf, '\n')
			}
			TcpFileWrite.Write(writeBuf)
			dataBuf.Reset()
		}
	}
}

func StartTcpListen(port int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	util.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	util.CheckError(err)
	defer listener.Close()
	logger.Info("listener for server. address: %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		logger.Debug("%d accept new connect, remote address: %s.", port, conn.RemoteAddr().String())
		go handleTcpConnect(conn)
	}
}
