package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io/ioutil"
	"math/rand"
	"github.com/joho/godotenv"
	"github.com/jrm780/gotirc"
)


type Player struct {
	Name string
	cards []string
}

func generateCard() string {
	cards := [53]string{"red 1","red 2","red 3","red 4","red 5","red 6","red 7","red 8","red 9","blue 1","blue 2","blue 3","blue 4","blue 5","blue 6","blue 7","blue 8","blue 9", "green 1","green 2","green 3","green 4","green 5","green 6","green 7","green 8","green 9", "yellow 1","yellow 2","yellow 3","yellow 4","yellow 5","yellow 6","yellow 7","yellow 8","yellow 9", "wild 4", "red reverse", "yellow reverse", "blue reverse", "green reverse", "red +2", "green +2", "yellow +2", "blue +2"}

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

		d1 := []byte("true")
		err := ioutil.WriteFile("unoStarted.txt", d1, 0644)
		check(err)

	}
}

func enterUNO(messageString []string, tags map[string]string, players []Player, unoStarted string) Player{
	if messageString[0] == "!enter" && unoStarted == "true"{
		var player Player
		player.Name = tags["display-name"]

		return player
	}
	var player Player
	player.Name = ""

	return player
}

func endUno(messageString []string, tags map[string]string, channel string, players []Player, unoStarted string){
	channelClean := strings.ReplaceAll(channel, "#", "")
	fmt.Println( messageString[0] == "!enduno" && tags["display-name"] == channelClean)
	if messageString[0] == "!enduno" && tags["display-name"] == channelClean{

		d1 := []byte("false")
		err := ioutil.WriteFile("unoStarted.txt", d1, 0644)
		check(err)

		for i := 0;  i < len(players); i++{
			player := players[i]

			if player.Name != ""{
				for i := 0; i < 7; i++{
					pick := generateCard()
					player.cards = append(player.cards, pick)
				}
			}
		}

		fmt.Println("players (endUno): ", players)

	}
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
	client.OnChat(func(channel string, tags map[string]string, msg string) {

		fmt.Println(msg)

		messageString := strings.Split(msg, " ")

		unoStarted := checkUnoStarted()

		initUNO(messageString, tags, channel)
		players = append(players, enterUNO(messageString, tags, players, unoStarted))
		fmt.Println("Players (main): ", players)
		endUno(messageString, tags, channel, players, unoStarted)

	})

	client.Connect(user, token)
}
