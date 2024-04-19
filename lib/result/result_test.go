package result

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	fmt.Println(NewResult[int](ret()).UnwrapDefault(3121312))
}
func ret() (int, error) {
	if time.Now().Unix()%2 == 0 {
		return 0, fmt.Errorf("111")
	}
	return 10, nil
}
