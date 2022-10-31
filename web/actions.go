package main

const (
	skipmove = "skipmove"
	takecard = "takecard"
	putcard  = "putcard"
)

type action struct {
	Useraction string `json:"useraction"`
	Card       card   `json:"card"`
}
type user struct {
	Id         int      `json:"id"`
	Actions    []action `json:"actions"`
	Deckinhand deck     `json:"deckinhand"`
}
type game struct {
	Id          int    `json:"id"`
	Users       []user `json:"users"`
	CurrentDeck deck   `json:"currentdeck"`
	DropDeck    deck   `json:"dropdeck"`
}
