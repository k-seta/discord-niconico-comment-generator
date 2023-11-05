package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kurehajime/flagstone"
)

// Variables used for command line parameters
var (
	Token string
)

type Comment struct {
	No      uint32 `xml:"no,attr"`
	Time    int64  `xml:"time,attr"`
	Owner   int    `xml:"owner,attr"`
	Service string `xml:"service,attr"`
	Handle  string `xml:"handle,attr"`
	Message string `xml:",innerxml"`
}

type CommentXml struct {
	XMLName xml.Name  `xml:"log"`
	Log     []Comment `xml:"comment"`
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// start Web GUI
	flagstone.Launch("Discord NicoNico Comment Generator", "comment.xml generator for NiCommentGenerator")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!" for healthcheck
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// generate xml
	comment := Comment{
		No:      0,
		Time:    time.Now().Unix(),
		Owner:   0,
		Service: "discord",
		Handle:  m.Author.Username,
		Message: m.Content,
	}
	fmt.Println(comment)

	data, _ := os.ReadFile("comment.xml")
	comments := CommentXml{}
	xml.Unmarshal(data, &comments)

	comments.Log = append(comments.Log, comment)

	var buf bytes.Buffer
	buf.Write([]byte(xml.Header))
	b, _ := xml.MarshalIndent(comments, "", "  ")
	buf.Write(b)
	os.WriteFile("comment.xml", buf.Bytes(), 0666)
}
