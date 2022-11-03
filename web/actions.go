package main

const (
	skipmove   = "skipmove"
	takecard   = "takecard"
	playcard   = "playcard"
	newuser    = "newuser"
	reversmove = "reversmove"

	right = 1
	left  = -1
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

func (u *user) WriteAction(s string, v interface{}) {
	u.Actions = append(u.Actions, action{Useraction: s, Data: v})
}

type game struct {
	Id            int    `json:"id"`
	Users         []user `json:"users"`
	DirectionGame int    `josn:"directiongame"`
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
func (g *game) WriteActionForAll(s string, v interface{}) {
	for _, k := range g.Users {
		k.Actions = append(k.Actions, action{Useraction: s, Data: v})
	}
}

func (g *game) WriteActionForAllExceptOne(id int, s string, v interface{}) {
	for _, k := range g.Users {
		if id != k.Id {
			k.Actions = append(k.Actions, action{Useraction: s, Data: v}) //Еблан добавь вывод юзера
		}
	}
}

func (g *game) NextPlayerTurnSkip(id int) int {
	g.MoveInThisTurnFalse(id)
	switch g.DirectionGame {
	case right:
		if id == len(g.Users)-1 {
			id = -2
		}
		g.MoveInThisTurnTrue(id + 2)
		id = id + 2
	case left:
		if id == 0 && len(g.Users) > 2 {
			id = len(g.Users)
		} else if id == 0 && len(g.Users) > 2 {
			id = len(g.Users) + 1
		} else if (id == 0 || id == 1) && len(g.Users) == 2 {
			id += 2
		}
		g.MoveInThisTurnTrue(id - 2)
		id = id - 2
	}
	return id
}
func (g *game) MoveInThisTurnFalse(id int) {
	g.Users[id].MoveInThisTurn = false
}

func (g *game) MoveInThisTurnTrue(id int) {
	g.Users[id].MoveInThisTurn = true
}

func (g *game) NextPlayerTurnNumeric(id int) int {
	g.MoveInThisTurnFalse(id)
	switch g.DirectionGame {
	case right:
		if id == len(g.Users)-1 {
			id = -1
		}
		g.MoveInThisTurnTrue(id + 1)
		id = id + 1
	case left:
		if id == 0 {
			id = len(g.Users)
		}
		g.MoveInThisTurnTrue(id - 1)
		id = id - 1
	}
	return id
}

func (g *game) ClearCard(idcard int) {
	lendec := len(g.CurrentDeck.Deckcard)

	if lendec == 1 {
		lendecdrop := len(g.DropDeck.Deckcard)
		lastcard := g.DropDeck.Deckcard[lendecdrop-1]

		g.DropDeck.Deckcard = append(g.CurrentDeck.Deckcard[:idcard-1], g.CurrentDeck.Deckcard[idcard-1:idcard]...)

		g.CurrentDeck.Deckcard = g.DropDeck.Deckcard

		g.DropDeck.Deckcard = nil

		g.DropDeck.Deckcard = append(g.DropDeck.Deckcard, lastcard)
		return
	}

	if idcard == 0 {
		g.CurrentDeck.Deckcard = append(g.CurrentDeck.Deckcard[1:2], g.CurrentDeck.Deckcard[2:]...)
	} else if idcard == lendec-1 {
		g.CurrentDeck.Deckcard = append(g.CurrentDeck.Deckcard[:idcard-1], g.CurrentDeck.Deckcard[idcard-1:idcard]...)
	} else {
		g.CurrentDeck.Deckcard = append(g.CurrentDeck.Deckcard[:idcard-1], g.CurrentDeck.Deckcard[idcard+1:]...)
	}
}
