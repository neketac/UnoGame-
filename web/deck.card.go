package main

import (
	"math/rand"
	"time"
)

const (
	//Типы карт
	numeric  = "numeric"
	revers   = "revers"
	drawtwo  = "drawtwo"
	drawfore = "drawfore"
	skip     = "skip"
	//Цвета карт
	blue   = "blue"
	green  = "green"
	red    = "red"
	yellow = "yellow"
	wild   = "wild"
)

type cardtype string

type equalscard struct {
	equalscolor  bool
	equalsnumber bool
}

type card struct {
	CardType cardtype `json:"cardType"`
	Number   int      `json:"number"`
	Color    string   `json:"color"`
}

type deck struct {
	Deckcard []card `json:"deckcard"`
}

func NewDeck() *deck {
	deck := deck{Deckcard: make([]card, 0)}
	return &deck
}

func (d *deck) GenerateDeck() {
	cl := [...]string{blue, green, red, yellow}
	tc := [...]cardtype{revers, drawtwo, skip}
	for _, k := range cl {
		d.Deckcard = append(d.Deckcard, card{CardType: numeric, Number: 0, Color: k})
		for i := 1; i <= 9; i++ {
			d.Deckcard = append(d.Deckcard, card{CardType: numeric, Number: i, Color: k})
			d.Deckcard = append(d.Deckcard, card{CardType: numeric, Number: i, Color: k})
		}
	}

	for _, t := range tc {
		for _, k := range cl {
			d.Deckcard = append(d.Deckcard, card{CardType: t, Color: k})
			d.Deckcard = append(d.Deckcard, card{CardType: t, Color: k})
		}
	}

	for i := 0; i < 4; i++ {
		d.Deckcard = append(d.Deckcard, card{CardType: drawfore, Color: wild})
	}
}

func (d *deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Deckcard),
		func(i, j int) { d.Deckcard[i], d.Deckcard[j] = d.Deckcard[j], d.Deckcard[i] })
}

func CheckCard(c card, dropdeck []card) bool {
	if c.Color == wild {
		return true
	}

	lendropdeck := len(dropdeck)
	equals := equalscard{}

	switch c.Color {
	case dropdeck[lendropdeck-1].Color:
		equals.equalscolor = true
	default:
		equals.equalscolor = false
	}

	switch c.Number {
	case dropdeck[lendropdeck-1].Number:
		equals.equalscolor = true
	default:
		equals.equalscolor = false
	}

	if (c.CardType == drawtwo || c.CardType == skip || c.CardType == revers) && equals.equalscolor {
		return true
	}

	if c.CardType == numeric && (equals.equalscolor || equals.equalsnumber) {
		return true
	}

	return false
}

// arr []card, idcard int
func (arr *deck) ClearCard(idcard int) {
	len := len(arr.Deckcard)
	if idcard == 0 {
		arr.Deckcard = append(arr.Deckcard[1:2], arr.Deckcard[2:]...)
	} else if idcard == len-1 {
		arr.Deckcard = append(arr.Deckcard[:idcard-1], arr.Deckcard[idcard-1:idcard]...)
	} else {
		arr.Deckcard = append(arr.Deckcard[:idcard-1], arr.Deckcard[idcard+1:]...)
	}
}

func (arr *deck) AddingCard(cardinarr card) {
	arr.Deckcard = append(arr.Deckcard, cardinarr)
}
