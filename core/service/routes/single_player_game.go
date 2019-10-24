package routes

import (
	"net/http"

	"github.com/hayrullahcansu/mirana/core/comm/netsp"
	"github.com/hayrullahcansu/mirana/core/mng/sp"
)

//JoinRoomMode1Handler hnadles login requests and authorize user who is valid
func JoinRoomNormalGameHandler(w http.ResponseWriter, r *http.Request) {
	c := netsp.NewClient()
	c.ServeWs(w, r)
	sp.Manager().RequestPlayGame(c)
}
