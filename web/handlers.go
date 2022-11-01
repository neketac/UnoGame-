package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (ap *aplication) GetDeckInHand(w http.ResponseWriter, req *http.Request) { // Раздача карт
	rand.Seed(time.Now().UnixNano())

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	for g := range ap.Games[id].Users { //Раздача карт в руки
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
	//Первая карта на столе
	ap.Games[id].DropDeck.Deckcard = append(ap.Games[id].DropDeck.Deckcard, ap.Games[id].CurrentDeck.Deckcard[0])
	ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...)

	ap.Games[id].Users[rand.Intn(len(ap.Games[id].Users))].FirstMove = true //Определение первого кто ходит
	ap.Games[id].Users[rand.Intn(len(ap.Games[id].Users))].MoveInThisTurn = true

	ap.Games[id].GameStart = true //Обьявляет что игра началась

	js, err := json.MarshalIndent(ap.Games[id].Users, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ap *aplication) GetHighCard(w http.ResponseWriter, req *http.Request) { //Получение верхней карты при доборе
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	type userid struct {
		Id int
	}

	var user userid

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log.Printf(strconv.FormatInt(int64(user.id), 10))
	// , err := strconv.Atoi(req.URL.Query().Get("user"))
	// if err != nil {
	// 	http.Error(w, "chto za kal v uzere", http.StatusInternalServerError)
	// 	return
	// }
	// if user.id == 0 {
	// 	http.Error(w, "User ukazhi daun", http.StatusBadRequest)
	// 	return
	// }

	var iduser int
	for i, k := range ap.Games[id].Users {
		if (user.Id) == k.Id {
			iduser = i
		}
	}

	// log.Println(userid)

	v := ap.Games[id].CurrentDeck.Deckcard[0]                                                                                    //Получение верхней карты
	ap.Games[id].CurrentDeck.Deckcard = append(ap.Games[id].CurrentDeck.Deckcard[1:2], ap.Games[id].CurrentDeck.Deckcard[2:]...) //Удаление верхней карты их текущей деки

	ap.Games[id].Users[iduser].Deckinhand.Deckcard = append(ap.Games[id].Users[iduser].Deckinhand.Deckcard, v)             //Добавление карты в руку пользователю
	ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: takecard, Data: v}) //Записываем действие "Взял карту"
	for iduser, k := range ap.Games[id].Users {
		if user.Id != k.Id {
			ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: takecard}) //Еблан добавь вывод юзера
		}
	}

	// ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: })
	js, err := json.MarshalIndent(ap.Games[id].Users[iduser].Actions, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ap *aplication) PutCardsPlayedInDeck(w http.ResponseWriter, req *http.Request) {

}

func (ap *aplication) GetCurrentState(w http.ResponseWriter, req *http.Request) { //Отправка инфы по игре
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	js, err := json.MarshalIndent(ap.Games[id], "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	for g := range ap.Games[id].Users { //удаляем действия пользователей
		ap.Games[id].Users[g].Actions = nil
	}
}

func (ap *aplication) CreateGame(w http.ResponseWriter, req *http.Request) {
	ap.idforgame++

	decks := NewDeck()
	decks.GenerateDeck()
	decks.Shuffle()

	ap.Games = append(ap.Games, game{Id: ap.idforgame,
		Users:       make([]user, 0),
		CurrentDeck: *decks,
		DropDeck: deck{
			Deckcard: make([]card, 0)}})

	type ResponseId struct {
		Id int `json:"id"`
	}
	id := ResponseId{Id: ap.idforgame}

	js, err := json.MarshalIndent(id, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func (ap *aplication) CreateUser(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	ap.idforuser++

	ap.Games[id].Users = append(ap.Games[id].Users, user{Id: ap.idforuser, Actions: make([]action, 0), Deckinhand: deck{Deckcard: make([]card, 0)}})

	type UserResponse struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type Response struct {
		CreatIdUser int
		Users       []UserResponse `json:"User"`
	}

	user := UserResponse{}
	for _, k := range ap.Games[id].Users {
		if ap.idforuser == k.Id {
			user.Id = ap.idforuser
			user.Name = k.Name
		}
	}

	for iduser := range ap.Games[id].Users {
		ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: newuser, Data: user})
	}

	UserArray := Response{CreatIdUser: ap.idforuser, Users: make([]UserResponse, 0)}

	for _, k := range ap.Games[id].Users {
		UserArray.Users = append(UserArray.Users, UserResponse{Id: k.Id, Name: k.Name})
	}

	js, err := json.MarshalIndent(UserArray, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

//методы: Вернуть список комнат, Вернуть список игроков ИД
//Переделать переддачу черех экшоны
