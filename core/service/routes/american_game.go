package routes

import (
	"fmt"
	"net/http"

	"bitbucket.org/digitdreamteam/mirana/core/mng/american"
)

//JoinRoomAmericanGameHandler hnadles login requests and authorize user who is valid
func JoinRoomAmericanGameHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user-id")
	if userId != "" {
		fmt.Println("UserId:" + userId)
		c := american.NewClient(userId)
		c.ServeWs(w, r)
		american.Manager().RequestPlayGame(c)
	}
}
