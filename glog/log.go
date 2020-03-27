package glog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

const (
	EMERGENCY = "[emergency]" //致命的错误, 导致系统奔溃
	ALTER     = "[alter]"     //严重的错误, 如mysql连接不可用
	CRITICAL  = "[critical]"  //警戒性错误, 如某个组件不可用, 服务不可用
	ERROR     = "[error]"     //一般错误, 数值不对等等
	WARN      = "[warning]"   //警告性错误, 程序纯在确定情况, 如switch中间case未处理
	NOTICE    = "[notice]"    //通知, 不一定是错误, 执行超时等等
	INFO      = "[info]"      //程序输出信息
	DEBUG     = "[debug]"     //调试信息
)

var (
	suffix               = ".log"                //文件后缀
	fileName             = "log"                 //文件名称
	dir                  = ""                    //文件地址
	dateTime             = "2006/01/02 15:04:05" //date
	megabyte       int64 = 1024 * 1024           //1MB
	defaultMaxSize int64 = 512                   //默认最大单位MB
	split                = false                 //是否切割
	traceFileLine        = true
)

var file *os.File

func InitLog() {
	file = openCreateFile()
	go waitWrite()
}

var readWriteChan = make(chan []byte)
var readChan <-chan []byte = readWriteChan
var writeChan chan<- []byte = readWriteChan

func waitWrite() {
	for {
		select {
		case content := <-readChan:
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				_, err := write(content)

				panic(err)
			}()
		}
	}
}

type Logger struct{}

type content struct {
	Level     string
	TimeLocal string
	Msg       string
	LineNo    int
	FilePath  string
	Context   interface{}
}

func openCreateFile() *os.File {
	findCreateDir()

	filePath := filepath.Join(dir, fileName+suffix)

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

//查找创建目录
func findCreateDir() {
	if dir == "" {
		dir = os.TempDir()
	}

	_, err := os.Stat(dir)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		if ok := os.Mkdir(dir, 0755); ok != nil {
			panic(ok)
		}

		return
	}

	panic(err)
}

//设置dir
func SetDir(d string) {
	dir = d
	findCreateDir()
}

func write(p []byte) (n int, err error) {
	n, err = file.Write(p)
	return
}

func SetFileName(name string) {
	fileName = name
}

func ReloadFile() {
	file = openCreateFile()
}

func (l *Logger) Write(p []byte) (n int, err error) {
	writeChan <- p
	return len(p), err
}

func Debug(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = DEBUG
	c.Context = v
	writeChan <- formatContext(c)
}

func Info(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = INFO
	c.Context = v
	writeChan <- formatContext(c)
}

func Notice(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = NOTICE
	c.Context = v
	writeChan <- formatContext(c)
}

func Warn(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = WARN
	c.Context = v
	writeChan <- formatContext(c)
}

func Error(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = ERROR
	c.Context = v
	writeChan <- formatContext(c)
}

func Critical(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = CRITICAL
	c.Context = v
	writeChan <- formatContext(c)
}

func Alter(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = ALTER
	c.Context = v
	writeChan <- formatContext(c)
}

func Emergency(msg string, v interface{}) {
	c := newContext()
	c.Msg = msg
	c.Level = EMERGENCY
	c.Context = v
	writeChan <- formatContext(c)
}

func newContext() *content {
	return &content{
		Level:     "",
		TimeLocal: time.Now().Format(dateTime),
		Msg:       "",
		LineNo:    0,
		FilePath:  "",
		Context:   nil,
	}
}

func formatContext(c *content) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(c.Level)
	buffer.WriteString(c.TimeLocal + " ")
	buffer.WriteString(c.Msg + " ")

	if c.Context != nil {
		buffer.WriteString("\r")

		switch c.Context.(type) {
		case string:
			buffer.WriteString(reflect.ValueOf(c.Context).String())
		case []byte:
			buffer.Write(reflect.ValueOf(c.Context).Bytes())
		default:
			s, _ := json.Marshal(c.Context)
			buffer.Write(s)
		}
	}

	buffer.WriteString("\r")
	return buffer.Bytes()
}
