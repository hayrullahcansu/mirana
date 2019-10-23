package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hayrullahcansu/mirana/abc"
	"github.com/hayrullahcansu/mirana/core/service"
	"github.com/hayrullahcansu/mirana/cross"
	"github.com/hayrullahcansu/mirana/xyz"
)

func main() {
	t1 := &xyz.XYZ{}

	fmt.Println(t1.TestMessage("TEST1"))
	t2 := &abc.ABC{}
	fmt.Println(t2.TestMessage("TEST2"))
	t1.Reis()
	t2.Reis()

	name := flag.String("name", "Mirana Game Server", "Mirana Game Server")
	// configPath := flag.String("config", "app.json", "config file")
	flag.Parse()

	// arguments stuffs
	// config.SetConfigFilePath(*configPath)

	log.Printf("Starting service for %s%s", *name, cross.NewLine)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s%s", s, cross.NewLine)
		AppCleanup()
		os.Exit(1)
	}()
	service.RunHandlers()
	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(1000))
	log.Println("CLEANUP APP BEFORE EXIT!!!")
}

// func main() {

// 	for index := 0; index < count; index++ {

// 	}

// 	if x := recover(); x != nil {
// 		log.Printf("run time panic: %v", x)
// 		//TODO: save the state and initlaize again.
// 		service.RunHandlers()
// 	}
// }
