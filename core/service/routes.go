package service

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hayrullahcansu/mirana/core/service/routes"
)

var addr = flag.String("addr", "127.0.0.1:3535", "http service address")

var loginHandlerRoute = flag.String("loginHandlerRoute", "/login", "login handler function route")
var registerHandlerRoute = flag.String("registerHandlerRoute", "/register", "login handler function route")
var joinLobbyHandlerRoute = flag.String("joinLobbyHandlerRoute", "/joinlobby", "login handler function route")
var joinRoomNormalGameHandlerRoute = flag.String("joinRoomMode1HandlerRoute", "/join_sp_game", "login handler function route")
var joinRoomRankedGameHandlerRoute = flag.String("joinRoomMode2HandlerRoute", "/join_mp_game", "login handler function route")

var appLiveVersion = flag.String("app_live_version", "1.0.0", "http spinner handler function path")

//RunHandlers provite to handle requests and redirect
func RunHandlers(gameServer *GameServer) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	http.HandleFunc(*loginHandlerRoute, routes.LoginHandler)
	http.HandleFunc(*registerHandlerRoute, routes.RegisterHandler)

	http.HandleFunc(*joinLobbyHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinLobbyHandler(gameServer, w, r)
	})
	http.HandleFunc(*joinRoomNormalGameHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinRoomNormalGameHandler(gameServer, w, r)
	})
	http.HandleFunc(*joinRoomRankedGameHandlerRoute, func(w http.ResponseWriter, r *http.Request) {
		routes.JoinRoomRankedGameHandler(gameServer, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		log.Fatal("ListenAndServe: ", err)
	}
}
