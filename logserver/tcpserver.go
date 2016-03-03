package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/zhoufuture/golite/logger"
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

	var dataBuf bytes.Buffer
	buf := make([]byte, 2048)
	for {
		//		conn.SetDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			//			conn.SetDeadline(time.Time{})
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
			g_fileLog.Normal(string(writeBuf))
			dataBuf.Reset()
		}
	}
}

func checkError(err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func StartTcpListen(port int) {
	hostandip := fmt.Sprintf(":%d", port)
	//	listener, err := net.Listen("tcp", addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", hostandip)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	defer listener.Close()

	logger.Info("listener for server. local address: %s", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		logger.Debug("%d accept new connect, remote address: %s.", port, conn.RemoteAddr().String())
		go handleTcpConnect(conn)
	}
}
