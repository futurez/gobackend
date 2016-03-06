package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhoufuture/golite/logger"
)

const (
	MAXSIZE = (1 << 30)
	MAXDAYS = 7
)

type FileWriter struct {
	fd         *os.File
	curFileNum int
	curSize    int
	curDate    int
	fileName   string
	pathName   string
	bufChan    chan []byte
}

func NewFileWriter(pathname string) (*FileWriter, error) {
	if len(pathname) == 0 {
		return nil, errors.New("pathname is null.")
	}

	fw := &FileWriter{
		fileName: pathname + "/log",
		pathName: pathname,
		bufChan:  make(chan []byte, 1000),
	}
	go fw.Save()
	return fw, fw.openFile()
}

func (fw *FileWriter) Save() {
	for buf := range fw.bufChan {
		fw.fd.Write(buf)
		fw.curSize += len(buf)
		fw.docheck()
	}
}

//every-time, call write must new []byte, because we use chan []byte,
//and []byte is reference, not value pass.
func (fw *FileWriter) Write(b []byte) {
	fw.bufChan <- b
}

func (fw *FileWriter) docheck() {
	if (fw.curSize >= MAXSIZE) || (time.Now().Day() != fw.curDate) {
		if err := fw.backupFile(); err != nil {
			logger.Error(err.Error())
		}
	}
}

func (fw *FileWriter) openFile() error {
	err := os.MkdirAll(fw.pathName, 0660)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	fd, err := os.OpenFile(fw.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stat, err := fd.Stat()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if fw.fd != nil {
		fw.fd.Close()
	}
	fw.curDate = time.Now().Day()
	fw.curSize = int(stat.Size())
	fw.fd = fd
	return nil
}

func (fw *FileWriter) backupFile() error {
	_, err := os.Lstat(fw.fileName)
	if err != nil { // file not exists
		logger.Error("why log file not exists.")
		return err
	}

	// Find the next available number
	if time.Now().Day() != fw.curDate {
		fw.curFileNum = 0
	}

	fname := ""
	for err == nil {
		fname = fmt.Sprintf("%s.%s.%03d", fw.fileName, time.Now().Format("2006-01-02"), fw.curFileNum)
		_, err = os.Lstat(fname)
		fw.curFileNum++
	}

	fw.fd.Close()
	err = os.Rename(fw.fileName, fname)
	if err != nil {
		return err
	}

	// re-create logger
	err = fw.openFile()
	if err != nil {
		return err
	}

	go fw.cleanUp()
	return nil
}

func (fw *FileWriter) cleanUp() {
	dir := filepath.Dir(fw.fileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		defer func() {
			if r := recover(); r != nil {
				e := fmt.Errorf("Unable to delete old log %s, error: %v", path, r)
				logger.Error(e.Error())
			}
		}()
		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*MAXDAYS) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(fw.fileName)) {
				os.Remove(path)
			}
		}
		return nil
	})
}
