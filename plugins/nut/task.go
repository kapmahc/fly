package nut

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

func monitorTask() error {
	beego.Info("start monitor task")
	return nil
}

func init() {
	for _, t := range []*toolbox.Task{
		toolbox.NewTask("monitor", "0 */5 * * * *", monitorTask),
	} {
		toolbox.AddTask(t.Taskname, t)
	}
}
