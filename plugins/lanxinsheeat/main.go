package lanxinsheeat

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
	"sync"
	"time"
)

var (
	waitSendImage sync.Map
)

func init() {
	control.Register("lanxinsheeat", &control.Options{
		Alias: "澜心社吃饭推送",
		Help: "指令:\n" +
			"* 早报 -> 获取每天60s读懂世界\n" +
			"* 每日早报 -> 获取每天60s读懂世界\n" +
			"* 早报定时 -> 专门用于定时任务的指令，请不要手动调用",
		DataFolder: "lanxinsheeat",
		OnEnable: func(ctx *robot.Ctx) {
			// todo 启动将定时任务加入到定时任务列表
			ctx.ReplyText("启用成功")
		},
		OnDisable: func(ctx *robot.Ctx) {
			// todo 停止将定时任务从定时任务列表移除
			ctx.ReplyText("禁用成功")
		},
	})

	go pollingTask()

}

func pollingTask() {
	// 计算下一个整点
	now := time.Now().Local()
	next := now.Add(10 * time.Minute).Truncate(10 * time.Minute)
	diff := next.Sub(now)
	timer := time.NewTimer(diff)
	<-timer.C
	timer.Stop()

	// 任务
	doSendImage := func(text string) {
		waitSendImage.Range(func(key, val interface{}) bool {
			ctx := val.(*robot.Ctx)
			ctx.ReplyText(text)

			time.Sleep(10 * time.Second)
			return true
		})
	}

	// 轮询任务
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		// 避开0点-5点(应该不会有人设置这个时间吧)
		if time.Now().Hour() < 5 {
			continue
		}
		doSendImage("我是测试的主动发送")
	}
}
