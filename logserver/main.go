package main

import (
	"encoding/json"
	"os"

	"github.com/zhoufuture/golite/config"
	"github.com/zhoufuture/golite/logger"
)

var (
	g_tcpPort     int
	g_udpPort     int
	g_tcpfilepath string
	g_udpfilepath string
	g_fileLog     *logger.Logger
)

func init() {
	iniconf, err := config.NewConfig(config.IniProtocol, "config.ini")
	if err != nil {
		os.Exit(-1)
	}

	g_tcpPort, _ = iniconf.GetInt("logserver.tcpport", 30000)
	g_udpPort, _ = iniconf.GetInt("logserver.udpport", 40000)
	g_tcpfilepath = iniconf.GetString("logserver.tcpfilepath", "data/tcplog/log")
	g_udpfilepath = iniconf.GetString("logserver.udpfilepath", "data/udplog/log")

	g_fileLog = logger.NewLogger(10000)

	var fileconf logger.FileLogConfig
	fileconf.FileName = g_tcpfilepath
	fileconf.LogFlag = 0
	fileconf.MaxDays = 7
	fileconf.MaxSize = 1 << 30
	fileconfbuf, _ := json.Marshal(fileconf)
	g_fileLog.SetLogger(logger.FILE_PROTOCOL_LOG, string(fileconfbuf))

	g_fileLog.SetFuncCallDepth(0)
	g_fileLog.SetEnableFuncCall(false)
	g_fileLog.SetLogLevel(logger.LevelNormal)
	g_fileLog.Async()
}

func main() {

	StartTcpListen(g_tcpPort)
	//StartUdpListen(g_udpport)
}
