package main

import (
	"fmt"
)

func main() {

	// mux := http.NewServeMux()
	// err := http.ListenAndServe(":4000", mux)
	// if err != nil {

	// }
	v := NewDeck()
	v.GenerateDeck()
	v.Shuffle()
	fmt.Println(v.Deckcard)
}
