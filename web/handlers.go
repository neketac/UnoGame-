package main

import (
	help "UnoGame/Help"
	"encoding/json"
	"fmt"
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

			ap.Games[id].Users[g].Deckinhand.AddingCard(ap.Games[id].CurrentDeck.Deckcard[idcard])

			ap.Games[id].CurrentDeck.ClearCard(idcard)
		}
	}
	//Первая карта на столе
	ap.Games[id].DropDeck.AddingCard(ap.Games[id].CurrentDeck.Deckcard[0])
	ap.Games[id].CurrentDeck.ClearCard(0)

	ap.Games[id].Users[rand.Intn(len(ap.Games[id].Users))].FirstMove = true //Определение первого кто ходит
	ap.Games[id].Users[rand.Intn(len(ap.Games[id].Users))].MoveInThisTurn = true

	ap.Games[id].GameStart = true //Обьявляет что игра началась

	help.RenderAndWrite(w, ap.Games[id].Users)
}

func (ap *aplication) GetHighCard(w http.ResponseWriter, req *http.Request) { //Получение верхней карты при доборе
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	if len(ap.Games[id].Users) == 0 {
		http.Error(w, "Userov Net Eblan", http.StatusBadRequest)
		return
	}

	type userid struct {
		Id int
	}

	var user userid

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\nchto za kal v uzere", err.Error()), http.StatusBadRequest)
		return
	}

	iduser, _ := ap.Games[id].SearchIdUser(user.Id)

	type UserResponse struct {
		Id        int      `json:"id"`
		Name      string   `json:"name"`
		Permitted bool     `json:"permitted"`
		Actions   []action `json:"action"`
	}
	if !ap.Games[id].Users[iduser].MoveInThisTurn {
		notpermited := UserResponse{Id: ap.Games[id].Users[iduser].Id, Name: ap.Games[id].Users[iduser].Name, Permitted: false}
		help.RenderAndWrite(w, notpermited)
	}

	type UserForActions struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	userforactions := UserForActions{Id: ap.Games[id].Users[iduser].Id, Name: ap.Games[id].Users[iduser].Name}

	v := ap.Games[id].CurrentDeck.Deckcard[0]
	ap.Games[id].CurrentDeck.ClearCard(0)
	ap.Games[id].Users[iduser].Deckinhand.AddingCard(v) //Добавление карты в руку пользователю
	ap.Games[id].Users[iduser].WriteAction(takecard, v) //Записываем действие "Взял карту"
	for iduser, k := range ap.Games[id].Users {
		if user.Id != k.Id {
			ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: takecard, Data: userforactions}) //Еблан добавь вывод юзера
		}
	}

	resuser := UserResponse{Id: ap.Games[id].Users[iduser].Id, Name: ap.Games[id].Users[iduser].Name, Permitted: true, Actions: ap.Games[id].Users[iduser].Actions}

	// ap.Games[id].Users[iduser].Actions = append(ap.Games[id].Users[iduser].Actions, action{Useraction: })
	help.RenderAndWrite(w, resuser)
}

func (ap *aplication) GetCurrentState(w http.ResponseWriter, req *http.Request) { //Отправка инфы по игре
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	type userid struct {
		Id int
	}

	if len(ap.Games[id].Users) == 0 {
		http.Error(w, "Userov Net Eblan", http.StatusBadRequest)
		return
	}

	var user userid
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\nchto za kal v uzere", err.Error()), http.StatusBadRequest)
		return
	}

	useridget, _ := ap.Games[id].SearchIdUser(user.Id)

	type Response struct {
		Id      int      `json:"id"`
		Actions []action `json:"actions"`
	}
	ResponseUser := Response{Id: ap.Games[id].Users[useridget].Id, Actions: ap.Games[id].Users[useridget].Actions}

	help.RenderAndWrite(w, ResponseUser)
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

	help.RenderAndWrite(w, id)
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
		ap.Games[id].Users[iduser].WriteAction(newuser, user)
	}

	UserArray := Response{CreatIdUser: ap.idforuser, Users: make([]UserResponse, 0)}

	for _, k := range ap.Games[id].Users {
		UserArray.Users = append(UserArray.Users, UserResponse{Id: k.Id, Name: k.Name})
	}

	help.RenderAndWrite(w, UserArray)
}

func (ap *aplication) GetListGame(w http.ResponseWriter, req *http.Request) {
	if len(ap.Games) == 0 {
		http.Error(w, "Igor net Daunich", http.StatusBadRequest)
	}

	type Respone struct {
		GamesArray []int
	}

	ResponeGames := Respone{GamesArray: make([]int, 0)}
	for _, k := range ap.Games {
		ResponeGames.GamesArray = append(ResponeGames.GamesArray, k.Id)
	}

	help.RenderAndWrite(w, ResponeGames)
}

func (ap *aplication) GetListUsers(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	if len(ap.Games[id].Users) == 0 {
		http.Error(w, "Igrokov net Daunich", http.StatusBadRequest)
	}

	type ResUser struct {
		Id   int
		Name string
	}

	type Respone struct {
		UsersArray []ResUser
	}

	ResponseUser := Respone{UsersArray: make([]ResUser, 0)}
	for _, k := range ap.Games[id].Users {
		ResponseUser.UsersArray = append(ResponseUser.UsersArray, ResUser{Id: k.Id, Name: k.Name})
	}

	help.RenderAndWrite(w, ResponseUser)
}

func (ap *aplication) PlayCard(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	type Request struct {
		Id       int
		Dropcard card
	}

	var requestcard Request
	err := json.NewDecoder(req.Body).Decode(&requestcard)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\nchto za karta uebumba?", err.Error()), http.StatusBadRequest)
		return
	}

	userid, _ := ap.Games[id].SearchIdUser(requestcard.Id)

	type Respone struct {
		Permitted bool `json:"permitted"`
	}

	responsebool := Respone{}

	if !ap.Games[id].Users[userid].MoveInThisTurn {
		help.RenderAndWrite(w, responsebool)
	}

	equals := CheckCard(requestcard.Dropcard, ap.Games[id].DropDeck.Deckcard)

	if requestcard.Dropcard.CardType == numeric && equals {
		var idcard int
		for i, k := range ap.Games[id].Users[userid].Deckinhand.Deckcard {
			if k == requestcard.Dropcard {
				idcard = i
			}
		}
		ap.Games[id].Users[userid].Deckinhand.ClearCard(idcard)
		ap.Games[id].DropDeck.AddingCard(requestcard.Dropcard)

	}

}

//методы: Вернуть список комнат, Вернуть список игроков ИД
//Переделать переддачу черех экшоны
