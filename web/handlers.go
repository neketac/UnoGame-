package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (ap *aplication) GetDeckInHand(w http.ResponseWriter, req *http.Request) {
	rand.Seed(time.Now().UnixNano())

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	for _, k := range ap.Games[id].Users {
		len := len(ap.Games[id].CurrentDeck.Deckcard)
		idcard := rand.Intn(len)

		k.Deckinhand.Deckcard = append(k.Deckinhand.Deckcard, ap.Games[id].CurrentDeck.Deckcard[idcard])

		if idcard == 0 {
			ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...)
		} else if idcard == len-1 {
			ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[:idcard-1], ap.Games[id].CurrentDeck.Deckcard[idcard-1:idcard]...)
		} else {
			ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[:idcard-1], ap.Games[id].CurrentDeck.Deckcard[idcard+1:]...)
		}

	}
}

func (ap *aplication) GetHighCard(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	userid := req.URL.Query().Get("user")
	log.Println(userid)
	v := ap.Games[id].CurrentDeck.Deckcard[0] //Надо id игры и id пользователя соотвественно
	ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...)

	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ap *aplication) PutCardsPlayedInDeck(w http.ResponseWriter, req *http.Request) {

}

func (ap *aplication) GetCurrentState(w http.ResponseWriter, req *http.Request) {
	//игрок  ИД Количество карт в руке
	//масив действий игроков которые произошли с момента последнего запроса записывать действия игроков в масив
	//

}
