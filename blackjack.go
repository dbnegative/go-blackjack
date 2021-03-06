package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var cardtypes = map[string]int{
	"Ace":   11,
	"King":  10,
	"Queen": 10,
	"Jack":  10,
	"10":    10,
	"9":     9,
	"8":     8,
	"7":     7,
	"6":     6,
	"5":     5,
	"4":     4,
	"3":     3,
	"2":     2,
}
var suite = map[int]string{
	0: "\u2662",
	1: "\u2661",
	2: "\u2667",
	3: "\u2664",
}

var keys = []string{
	"g",
	"h",
	"y",
	"n",
}

//Game of blackjack
type Game struct {
	decks   []Deck
	players []Player
	turn    Turn
	running bool
}

//Player has a turn, cards and a score
type Player struct {
	Name         string
	Cards        []Card
	Score        int
	IsTurn       bool
	HasAce       bool
	AceAmount    int
	ConvertedAce int
}

//Turn that can be taken on a turn
type Turn struct {
	numberofturns int
	whosturn      int
}

//Deck of crades
type Deck struct {
	Cards []Card
	Order []int
}

//Card that can be played
type Card struct {
	Type      string
	Value     int
	Suite     string
	Available bool
}

var gscore1 = 0
var gscore2 = 0

func main() {
	anothergame := true
	for anothergame {
		rand.Seed(time.Now().UTC().UnixNano())
		bjgame := new(Game)
		bjgame.InitGame()

		for bjgame.running {
			for i := 0; bjgame.turn.numberofturns >= i; i++ {
				activeplayer := bjgame.turn.whosturn
				turnactive := true
				for turnactive {
					//clearscreen()
					bjgame.checkstate()
					if !bjgame.running {
						break
					}
					if bjgame.players[activeplayer].Name == "Dealer" {
						bjgame.strategy()
					} else {
						fmt.Println("", bjgame.players[activeplayer].Name, "Turn!")
						fmt.Println("+++++++++++++++++++++++++++++++++++++++")
						printgfxcards(bjgame.players[activeplayer].Cards, false)
						fmt.Println("+++++++++++++++++++++++++++++++++++++++")
						fmt.Println("", bjgame.players[activeplayer].Name, " total is: ", bjgame.players[activeplayer].Score)
						fmt.Println("+++++++++++++++++++++++++++++++++++++++")
						printgfxcards(bjgame.players[activeplayer+1].Cards, true)
						fmt.Println("+++++++++++++++++++++++++++++++++++++++")
						bjgame.action(activeplayer, getresponse(" (g) Hit, (h) Hold: "))

					}
					if !bjgame.players[activeplayer].IsTurn { //has player ended turn
						if bjgame.turn.whosturn == 0 {
							bjgame.turn.whosturn = 1
						} else {
							bjgame.turn.whosturn = 0
						}
						bjgame.players[activeplayer].IsTurn = true //reset
						turnactive = false                         //end turn
					}
					fmt.Println("A1:", bjgame.players[0].Score)
				}
			}
			fmt.Println("B1:", bjgame.players[0].Score)
			bjgame.whowins()
			bjgame.running = false
			//debug
			//fmt.Println(bjgame.players)
		}
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		resp := getresponse("Want to play again? (y/n):")
		resp = strings.ToLower(resp)
		if resp == "n" {
			anothergame = false
		}
	}
}

//InitGame inits a new game
func (game *Game) InitGame() *Game {
	game.decks = make([]Deck, 1)
	game.decks[0].NewDeck()
	game.players = make([]Player, 2)
	game.players[0].NewPlayer("Player", game.decks[0])
	game.players[1].NewPlayer("Dealer", game.decks[0])
	game.players[0].setscore()
	game.players[1].setscore()
	game.players[0].IsTurn = true
	game.players[1].IsTurn = true
	game.turn.newturn()
	game.turn.numberofturns = 1 //each
	game.running = true
	return game
}

//NewDeck init's a deck
func (deck *Deck) NewDeck() *Deck {
	deck.Cards = make([]Card, 52)
	deck.Order = make([]int, 52)
	deck.createdeck(1)
	deck.randomorder()
	return deck
}
func (turn *Turn) newturn() *Turn {
	turn.numberofturns = 1
	turn.whosturn = 0
	return turn

}

//NewPlayer init's a new player
func (player *Player) NewPlayer(name string, deck Deck) *Player {
	player.Cards = make([]Card, 2)
	player.HasAce = false //player has no cards yet
	player.AceAmount = 0  //player has no cards yet
	player.assigncards(deck, 1)
	player.IsTurn = false
	player.Score = 0
	player.Name = name
	player.ConvertedAce = 0
	return player
}

func (player *Player) hit(deck Deck) *Player {
	rand.Seed(48)
	//assign 2 cards
	assigned := false

	for i := 0; 1 > i; i++ {
		for !assigned {
			orderindex := deck.Order[rand.Intn(51)]
			if deck.Cards[orderindex].Available {
				player.Cards = append(player.Cards, deck.Cards[orderindex])
				deck.Cards[orderindex].Available = false
				assigned = true
			}
		}
		assigned = false
	}
	return player
}
func (player *Player) hasace() *Player {
	for i := range player.Cards {
		if player.Cards[i].Type == "Ace" {
			player.HasAce = true
			player.AceAmount = player.AceAmount + 1
		}
	}
	return player
}

func (player *Player) hold() *Player {
	player.IsTurn = false
	return player
}

func (player *Player) setscore() *Player {
	player.Score = 0

	for i := 0; len(player.Cards) > i; i++ {
		player.Score = player.Score + player.Cards[i].Value
	}
	//fmt.Println("score is", player.Score)
	if player.ConvertedAce > 0 {
		for i := 0; player.ConvertedAce > i; i++ {
			player.Score = player.Score - 10
		}
	}
	return player
}
func (deck *Deck) createdeck(numberofdecks int) *Deck {
	count := 0
	//cardtypes =13
	for v, t := range cardtypes {
		//suite = 4
		for _, n := range suite {
			if count < 52 {
				deck.Cards[count].Suite = n
				deck.Cards[count].Type = v
				deck.Cards[count].Value = t
				deck.Cards[count].Available = true
			}
			count++
		}
	} //suite
	return deck
}

//RandomOrder returns a deck of shuffled cards
func (deck *Deck) randomorder() *Deck {
	//rand.Seed(47)
	fmt.Println("Shuffeling Deck!")
	deck.Order = rand.Perm(51)
	return deck
}

//AssignCards assigns cards to players
func (player *Player) assigncards(deck Deck, amount int) (*Player, Deck) {
	//fmt.Println("working with this deck", deck)
	//fmt.Println(" ")
	if amount < 0 {
		amount = 1
	}
	//rand.Seed(time.Second(time.Now()))
	//assign 2 cards
	assigned := false

	for i := 0; amount >= i; i++ {
		for !assigned {
			orderindex := deck.Order[rand.Intn(51)]
			if deck.Cards[orderindex].Available {
				//debug
				//player.Cards = []Card{Card{"Ace", 11, "\u2661", true}, Card{"8", 8, "\u2661", true}}
				player.Cards[i] = deck.Cards[orderindex]
				deck.Cards[orderindex].Available = false
				assigned = true
			}

		}
		assigned = false
	}

	player.hasace()
	//fmt.Println("Player :", player, "\n\n")
	return player, deck
}
func (game *Game) action(playerid int, input string) *Game {
	//fmt.Println("Choice is: ", input)
	switch {
	case input == "g":
		game.players[playerid].hit(game.decks[0])
		game.players[playerid].setscore()
	case input == "h":
		game.players[playerid].hold()
		game.players[playerid].setscore()
	}

	return game
}

func getresponse(instruction string) string {
	reader := bufio.NewReader(os.Stdin)
	isvalid := false
	var text string
	for !isvalid {
		fmt.Print(instruction)
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.Trim(text, "\n"))
		if checkresponse(keys, text) {
			isvalid = true
			//fmt.Println("Text = ", text)
			return text
		}
		fmt.Println("Invalid Option, Please try again!")

	}
	return text
}

func checkresponse(validinput []string, input string) bool {
	for i := range validinput {
		//fmt.Println("Input=", input, " Validinput @ ", validinput[i])
		if input == validinput[i] {
			return true
		}
	}
	return false
}

func (game *Game) whowins() {
	player1score := game.players[0].Score
	player2score := game.players[1].Score
	player1name := game.players[0].Name
	player2name := game.players[1].Name

	fmt.Println("Score1 =", player1score)
	fmt.Println("Score2 =", player2score)

	switch {
	case player1score == 21 && game.players[0].HasAce && len(game.players[0].Cards) < 3:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println("BLACKJACK!! " + player1name + " Has Won!!! ")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		gscore1 = gscore1 + 1
		printgfxcards(game.players[1].Cards, false)
	case player2score == 21 && game.players[1].HasAce && len(game.players[1].Cards) < 3:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println("BLACKJACK!! " + player2name + " Has Won!!! ")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		printgfxcards(game.players[1].Cards, false)
		gscore2 = gscore2 + 1
	case player1score > 21:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println(player1name, "is Bust!! "+player2name+" Has Won!! "+player1name+" Score:", player1score)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		printgfxcards(game.players[1].Cards, false)
		gscore2 = gscore2 + 1
	case player2score > 21:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println(player2name, "is Bust!! "+player1name+" Has Won!! "+player2name+" Score:", player2score)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		gscore1 = gscore1 + 1
		printgfxcards(game.players[1].Cards, false)
	case player1score > player2score:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println(player1name, "has Won!! "+player2name+" Score:", player2score)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		printgfxcards(game.players[1].Cards, false)
		gscore1 = gscore1 + 1
	case player2score > player1score:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println(player2name, "has Won!! Score:", player2score)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		printgfxcards(game.players[1].Cards, false)
		gscore2 = gscore2 + 1
	case player1score == player2score:
		printgfxcards(game.players[0].Cards, false)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		fmt.Println("It's a Draw :'(")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++")
		printgfxcards(game.players[1].Cards, false)

	}
}

//TODO

func (game *Game) checkstate() *Game {

	bjprompt(game.players[0].Name, game.players[0].Score, game.players[1].Name, game.players[1].Score)
	for i := range game.players {
		//dealt a natural blackjack i.e Ace and 10 or above?
		if game.players[i].HasAce && game.players[i].Score == 21 && len(game.players[i].Cards) < 3 {
			game.running = false
			return game
		}
		if game.players[i].Score > 21 {
			if game.players[i].HasAce && game.players[i].AceAmount > 0 {
				game.players[i].Score = game.players[i].Score - 10              //give the player ace's lower value
				game.players[i].ConvertedAce = game.players[i].ConvertedAce + 1 //reduce the amount so check is not triggered
				game.players[i].AceAmount = game.players[i].AceAmount - 1
				return game
				break
			}
			game.running = false
			return game
		}
	}
	return game
}
func (game *Game) strategy() *Game {
	if game.players[1].Score <= 16 {
		game.action(1, "g")
		fmt.Println("Dealer takes another card")
	}
	if game.players[1].Score >= 17 && game.players[1].Score <= 21 {
		game.action(1, "h")
		fmt.Println("Dealer is holding")

	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
	return game
}

func bjprompt(name1 string, score1 int, name2 string, score2 int) {
	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
	fmt.Println("", name1, ":", gscore1, " vs ", name2, ":", gscore2)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
}

func printgfxcards(cards []Card, hide bool) {
	line1 := make([]string, 5)
	line2 := make([]string, 5)
	if !hide {
		for i := range cards {
			ctype1 := strings.Split(cards[i].Type, "")[0]
			if _, err := strconv.Atoi(ctype1); err == nil {
				ctype1 = cards[i].Type
			}
			line1[0] = line1[0] + "  ----- "
			if len(ctype1) > 1 {
				line1[1] = line1[1] + " | " + ctype1 + "  |"
			} else {
				line1[1] = line1[1] + " | " + ctype1 + "   |"
			}
			line1[2] = line1[2] + (" |     |")
			//line[3] = line[3] + " |  " + strconv.Itoa(cards[i].Value) + "  |"
			line1[3] = line1[3] + " |   " + cards[i].Suite + " |"

			line1[4] = line1[0]
		}
		for i := range line1 {
			fmt.Println(line1[i])
		}

	} else {

		ctype2 := strings.Split(cards[0].Type, "")[0]
		if _, err := strconv.Atoi(ctype2); err == nil {
			ctype2 = cards[0].Type
		}
		line2[0] = line2[0] + "  -----   -----"
		if len(ctype2) > 1 {
			line2[1] = line2[1] + " | " + ctype2 + "  | | * \u0398 |"
		} else {
			line2[1] = line2[1] + " | " + ctype2 + "   | | * \u0398 |"
		}
		line2[2] = line2[2] + (" |     | |  *  |")
		//line[3] = line[3] + " |  " + strconv.Itoa(cards[i].Value) + "  |"
		line2[3] = line2[3] + " |   " + cards[0].Suite + " | | \u0394 * |"

		line2[4] = line2[0]

		for i := range line2 {
			fmt.Println(line2[i])
		}

	}
}
func clearscreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

/*
1  -----
2 | \ j |
3 |  \  |
4 | b \ |
5  -----
*/
