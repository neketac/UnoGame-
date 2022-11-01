package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type aplication struct {
	idforgame int
	idforuser int
	Games     []game `json:"games"`
}

func main() {
	app := &aplication{idforgame: 0, Games: make([]game, 0)}
	v := NewDeck()
	v.GenerateDeck()
	v.Shuffle()

	app.Games = append(app.Games, game{Id: 10, Users: make([]user, 0), CurrentDeck: *v})

	// for i := 0; i < 4; i++ {
	// 	app.Games[0].Users = append(app.Games[0].Users, user{Id: i})
	// }
	app.Games[0].CurrentDeck.Deckcard = append(app.Games[0].CurrentDeck.Deckcard[:102], app.Games[0].CurrentDeck.Deckcard[102:104]...)
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/game/{id:[0-9]+}/GetHighCard", app.GetHighCard).Methods("GET")         //GET
	r.HandleFunc("/game/{id:[0-9]+}/GetDeckInHand", app.GetDeckInHand).Methods("GET")     //GET
	r.HandleFunc("/game/{id:[0-9]+}/GetCurrentState", app.GetCurrentState).Methods("GET") //GET
	r.HandleFunc("/CreatGame", app.CreatGame).Methods("POST")                             //POST
	r.HandleFunc("/game/{id:[0-9]+}/CreatUser", app.CreatUser).Methods("GET")             //POST
	// r.HandleFunc("/articles", ArticlesHandler)
	// http.Handle("/", r)
	log.Printf("Запуск сервера на http://127.0.0.1:80")
	http.ListenAndServe("", r)

}
