package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/seternate/lanty-api-golang/pkg/user"
)

var userList *user.Users

func InitUsersHandler() {
	userList = &user.Users{}
	*userList = append(*userList, user.User{Id: 0, Name: "test1"})
}

func HandleGetUsers(w http.ResponseWriter, req *http.Request) {
	userJson, err := json.Marshal(*userList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Can not marshal Users")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func getUser(w http.ResponseWriter, req *http.Request) (user.User, error) {
	vars := mux.Vars(req)
	idAsString := vars["userId"]

	id, err := strconv.Atoi(idAsString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Can not convert userId to an integer: %s\n", idAsString)
		fmt.Println(err)
		return user.User{}, err
	}

	userObj, err := userList.FromId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return user.User{}, err
	}

	return userObj, nil
}

func HandleGetUser(w http.ResponseWriter, req *http.Request) {
	user, err := getUser(w, req)
	if err != nil {
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Can not marshal User: %s, %s\n", strconv.Itoa(user.Id), user.Name)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func HandlePostUser(w http.ResponseWriter, req *http.Request) {
	user := &user.User{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Can not unmarshl User POST request")
		fmt.Println(err)
		return
	}

	//TODO ADD Id to user
	*userList = append(*userList, *user)

	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Can not marshl User")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJson)
}
