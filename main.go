package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	ChannelID string
	Filepath  string
	Message   *widget.Label
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

func main() {

	// Fyne
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})

	myWindow := myApp.NewWindow("Discord Niconico Comment Generator")
	myWindow.Resize(fyne.NewSize(300, 400))

	channelID := widget.NewEntry()
	filepath := widget.NewEntry()
	token := widget.NewEntry()

	status := widget.NewLabel("DisConnected.")
	Message = widget.NewLabel("")

	var dg *discordgo.Session

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Channel ID", Widget: channelID},
			{Text: "Filepath", Widget: filepath},
			{Text: "Discord Token", Widget: token},
		},
		OnSubmit: func() { // optional, handle form submission
			ChannelID = channelID.Text
			Filepath = filepath.Text
			dg = connect(token.Text)
			status.SetText("Connected.")
		},
	}

	content := container.NewVBox(
		form,
		status,
		Message,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

	// Cleanly close down the Discord session.
	dg.Close()
}

func connect(token string) *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil
	}

	return dg
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || m.ChannelID != ChannelID {
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
	commentBuf, _ := xml.Marshal(comment)
	Message.SetText(string(commentBuf))

	data, _ := os.ReadFile(Filepath)
	comments := CommentXml{}
	xml.Unmarshal(data, &comments)

	comments.Log = append(comments.Log, comment)

	var buf bytes.Buffer
	buf.Write([]byte(xml.Header))
	b, _ := xml.MarshalIndent(comments, "", "  ")
	buf.Write(b)
	os.WriteFile(Filepath, buf.Bytes(), 0666)
}
