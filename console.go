package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sort"
	"time"

	"github.com/arturo-source/poker-engine"
)

type Color string

const (
	Reset  Color = "0"
	Black  Color = "30"
	Red    Color = "31"
	Green  Color = "32"
	Yellow Color = "33"
	Blue   Color = "34"
	Purple Color = "35"
	Cyan   Color = "36"
	White  Color = "37"
	Gray   Color = "90"
)

// https://upload.wikimedia.org/wikipedia/commons/b/b8/4coloraces1.jpg
const (
	NoSuit        = Yellow
	SpadesColor   = Black
	ClubsColor    = Green
	HeartsColor   = Red
	DiamondsColor = Blue
)

func colorize(txt string, color Color) string {
	const START_CODE = "\033["
	const END_CODE = "m"

	if disableColor {
		return txt
	}

	return fmt.Sprint(START_CODE, color, END_CODE, txt, START_CODE, Reset, END_CODE)
}

func colorizeCards(cards poker.Cards) string {
	colorizeSuit := func(suit poker.Cards, color Color) string {
		var cardsStr string

		suitStr := poker.SUIT_VALUES[suit]
		cardsSuited := cards & suit

		for cardNum, cardNumStr := range poker.NUMBER_VALUES {
			card := cardsSuited & cardNum
			if card != poker.NO_CARD {
				cardsStr += colorize(cardNumStr+suitStr, color)
			}
		}

		return cardsStr
	}

	spadesStr := colorizeSuit(poker.SPADES, SpadesColor)
	clubsStr := colorizeSuit(poker.CLUBS, ClubsColor)
	heartsStr := colorizeSuit(poker.HEARTS, HeartsColor)
	diamondsStr := colorizeSuit(poker.DIAMONDS, DiamondsColor)

	return fmt.Sprint(spadesStr, clubsStr, heartsStr, diamondsStr)
}

var disableColor bool

func parseCommandLine() (hands []poker.Cards, board poker.Cards, err error) {
	flag.Usage = func() {
		fmt.Println("Available numbers: A K Q J T 9 8 7 6 5 4 3 2")
		fmt.Println("Available suits: s c h d")
		fmt.Println("Example of usage (hands must be last always):\n  poker-odds --board AcTh6h Ah3h KdQd")
		fmt.Println()
		fmt.Println("Available flags:")
		flag.PrintDefaults()
	}

	var boardStr string
	flag.StringVar(&boardStr, "board", "", "The cards with their suits in the board")
	flag.BoolVar(&disableColor, "no-color", false, "Disable color output (raw data for file saving)")

	flag.Parse()

	if runtime.GOOS == "windows" {
		disableColor = true
	}

	// Read all Args input and transform them into cards
	var allCards []poker.Cards
	handsStr := flag.Args()
	if len(handsStr) == 0 {
		flag.Usage()

		err = fmt.Errorf("at least one hand is needed")
		return
	}

	for _, handStr := range handsStr {
		if len(handStr) != 4 {
			err = fmt.Errorf("%s hand is not valid, hands must have 2 cards with suit", colorize(handStr, NoSuit))
			return
		}

		firstCardStr, secondCardStr := handStr[:2], handStr[2:]
		firstCard, secondCard := poker.NewCard(firstCardStr), poker.NewCard(secondCardStr)
		if firstCard == poker.NO_CARD {
			err = fmt.Errorf("%s card (%s hand) is not valid", colorize(firstCardStr, NoSuit), colorize(handStr, NoSuit))
			return
		}
		if secondCard == poker.NO_CARD {
			err = fmt.Errorf("%s card (%s hand) is not valid", colorize(secondCardStr, NoSuit), colorize(handStr, NoSuit))
			return
		}

		hand := poker.JoinCards(firstCard, secondCard)
		hands = append(hands, hand)

		allCards = append(allCards, firstCard, secondCard)
	}

	// Read --board input and transform them into cards
	for i := 0; i < len(boardStr); i += 2 {
		end := i + 2
		if end > len(boardStr) {
			end = len(boardStr)
		}

		cardStr := boardStr[i:end]
		card := poker.NewCard(cardStr)
		if card == poker.NO_CARD {
			err = fmt.Errorf("%s card (%s board) is not valid", colorize(cardStr, NoSuit), colorize(boardStr, NoSuit))
			return
		}

		board = board.AddCards(card)

		allCards = append(allCards, card)
	}

	// Check if any card is repeated
	allCardsJoined := poker.JoinCards(allCards...)
	for _, card := range allCards {
		if !allCardsJoined.CardsArePresent(card) {
			err = fmt.Errorf("card %s is duplicated", colorizeCards(card))
			return
		}

		allCardsJoined = allCardsJoined.QuitCards(card)
	}

	return hands, board, err
}

func printResults(board poker.Cards, equities map[*poker.Player]equity, nCombinations uint, timeElapsed time.Duration) {
	// Sort players to get always winner first and because `equities map[*poker.Player]equity` range do not return always the same order
	var playerPointers = make([]*poker.Player, 0, len(equities))
	for player := range equities {
		playerPointers = append(playerPointers, player)
	}
	sort.Slice(playerPointers, func(i, j int) bool {
		return equities[playerPointers[i]].wins > equities[playerPointers[j]].wins
	})

	// Print board
	if board != poker.NO_CARD {
		fmt.Println()
		fmt.Println(colorize("board", Gray))
		fmt.Println(colorizeCards(board))
	}

	// Print player equities
	fmt.Println()
	pad := 14
	if disableColor {
		pad = 5
	}
	fmt.Printf("%s %*s %*s\n", colorize("hand", Gray), pad, colorize("win", Gray), pad+2, colorize("tie", Gray))
	for _, player := range playerPointers {
		eq := equities[player]
		winPercentage := float64(eq.wins) / float64(nCombinations) * 100
		tiePercentage := float64(eq.ties) / float64(nCombinations) * 100
		fmt.Printf("%s %5.1f%% %5.1f%%\n", colorizeCards(player.Hand), winPercentage, tiePercentage)
	}

	// Print hands equities
	fmt.Println()
	fmt.Printf("%-16s", "")
	pad = 26
	if disableColor {
		pad = 8
	}
	for _, player := range playerPointers {
		fmt.Printf("%*s", pad, colorizeCards(player.Hand))
	}

	fmt.Println()
	handKinds := []poker.HandKind{poker.HIGHCARD, poker.PAIR, poker.TWOPAIR, poker.THREEOFAKIND, poker.STRAIGHT, poker.FLUSH, poker.FULLHOUSE, poker.FOUROFAKIND, poker.STRAIGHTFLUSH, poker.ROYALFLUSH}
	for _, hk := range handKinds {
		fmt.Printf("%-16s", hk)
		for _, player := range playerPointers {
			eq := equities[player]
			handEqPercentage := float64(eq.hands[hk]) / float64(eq.wins+eq.ties) * 100
			if math.IsNaN(handEqPercentage) || handEqPercentage == 0.0 {
				fmt.Printf("%7s%s", "", colorize(".", Gray))
			} else if handEqPercentage < 0.1 {
				fmt.Printf("%4s%s", "", colorize(">0.1", Gray))
			} else {
				fmt.Printf("%7.1f%%", handEqPercentage)
			}
		}
		fmt.Println()
	}

	// Print program stats
	fmt.Println()
	fmt.Println(colorize(fmt.Sprintf("%d combinations calculated in %s", nCombinations, timeElapsed), Gray))
}
