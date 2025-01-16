package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/pkg/os"
	ctl "github.com/wike2019/wike_go/pkg/service/ctl"
)

func (this CoreCtl) SystemInfo(context *gin.Context) {
	c := this.SetContext(context)
	path := c.DefaultQuery("path", "/")
	var s os.Server
	s.Os = os.InitOS()
	cpu, err := os.InitCPU()
	if err != nil {
		ctl.Error(err, 400)
	}
	s.Cpu = cpu
	ram, err := os.InitRAM()
	if err != nil {
		ctl.Error(err, 400)
	}
	s.Ram = ram
	disk, err := os.InitDisk(path)
	if err != nil {
		ctl.Error(err, 400)
	}
	s.Disk = disk
	c.Success("获取服务器信息成功", s)
}
