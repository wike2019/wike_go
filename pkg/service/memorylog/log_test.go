package memorylog

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	count := 0
	f := &LogStruct{}
	for {
		count = count + 1

		f.Write([]byte(fmt.Sprintf("%d", count)))
		//fmt.Println(count)
		if count == 3000 {
			break
		}
	}
	fmt.Println("ok")

	f.Show()
}
