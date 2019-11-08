package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/service"
	"bitbucket.org/digitdreamteam/mirana/cross"
)

type aut struct {
	name    string
	article int
	id      int
}

func main() {

	a := []int{1, 2, 4, 5, 6}
	b := 3

	// Make space in the array for a new element. You can assign it any value.
	a = append(a, 0)
	fmt.Println(a)

	// Copy over elements sourced from index 2, into elements starting at index 3.
	copy(a[3:], a[2:])
	fmt.Println(a)

	a[2] = b
	fmt.Println(a)
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
		{"Rohit", 30, 107},
		{"Mina", 20, 1098},
		{"Cina", 21, 102},
		{"Tina", 25, 105},
		{"Rina", 18, 108},
		{"Mohit", 21, 104},
		{"Riya", 10, 101},
		{"Sohit", 5, 110},
		{"Tina", 25, 105},
	}
	s := []int{30, 2, 20, 21, 30, 21, 22, 25, 17, 19} // unsorted
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j] && s[i] <= 21
	})

	for _, val := range s {
		fmt.Printf("%d \n", val)
	}
	// // Sorting Author by their name
	// Using Slice() function
	sort.Slice(Author, func(p, q int) bool {
		return Author[p].a_article < Author[q].a_article
	})
	// for _, val := range Author {
	// 	fmt.Printf("%s %d \n", val.a_name, val.a_article)
	// }

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
