package recovery

import (
	"os"
	"os/signal"

	"../log"
	"../memory"
	"../table"
)

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
	log.WriteLog("sys", "Begin Restore.")
	defer log.WriteLog("sys", "Restore Finished.")
	memory.Restore()
	table.Restore()
}
