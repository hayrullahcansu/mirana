package routes

import (
	"net/http"

	"bitbucket.org/digitdreamteam/mirana/core/comm/netsp"
	"bitbucket.org/digitdreamteam/mirana/core/mng/singledeck"
)

//JoinRoomMode1Handler hnadles login requests and authorize user who is valid
func JoinRoomNormalGameHandler(w http.ResponseWriter, r *http.Request) {
	c := netsp.NewClient()
	c.ServeWs(w, r)
	singledeck.Manager().RequestPlayGame(c)
}
