package charger

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"

	"github.com/drewwells/charger/store"
)

func init() {
	http.HandleFunc("/users/", handleUser)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleNewUser(w, r)
	case "PUT":
		handleUpdateUser(w, r)
	case "GET":
		handleUsersToNotify(w, r)
	default:
		http.Error(w, "", http.StatusNotImplemented)
	}
}

func handleNewUser(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var u store.User
	err = json.Unmarshal(bs, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := appengine.NewContext(r)
	err = store.NewUser(ctx, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var u store.User
	err = json.Unmarshal(bs, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	err = store.UpdateUser(ctx, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(u)
}

type UserResponse struct {
	Users []store.User
}

func handleUsersToNotify(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	users, err := store.GetEmails(ctx, &store.User{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := UserResponse{
		Users: users,
	}
	json.NewEncoder(w).Encode(resp)
}
