//
//	Boot and authentication module
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var (
	s				*discordgo.Session
	token =			os.Getenv("TOKEN_CARA")
	register =		flag.Bool("r", false, "Register (or re-register) bot commands on startup")
	unregister =	flag.Bool("u", false, "Unregister bot commands on shutdown")
)

func init() {
	// Print boot "splash"
	fmt.Println("╔═══════════════════════╗")
	fmt.Println("║ Cara by cyckl         ║")
	fmt.Println(fmt.Sprintf("║ Running version %s ║", version))
	fmt.Println("╚═══════════════════════╝")
	log.Println("[Info] Minimum permissions are 8")
	
	// Pass args in
	flag.Parse()
	
	configParse()
}

func main() {
	// Declare error here so it can be set without :=
	var err error
	
	// Create bot client session
	log.Println("[Info] Logging in")
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("[Fatal] Error creating session: %v", err)
	}
	
	// Pass on command events to functions
	s.AddHandler(guildInfo)
	s.AddHandler(about)
	s.AddHandler(help)
	s.AddHandler(sticks)
	s.AddHandler(avatar)
	s.AddHandler(tweet)
	s.AddHandler(soundboard)
	s.AddHandler(restore)
	s.AddHandler(kick)
	s.AddHandler(respond)
	s.AddHandler(selfEsteem)
	s.AddHandler(status)
	s.AddHandler(gutEnable)
	s.AddHandler(gutSpam)

	// We only care about integration (command) intents
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildIntegrations | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildVoiceStates)

	// Open websocket connection to Discord and listen
	err = s.Open()
	if err != nil {
		log.Fatalf("[Fatal] Error opening connection: %v", err)
	}
	
	// Default status
	err = s.UpdateStatusComplex(*newUpdateStatusData(0, 3, "you", "https://cara.cyckl.net"))
	if err != nil {
		log.Println("[Error][Boot] Failed to set status:", err)
		return
	}
	
	if *register == true {
		registerCmd(s)
	}

	// Close Discord session cleanly
	defer s.Close()
	
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop
	
	if *unregister == true {
		unregisterCmd(s)
	}
	
	log.Println("[Info] Shutting down")
}
