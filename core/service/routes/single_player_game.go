package routes

import (
	"fmt"
	"net/http"

	"bitbucket.org/digitdreamteam/mirana/core/mng/blackjack"
	"bitbucket.org/digitdreamteam/mirana/utils"
	"github.com/sirupsen/logrus"
)

//JoinRoomMode1Handler hnadles login requests and authorize user who is valid
func JoinRoomNormalGameHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("A client joint SingleDeck Game Room\n%s", utils.FormatRequest(r))
	userId := r.URL.Query().Get("user-id")
	if userId != "" {
		fmt.Println("UserId:" + userId)
		c := blackjack.NewClient(userId)
		c.ServeWs(w, r)	
		blackjack.Manager().RequestSingleDeckPlayGame(c)
	}
}
