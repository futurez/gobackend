package main

import (
	"os"

	"github.com/zhoufuture/golite/config"
	"github.com/zhoufuture/golite/util"
)

var (
	TcpFileWrite *FileWriter
	UdpFileWrite *FileWriter
)

func main() {
	iniconf, err := config.NewConfig(config.IniProtocol, "config.ini")
	if err != nil {
		os.Exit(-1)
	}

	nettype := iniconf.GetString("logserver.nettype", "udp")
	tcpPort, _ := iniconf.GetInt("logserver.tcpport", 10000)
	udpPort, _ := iniconf.GetInt("logserver.udpport", 20000)
	tcpfilepath := iniconf.GetString("logserver.tcpfilepath", "/data/tcplog")
	udpfilepath := iniconf.GetString("logserver.udpfilepath", "/data/udplog")

	if nettype == "tcp" || nettype == "all" {
		TcpFileWrite, err = NewFileWriter(tcpfilepath)
		util.CheckError(err)
		go StartTcpListen(tcpPort)
	}

	if nettype == "udp" || nettype == "all" {
		UdpFileWrite, err = NewFileWriter(udpfilepath)
		util.CheckError(err)
		go StartUdpListen(udpPort)
	}

	waiting := make(chan byte)
	<-waiting
}
