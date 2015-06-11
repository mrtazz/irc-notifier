package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/deckarep/gosx-notifier"
	"github.com/gosexy/redis"
	"github.com/rakyll/globalconf"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

type Notification struct {
	Source string   `json:"@source"`
	Tags   []string `json:"@tags"`
	Fields struct {
		Time    []string `json:"irctime"`
		Sender  []string `json:"ircsender"`
		Message []string `json:"ircmessage"`
	} `json:"@fields"`
	TimeStamp  string `json:"@timestamp"`
	SourceHost string `json:"@source_host"`
	SourcePath string `json:"@source_path"`
	Message    string `json:"@message"`
	Type       string `json:"@type"`
}

type Message struct {
	Message  string
	Title    string
	Subtitle string
}

func main() {
	var (
		host      = flag.String("host", "", "the redis host")
		port      = flag.Int("port", 6379, "the redis port")
		auth      = flag.String("auth", "", "the auth key if set")
		appIcon   = flag.String("icon", "", "path to the app icon")
		show_help = flag.Bool("help", false, "show help")
	)

	conf, _ := globalconf.New("irc-notifier")
	conf.ParseAll()

	if *show_help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: irc-notifier\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	var client *redis.Client

	if *host == "" {
		fmt.Println("You need to specify a host with --host")
		return
	}

	client = redis.New()

	err := client.Connect(*host, uint(*port))

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return
	}

	if *auth != "" {
		_, err = client.Auth(*auth)
		if err != nil {
			log.Fatalf("Could not auth: %s\n", err.Error())
			return
		}
	}

	for {
		raw_notification, err := client.BLPop(1000, "notifications")
		if err != nil {
			log.Fatalf("Could not get notifications: %s\n", err.Error())
		} else if len(raw_notification) > 1 {
			var notification Notification
			err = json.Unmarshal([]byte(raw_notification[1]), &notification)
			if err != nil {
				log.Printf("%s", err.Error())
			}
			msg, _ := ParseLogLine(notification)
			go Notify(msg.Message, msg.Title, msg.Subtitle, *appIcon)
		}
	}

	client.Quit()
}

func ParseLogLine(notification Notification) (Message, error) {
	var ret Message
	defer func() {
		if r := recover(); r != nil {
			_, file, line, _ := runtime.Caller(2)
			log.Println(file)
			log.Println(line)
			log.Printf("Runtime panic: %v", r)
			log.Printf("Stack: %s", debug.Stack())
		}
	}()
	split_source := strings.Split(notification.Source, "/")
	split_source = strings.Split(split_source[len(split_source)-1], "_")
	var channel string
	if len(split_source) > 2 {
		channel = strings.Join(split_source[1:len(split_source)-1], "-")
	} else {
		channel = ""
	}
	if len(notification.Fields.Sender) > 0 {
		ret.Subtitle = fmt.Sprintf("%s in %s", notification.Fields.Sender[0], channel)
	}
	ret.Title = "Etsy IRC"
	if len(notification.Fields.Message) > 0 {
		ret.Message = notification.Fields.Message[0]
	} else {
		log.Printf("Empty message found for %v.", notification.Fields.Message)
		ret.Message = ""
	}
	return ret, nil
}

func Notify(message string, title string, subtitle string, icon string) error {
	//At a minimum specifiy a message to display to end-user.
	note := gosxnotifier.NewNotification(message)
	note.Group = "com.github.mrtazz.ircnotifier"

	//Optionally, set a title
	note.Title = title

	//Optionally, set a subtitle
	note.Subtitle = subtitle

	// set icon if we didn't get an empty string for it
	if icon != "" {
		note.AppIcon = icon
		note.ContentImage = icon
	} else {
		note.Sender = "net.limechat.LimeChat"
	}

	//Then, push the notification
	err := note.Push()

	//If necessary, check error
	if err != nil {
		log.Println("Uh oh!")
	}
	return nil
}
