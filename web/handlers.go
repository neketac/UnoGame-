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

func (ap *aplication) GetDeckInHand(w http.ResponseWriter, req *http.Request) { // Раздача карт
	rand.Seed(time.Now().UnixNano())

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	for g := range ap.Games[id].Users {
		for i := 0; i < 8; i++ {
			len := len(ap.Games[id].CurrentDeck.Deckcard)
			idcard := rand.Intn(len)

			ap.Games[id].Users[g].Deckinhand.Deckcard = append(ap.Games[id].Users[g].Deckinhand.Deckcard, ap.Games[id].CurrentDeck.Deckcard[idcard])

			if idcard == 0 {
				ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...)
			} else if idcard == len-1 {
				ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[:idcard-1], ap.Games[id].CurrentDeck.Deckcard[idcard-1:idcard]...)
			} else {
				ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[:idcard-1], ap.Games[id].CurrentDeck.Deckcard[idcard+1:]...)
			}
		}
	}

	ap.Games[id].Users[rand.Intn(len(ap.Games[id].Users))].FirstMove = true

	js, err := json.Marshal(ap.Games[id].Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ap *aplication) GetHighCard(w http.ResponseWriter, req *http.Request) { //Получение верхней карты при доборе
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	userid, _ := strconv.Atoi(req.URL.Query().Get("user"))

	log.Println(userid)

	v := ap.Games[id].CurrentDeck.Deckcard[0] //Надо id игры и id пользователя соотвественно
	ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...)

	ap.Games[id].Users[userid].Deckinhand.Deckcard = append(ap.Games[id].Users[userid].Deckinhand.Deckcard, v)
	ap.Games[id].Users[userid].Actions = append(ap.Games[id].Users[userid].Actions, action{Useraction: takecard, Card: v})

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
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	js, err := json.Marshal(ap.Games[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	for g := range ap.Games[id].Users {
		ap.Games[id].Users[g].Actions = nil
	}
	//игрок  ИД Количество карт в руке
	//масив действий игроков которые произошли с момента последнего запроса записывать действия игроков в масив
	//

}
