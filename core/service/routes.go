package service

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/digitdreamteam/mirana/core/service/routes"
)

var addr = flag.String("addr", "127.0.0.1:3535", "http service address")

var loginHandlerRoute = flag.String("loginHandlerRoute", "/login", "login handler function route")
var registerHandlerRoute = flag.String("registerHandlerRoute", "/register", "login handler function route")
var joinLobbyHandlerRoute = flag.String("joinLobbyHandlerRoute", "/joinlobby", "login handler function route")
var joinRoomNormalGameHandlerRoute = flag.String("joinRoomMode1HandlerRoute", "/join_sp_game", "login handler function route")
var joinRoomAmericanGameHandlerRoute = flag.String("joinRoomAmericanHandlerRoute", "/join_ap_game", "login handler function route")
var joinRoomRankedGameHandlerRoute = flag.String("joinRoomMode2HandlerRoute", "/join_mp_game", "login handler function route")

var appLiveVersion = flag.String("app_live_version", "1.0.0", "http spinner handler function path")

//RunHandlers provite to handle requests and redirect
func RunHandlers() {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
			//TODO: save the state and initlaize again.
			RunHandlers()
		}
	}()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	//TODO: init and initialize old game states / recover if the server broken
	http.HandleFunc(*loginHandlerRoute, routes.LoginHandler)
	http.HandleFunc(*registerHandlerRoute, routes.RegisterHandler)

	http.HandleFunc(*joinLobbyHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinLobbyHandler(w, r)
	})
	http.HandleFunc(*joinRoomNormalGameHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinRoomNormalGameHandler(w, r)
	})
	http.HandleFunc(*joinRoomAmericanGameHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinRoomAmericanGameHandler(w, r)
	})
	http.HandleFunc(*joinRoomRankedGameHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinRoomRankedGameHandler(w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		log.Fatal("ListenAndServe: ", err)
	}
}
