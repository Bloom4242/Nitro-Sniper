package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/bwmarrin/discordgo"

	"github.com/joho/godotenv"
)

var (
	// Token The user's authorization token used to make requests to discord's api.
	Token string

	// Sniped The amount of codes that have been sniped.
	Sniped int

	// GiftRegex The regular expression used to determine if the message contains a gift code.
	GiftRegex = regexp.MustCompile("(discord.com/gifts/|discord.gift/)([a-zA-Z0-9]+)")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	Token = os.Getenv("USER_TOKEN")

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
		code := GiftRegex.FindStringSubmatch(m.Content)[2]

		if len(code) < 16 {
			return
		}

		var strPost = []byte("POST")
		var strRequestURI = []byte("https://discordapp.com/api/v8/entitlements/gift-codes/" + code + "/redeem")
		var strRequestBody = []byte(`{"channel_id":` + m.ChannelID + `, "payment_source_id": null}`)

		req := fasthttp.AcquireRequest()
		req.Header.SetMethodBytes(strPost)
		req.Header.SetContentType("application/json")
		req.Header.Set("authorization", Token)
		req.SetBody(strRequestBody)
		req.Header.SetRequestURIBytes(strRequestURI)
		res := fasthttp.AcquireResponse()

		if err := fasthttp.Do(req, res); err != nil {
			log.Fatal(err)
		}

		fasthttp.ReleaseRequest(req)

		body := res.Body()

		bodyString := string(body)
		fasthttp.ReleaseResponse(res)

		if strings.Contains(bodyString, "10038") {
			return
		} else if strings.Contains(bodyString, "100011") || strings.Contains(bodyString, "50050") {
			fmt.Println("\n- recieved already redeemed code sent by " + m.Author.Username)
		} else {
			fmt.Println("\n+ Recieved valid code sent by " + m.Author.Username + " - " + code)

			Sniped++
			if Sniped == 3 {
				fmt.Println("Ended sniping at " + time.Now().Format("01-02-2006 15:04:05"))

				s.Close()
			}
		}
	}
}
