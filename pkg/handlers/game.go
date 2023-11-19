package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	fp "path/filepath"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/seternate/lanty-api-golang/pkg/game"
	lantyUtil "github.com/seternate/lanty-api-golang/pkg/util"
)

var gameList *game.Games

func InitGamesHandler(games *game.Games) {
	gameList = games
}

func HandleGetGames(w http.ResponseWriter, req *http.Request) {
	gamesJson, err := json.Marshal(*gameList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Can not marshal Games")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(gamesJson)
}

func getGame(w http.ResponseWriter, req *http.Request) (game.Game, error) {
	vars := mux.Vars(req)
	idAsString := vars["gameId"]

	id, err := strconv.Atoi(idAsString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Can not convert gameId to an integer: %s\n", idAsString)
		fmt.Println(err)
		return game.Game{}, err
	}

	gameObj, err := gameList.FromId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return game.Game{}, err
	}

	return gameObj, nil
}

func HandleGetGame(w http.ResponseWriter, req *http.Request) {
	game, err := getGame(w, req)
	if err != nil {
		return
	}

	gameJson, err := json.Marshal(game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Can not marshal Game: %s, %s\n", strconv.Itoa(game.Id), game.Slug)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(gameJson)
}

func getGameBinaryData(w http.ResponseWriter, req *http.Request, directory string) {
	game, err := getGame(w, req)
	if err != nil {
		return
	}

	gameSlug := game.Slug
	filepath := lantyUtil.SearchFileByName(gameSlug, directory)[0]
	filename := fp.Base(filepath)

	if len(filepath) == 0 || len(filename) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Can not find data for Game: %s, %s\n", strconv.Itoa(game.Id), game.Slug)
		return
	}

	contentType, _ := lantyUtil.DetectContentTypeOfFile(filepath)

	if err != nil {
		fmt.Printf("Can not detect content type of file: %s\n", filepath)
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	if req.Method == http.MethodHead {
		fi, err := os.Stat(filepath)
		if err == nil {
			w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
		}
		w.WriteHeader(200)
	} else if req.Method == http.MethodGet {
		http.ServeFile(w, req, filepath)
	}
}

func HandleGetGameFile(w http.ResponseWriter, req *http.Request) {
	getGameBinaryData(w, req, "data")
}

func HandleGetGameIcon(w http.ResponseWriter, req *http.Request) {
	getGameBinaryData(w, req, "icon")
}

func HandleGetGameCover(w http.ResponseWriter, req *http.Request) {
	getGameBinaryData(w, req, "cover")
}
