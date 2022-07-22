package main

import (
	"github.com/P4elme6ka/go-http-media-server/src/app"
	"github.com/P4elme6ka/go-http-media-server/src/param"
	"os"
	"os/signal"
	"syscall"
)

var appInst *app.App

func cleanupOnInterrupt() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT)

	go func() {
		<-chSignal
		appInst.Close()
		os.Exit(0)
	}()
}

func reOpenLogOnHup() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			appInst.ReOpenLog()
		}
	}()
}

func main() {
	cleanupOnInterrupt()

	params := param.ParseCli()
	appInst = app.NewApp(params)

	if appInst != nil {
		reOpenLogOnHup()
		appInst.Open()
		defer appInst.Close()
	}
}
