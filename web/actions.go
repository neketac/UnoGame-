package main

const (
	skipmove = "skipmove"
	takecard = "takecard"
	putcard  = "putcard"
	newuser  = "newuser"
)

type action struct {
	Useraction string      `json:"useraction"`
	Data       interface{} `json:"card"`
}
type user struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Actions        []action `json:"actions"`
	Deckinhand     deck     `json:"deckinhand"`
	FirstMove      bool     `json:"firstmove"`
	MoveInThisTurn bool     `json:"moveinthisturn"`
}
type game struct {
	Id            int    `json:"id"`
	Users         []user `json:"users"`
	DirectionGame bool   `josn:"directiongame"`
	CurrentDeck   deck   `json:"currentdeck"`
	DropDeck      deck   `json:"dropdeck"`
	GameStart     bool   `json:"gamestart"`
}
