package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"time"

	"github.com/hayrullahcansu/mirana/core/service"
	"github.com/hayrullahcansu/mirana/cross"
)

type aut struct {
	name    string
	article int
	id      int
}

func main() {
	// fmt.Println(t1.TestMessage("TEST1"))
	// t2 := &abc.ABC{}
	// fmt.Println(t2.TestMessage("TEST2"))
	// t1.Reis()
	// t2.Reis()
	Author := []struct {
		a_name    string
		a_article int
		a_id      int
	}{
		{"Mina", 304, 1098},
		{"Cina", 634, 102},
		{"Tina", 104, 105},
		{"Rina", 10, 108},
		{"Sina", 234, 103},
		{"Vina", 237, 106},
		{"Rohit", 56, 107},
		{"Mohit", 300, 104},
		{"Riya", 4, 101},
		{"Sohit", 20, 110},
	}

	// Sorting Author by their name
	// Using Slice() function
	sort.Slice(Author, func(p, q int) bool {
		return Author[p].a_name < Author[q].a_name
	})
	for _, val := range Author {
		fmt.Println(val.a_name)
	}
	// Checking the slice is sorted
	// according to their names
	// Using SliceIsSorted function
	res1 := sort.SliceIsSorted(Author, func(p, q int) bool {
		return Author[p].a_name < Author[q].a_name
	})

	if res1 == true {

		fmt.Println("Slice is sorted by their names")

	} else {

		fmt.Println("Slice is not sorted by their names")
	}

	// Checking the slice is sorted
	// according to their total articles
	// Using SliceIsSorted function
	res2 := sort.SliceIsSorted(Author, func(p, q int) bool {
		return Author[p].a_article < Author[q].a_article
	})

	if res2 == true {

		fmt.Println("Slice is sorted by " +
			"their total number of articles")

	} else {

		fmt.Println("Slice is not sorted by" +
			" their total number of articles")
	}

	// Sorting Author by their ids
	// Using Slice() function
	sort.Slice(Author, func(p, q int) bool {
		return Author[p].a_id < Author[q].a_id
	})

	// Checking the slice is sorted
	// according to their ids
	// Using SliceIsSorted function
	res3 := sort.SliceIsSorted(Author, func(p, q int) bool {
		return Author[p].a_id < Author[q].a_id
	})

	if res3 == true {

		fmt.Println("Slice is sorted by their ids")

	} else {

		fmt.Println("Slice is not sorted by their ids")
	}
	dd := make(map[string]*aut)
	dd["10"] = &aut{
		id: 10,
	}
	for i := 0; i < 5; i++ {
		d := strconv.Itoa(i)
		dd[d] = &aut{
			id: i,
		}
	}
	for key, _ := range dd {
		fmt.Println(key)
	}

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
