package recovery

import (
	"os"
	"os/signal"

	"../index"
	"../log"
	"../memory"
	"../table"
)

var Restoring = 0

func SafeExit() {
	Restore()
}

func init() {
	// Capture Ctrl+C signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		log.WriteLog("sys", "Get Signal:"+s.String()+", Safe Exiting...")
		SafeExit()
	}()
	LoadData()
}

func Restore() {
	Restoring = 1
	log.WriteLog("sys", "Begin Restore.")
	defer log.WriteLog("sys", "Restore Finished.")
	memory.Restore()
	table.Restore()
	index.Restore()
	Restoring = 0
}
