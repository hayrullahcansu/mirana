package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hayrullahcansu/zapper/man"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	server := man.NewServer()
	go server.Run()
	http.HandleFunc("/gameroom", func(w http.ResponseWriter, r *http.Request) {

		htmlData, err := ioutil.ReadAll(r.Body) //<--- here!
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// print out
		fmt.Println(string(htmlData)) //<-- here !
		man.ServeWs(server, w, r)
	})
	log.Print("ListenAndServe: ", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
