package routes

import (
	"net/http"
)

//JoinLobbyHandler hnadles login requests and authorize user who is valid
func JoinLobbyHandler(w http.ResponseWriter, r *http.Request) {

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
