package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/service"
	"bitbucket.org/digitdreamteam/mirana/cross"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

type aut struct {
	name    string
	article int
	id      int
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	logrus.SetReportCaller(false)

	name := flag.String("name", "Mirana Game Server", "Mirana Game Server")
	// configPath := flag.String("config", "app.json", "config file")
	flag.Parse()

	// arguments stuffs
	// config.SetConfigFilePath(*configPath)

	logrus.Infof("Starting service for %s%s", *name, cross.NewLine)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	go func() {
		s := <-sigs
		logrus.Infof("RECEIVED SIGNAL: %s%s", s, cross.NewLine)
		AppCleanup()
		os.Exit(1)
	}()
	service.RunHandlers()
	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(1000))
	logrus.Infof("CLEANUP APP BEFORE EXIT!!!")
}

// func main() {

// 	for index := 0; index < count; index++ {

// 	}

// 	if x := recover(); x != nil {
// 		log.Printf("run time panic: %v", x)
// 		service.RunHandlers()
// 	}
// }
