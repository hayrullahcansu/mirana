package routes

import (
	"net/http"

	"bitbucket.org/digitdreamteam/mirana/core/comm/netl"
	"bitbucket.org/digitdreamteam/mirana/core/mng/lobby"
)

//JoinLobbyHandler hnadles login requests and authorize user who is valid
func JoinLobbyHandler(w http.ResponseWriter, r *http.Request) {

	c := netl.NewClient()
	c.ServeWs(w, r)
	lobby.Manager().ConnectLobby(c)
	// t := netsp.NetSPClient{
	// 	BaseClient: c,
	// }
	// c := netsp.NetSPClient.ServeWs(w, r)
	// c := netsp.ServeWs(w, r)

	// _data, err := ioutil.ReadAll(r.Body)
	// if err == nil {
	// 	re := messaging.Response{Result: true, ContentCode: 1, Data: string(_data[:])}
	// 	data, _ := json.Marshal(re)
	// 	w.Write(data)
	// 	return
	// }

	//GetJWT

	//GETUSERDATA

	//SAVEONMAPSERVER
	// user := &dto.PlayerDto{
	// 	Id:         1,
	// 	SkinId:     "skin1",
	// 	Nick:       "test",
	// 	BodyRatio:  1.2,
	// 	ForceRatio: 1.2,
	// }
	// gs.Users[user.Id] = user
	// if lobbyPlayer, ok := management.CreateNewLobbyPlayer(gs, user, w, r); ok {
	// 	gs.Lobby.register <- lobbyPlayer
	// }

}
