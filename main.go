package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/joho/godotenv"
)

var (
	// GiftRegex The regular expression used to determine if the message contains a gift code.
	GiftRegex = regexp.MustCompile("(discord.com/gift/|discord.gifts/|discord.gift/)([a-zA-Z0-9]+)")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	Token := os.Getenv("USER_TOKEN")

	dg, err := discordgo.New(Token)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Logged in as " + dg.State.User.Username)
	fmt.Println("Running nitro sniper on " + strconv.Itoa(len(dg.State.Guilds)) + " servers.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if GiftRegex.Match([]byte(m.Content)) {
		code := GiftRegex.FindStringSubmatch(m.Content)

		fmt.Println(code)
	}
}
