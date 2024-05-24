package memorylog

import (
	"fmt"
)

// 内存日志
type Log struct {
	Text string
}

type LogList []*Log

type LogStruct struct {
	Log LogList
}

var LogInfo *LogStruct

func (this *LogStruct) Write(text []byte) (int, error) {
	var old LogList
	if len(this.Log) >= 2001 {
		old = this.Log[:2000]
	} else {
		old = this.Log
	}
	this.Log = append([]*Log{{Text: string(text)}}, old...)
	return len(text), nil
}

func (this *LogStruct) Show() {
	for _, item := range this.Log {
		fmt.Println(item.Text)
	}
}
func (this *LogStruct) All() LogList {
	return this.Log
}
