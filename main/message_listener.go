package main

import (
	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
	"net/url"
)

var (
	rxStrict = xurls.Strict()
)

// receiveMessage handles received messages from Discord's API.
func receiveMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	// don't check bot messages
	if event.Author.Bot {
		return
	}

	// find all URL's in the message
	detectedUrls := rxStrict.FindAllString(event.Content, -1)
	if len(detectedUrls) == 0 {
		// no URL's in message
		return
	}

	// construct a list of all valid URL's
	var addresses []*url.URL
	for _, rawUrl := range detectedUrls {
		addr, err := url.Parse(rawUrl)
		if err != nil {
			// not a URL
			continue
		}

		// ignore if domain in URL is discord.com
		if addr.Hostname() == "discord.com" {
			continue
		}

		// is a valid URL
		addresses = append(addresses, addr)
	}

	// audit message concurrently if we have at least 1 parsed URL
	if len(addresses) > 0 {
		go auditMessage(s, event, addresses)
	}
}
