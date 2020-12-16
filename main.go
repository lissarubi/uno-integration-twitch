package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"io/ioutil"
	"math/rand"
	"github.com/joho/godotenv"
	"github.com/jrm780/gotirc"
)

func remove(s []string, i int) []string {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

type Player struct {
	name string
	cards []string
}

func generateCard() string {

	cards := [53]string{"red 1","red 2","red 3","red 4","red 5","red 6","red 7","red 8","red 9","blue 1","blue 2","blue 3","blue 4","blue 5","blue 6","blue 7","blue 8","blue 9", "green 1","green 2","green 3","green 4","green 5","green 6","green 7","green 8","green 9", "yellow 1","yellow 2","yellow 3","yellow 4","yellow 5","yellow 6","yellow 7","yellow 8","yellow 9", "wild 4", "red reverse", "yellow reverse", "blue reverse", "green reverse", "red +2", "green +2", "yellow +2", "blue +2"}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(cards))
	pick := cards[randomIndex]

	return pick
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initUNO(messageString []string, tags map[string]string, channel string) {
	channelClean := strings.ReplaceAll(channel, "#", "")
	if messageString[0] == "!startuno" && tags["display-name"] == channelClean{

		d1 := []byte("waiting")
		err := ioutil.WriteFile("unoStarted.txt", d1, 0644)
		check(err)

	}
}

func endUno(messageString []string, tags map[string]string, channel string, players []Player, unoStarted string, client *gotirc.Client) []Player {
	channelClean := strings.ReplaceAll(channel, "#", "")
	if messageString[0] == "!enduno" && tags["display-name"] == channelClean{

		d1 := []byte("started")
		err := ioutil.WriteFile("unoStarted.txt", d1, 0644)
		check(err)

		for i := 0;  i < len(players); i++{

			for j := 0; j < 7; j++{
				pick := generateCard()
				players[i].cards = append(players[i].cards, pick)
			}
			
			cards := strings.Join(players[i].cards, ", ")
			client.Whisper(players[i].name, players[i].name + " suas cartas são: " + cards)
		}

		return players

	}
	return players
}

func checkUnoStarted() string{

	if fileExists("unoStarted.txt"){
		dat, err := ioutil.ReadFile("unoStarted.txt")
		check(err)
		unoStarted := string(dat)

		return unoStarted
	} else{
		return "false"
	}
}

func main() {

	var players []Player

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	user := os.Getenv("USER")
	channel := os.Getenv("CHANNEL")


	options := gotirc.Options{
		Host:     "irc.chat.twitch.tv",
		Port:     6667,
		Channels: []string{"#" + channel},
	}

	client := gotirc.NewClient(options)

	// Whenever someone sends a message, log it
	var stack []string
	currentPlayer := 0
	client.OnChat(func(channel string, tags map[string]string, msg string) {

		fmt.Println(msg)

		messageString := strings.Split(msg, " ")

		unoState := checkUnoStarted()

		initUNO(messageString, tags, channel)

		fmt.Println(unoState)
		if messageString[0] == "!enter" && unoState == "waiting"{
			var player Player
			player.name = tags["display-name"]

			players = append(players, player)
		}

		fmt.Println("Players (main): ", players)

		if messageString[0] == "!play" && len(messageString) == 3 && unoState == "started"{
			player := players[currentPlayer]
			fmt.Println(player)
			if player.name ==  tags["display-name"]{
				card := strings.ReplaceAll(messageString[1], " ", "") + " " + strings.ReplaceAll(messageString[2], " ", "")
				fmt.Println(card)
				fmt.Println(currentPlayer)

				for i := 0; i < len(player.cards); i++{
					if (card == player.cards[i]){
						stack = append(stack, card)
						players[currentPlayer].cards = remove(player.cards, i)

						cards := strings.Join(player.cards, ", ")
						client.Whisper(player.name, player.name + " suas cartas agora são: " + cards)

						if currentPlayer == len(players) - 1{
							currentPlayer = 0
						}else{
							currentPlayer++
						}
						fmt.Println(stack)
					}
				}
			}
	}
		players = endUno(messageString, tags, channel, players, unoState, client)

	})

	client.Connect(user, token)
}
