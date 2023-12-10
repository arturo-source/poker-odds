package main

import (
	"fmt"
	"time"

	"github.com/arturo-source/poker-engine"
)

type handEquity map[poker.HandKind]uint

type equity struct {
	wins  uint
	ties  uint
	hands handEquity
}

func getCombinations(hands []poker.Cards, board poker.Cards) <-chan poker.Cards {
	combc := make(chan poker.Cards)
	handsJoined := poker.JoinCards(hands...)

	var allCombinations func(currComb poker.Cards, start int, n int)
	allCombinations = func(currComb poker.Cards, start int, n int) {
		if n == 0 {
			combc <- currComb.AddCards(board)
			return
		}

		for i := start; i < poker.MAX_CARDS; i++ {
			if board.HasBit(i) {
				continue
			}
			if handsJoined.HasBit(i) {
				continue
			}

			currComb = currComb.BitToggle(i)
			allCombinations(currComb, i+1, n-1)
			currComb = currComb.BitToggle(i)
		}
	}

	cardsInBoard := board.Count()
	go func() {
		allCombinations(poker.NO_CARD, 0, poker.MAX_CARDS_IN_BOARD-cardsInBoard)
		close(combc)
	}()

	return combc
}

func calculateEquities(hands []poker.Cards, board poker.Cards) (equities map[*poker.Player]equity, nCombinations uint) {
	equities = make(map[*poker.Player]equity)
	players := make([]*poker.Player, 0, len(hands))
	for _, hand := range hands {
		player := &poker.Player{Hand: hand}
		players = append(players, player)
		equities[player] = equity{hands: make(handEquity)}
	}

	for comb := range getCombinations(hands, board) {
		winners := poker.GetWinners(comb, players)
		justOneWinner := len(winners) == 1

		for _, winner := range winners {
			playerEquity := equities[winner.Player]

			playerEquity.hands[winner.HandKind]++
			if justOneWinner {
				playerEquity.wins++
			} else {
				playerEquity.ties++
			}

			equities[winner.Player] = playerEquity
		}

		nCombinations++
	}

	return
}

func main() {
	hands, board, err := parseCommandLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	var start = time.Now()
	equities, nCombinations := calculateEquities(hands, board)
	timeElapsed := time.Since(start)

	printResults(board, equities, nCombinations, timeElapsed)
}
