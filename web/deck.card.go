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
