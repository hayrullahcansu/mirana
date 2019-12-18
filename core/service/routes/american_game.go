package routes

import (
	"fmt"
	"net/http"

	blackjack "bitbucket.org/digitdreamteam/mirana/core/mng/singledeck"
	"bitbucket.org/digitdreamteam/mirana/utils"
	"github.com/sirupsen/logrus"
)

//JoinRoomAmericanGameHandler hnadles login requests and authorize user who is valid
func JoinRoomAmericanGameHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("A client joint American Game Room\n%s", utils.FormatRequest(r))
	userId := r.URL.Query().Get("user-id")
	if userId != "" {
		fmt.Println("UserId:" + userId)
		c := blackjack.NewClient(userId)
		c.ServeWs(w, r)
		blackjack.Manager().RequestAmericanPlayGame(c)
	}
}
