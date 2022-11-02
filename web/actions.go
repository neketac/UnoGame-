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

func (g *game) SearchIdUser(id int) (int, bool) {
	for i, k := range g.Users {
		if id == k.Id {
			return i, true
		}
	}
	return 0, false
}

func (g *user) WriteAction(s string, v interface{}) {
	g.Actions = append(g.Actions, action{Useraction: s, Data: v})
}
