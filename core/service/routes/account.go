package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hayrullahcansu/mirana/core/types"
)

//LoginHandler handles login requests and authorize user who is valid
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var msg json.RawMessage
	env := types.Envelope{
		Message: &msg,
	}
	if err := json.Unmarshal(b, &env); err != nil {
		log.Fatal(err)
	}
	switch env.MessageCode {
	case types.ELocation:
		var location types.Location
		if err := json.Unmarshal(msg, &location); err != nil {
			log.Fatal(err)
		}
	case types.EEvent:
	}

	w.Header().Set("content-type", "application/json")
	// msg := messaging.Response{Result: true, ContentCode: 1, Data: "logged succesfully"}
	// if message, err := json.Marshal(msg); err == nil {
	// 	w.Write(message)
	// }
	//TODO: parse request json and authorize
}

//RegisterHandler handles register requests and authorize user who is valid
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	//TODO: parse request json and authorize
}
