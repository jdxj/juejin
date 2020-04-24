package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jdxj/juejin/db"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/juejin/app"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"juejin.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	coll := app.NewCollector()
	coll.Start()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-sig:
		logs.Info("receive signal")
	}

	coll.Stop()
	db.MySQL.Close()
	logs.GetBeeLogger().Flush()
}
