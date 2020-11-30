package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io/ioutil"
	"github.com/joho/godotenv"
	"github.com/jrm780/gotirc"
)

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

func enterUNO(messageString []string, tags map[string]string, players []string, unoStarted string) string{
	if messageString[0] == "!enter" && unoStarted == "true"{
		return tags["display-name"]
	}
	return ""
}

func endUno(messageString []string, tags map[string]string, channel string, players []string, unoStarted string){
	channelClean := strings.ReplaceAll(channel, "#", "")
	fmt.Println( messageString[0] == "!enduno" && tags["display-name"] == channelClean)
	if messageString[0] == "!enduno" && tags["display-name"] == channelClean{

    d1 := []byte("false")
    err := ioutil.WriteFile("unoStarted.txt", d1, 0644)
    check(err)

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

var players []string

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
