// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/4/16  9:38 下午

package logker

// file v1.0.6
// File is the implementation of logger interface.

import (
	"fmt"
	"os"
	"path"
)

const (
	// Error File Prefix
	errPrefix = "error_"
	// File Perfix
	suffix = ".log"
	// Backup File Suffix
	bakSuffix = "_bak.log"
	// Backup File Perfix
	bakPerfix = "log_"
)

// customize logging type
type logFileType int

const (
	// Common plain Logging file type
	plain logFileType = iota
	// Major error logging file type
	major
)

type fileLog struct {
	logLevel level
	// Whether enable error log file.
	wheError bool
	// Log save directory
	directory string
	// Log file name
	fileName string
	// Log file  pointer
	file *os.File
	// Error file log pointer
	errFile *os.File
	// Customize of time zone type
	tz *timeZone
	// Set running Time Zone
	timeZone logTimeZone
	// File system Power
	power os.FileMode
	// File Max size
	fileMaxSize int64
	// MessageMatchingCard
	formatting string
	asyncTask  *AsyncTask
}

// Initialization error file pointer
func (f *fileLog) initErrPtr() *os.File {
	savePath := path.Join(f.directory, errPrefix+f.fileName+suffix)
	file, e := os.OpenFile(savePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, f.power)
	if e != nil {
		panic("open file fail :" + savePath)
	}
	return file
}

//	Initialization file pointer
func (f *fileLog) initFilePtr() *os.File {
	savePath := path.Join(f.directory, f.fileName+".log")
	// fmt.Println(savePath)
	file, e := os.OpenFile(savePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, f.power)
	// fmt.Printf("3 %T %p \n", file, file)
	if e != nil {
		panic("open file fail :" + savePath)
	}
	return file
}

func (f *fileLog) isEnableErr() bool {
	return f.wheError
}

// TODO: Whether enable current level
func (f *fileLog) isEnableLevel(lev level) bool {
	// debug<info<warn<Error
	return f.logLevel <= lev
}

// TODO: Check log size & whether division log file
func (f *fileLog) checkSize() bool {
	info, e := f.file.Stat()
	if e != nil {
		return false
	}
	return info.Size() >= f.fileMaxSize
}
func (f *fileLog) checkErrSize() bool {
	info, e := f.errFile.Stat()
	if e != nil {
		return false
	}
	return info.Size() >= f.fileMaxSize
}
func (f *fileLog) Info(value string, args ...interface{}) {
	if f.isEnableLevel(INFO) {
		f.sendMsg(f.pack(INFO, fmt.Sprintf(value, args...)))
	}
}
func (f *fileLog) Debug(value string, args ...interface{}) {
	if f.isEnableLevel(DEBUG) {
		f.sendMsg(f.pack(DEBUG, fmt.Sprintf(value, args...)))
	}
}
func (f *fileLog) Error(value string, args ...interface{}) {
	if f.isEnableLevel(ERROR) {
		f.sendMsg(f.pack(ERROR, fmt.Sprintf(value, args...)))
	}
}
func (f *fileLog) Warning(value string, args ...interface{}) {
	if f.isEnableLevel(WARNING) {
		f.sendMsg(f.pack(WARNING, fmt.Sprintf(value, args...)))
	}
}

// Division logging file.
func (f *fileLog) divisionLogFile(fileType logFileType) {
	switch fileType {
	case plain:
		_ = f.file.Close()
		srcPath := path.Join(f.directory, f.fileName+suffix)
		newPath := path.Join(f.directory, bakPerfix+f.tz.NowTimeStrLogName()+bakSuffix)
		_ = os.Rename(srcPath, newPath)
		f.file = f.initFilePtr()
	case major:
		_ = f.errFile.Close()
		srcPath := path.Join(f.directory, errPrefix+f.fileName+suffix)
		newPath := path.Join(f.directory, errPrefix+f.tz.NowTimeStrLogName()+bakSuffix)
		_ = os.Rename(srcPath, newPath)
		f.errFile = f.initErrPtr()
	default:
		f.Error("division log file fail. Type: %v", fileType)
	}

}
