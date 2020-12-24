package LoadBalance

import (
	"fmt"
	"hash/crc32"
)

func(this *LoadBalance) SelectByIPHash(ip string) (NodeBalance,error) { //ip_hash算法

	if this.allDown() {
		return nil,fmt.Errorf("节点都不可用")
	}
	index:=int(crc32.ChecksumIEEE([]byte(ip))) % len(this.nodes)
	if this.nodes[index].GetStatus()!=Ready {
		return  nil,fmt.Errorf("当前节点都不可用")
	}
	return  this.nodes[index],nil
}