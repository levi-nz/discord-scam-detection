package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// auditMessage audits a message to check for malicious content.
func auditMessage(s *discordgo.Session, event *discordgo.MessageCreate, addresses []*url.URL) {
	// markMalicious logs the message as malicious and deletes it.
	markMalicious := func(reason string) {
		log.Println("Malicious image detected:", reason)

		if err := s.ChannelMessageDelete(event.ChannelID, event.Message.ID); err != nil {
			log.Println("Error deleting message:", err)
		}

		log.Println("Deleted message:", event.Message.Content)
	}

	client := http.Client{}

	for _, addr := range addresses {
		response, err := client.Get(addr.String())
		if err != nil {
			// log for debug purposes and move on
			log.Printf("Failed to GET %s: %v\n", addr.String(), err)
			continue
		}

		// mark as malicious if non-OK status is received
		if response.StatusCode != http.StatusOK {
			markMalicious("Bad HTTP status code received from server: " + response.Status)
			return
		}

		// parse response body as HTML
		doc, err := html.Parse(response.Body)
		if err != nil {
			// log and move on
			log.Printf("Failed to read response body from %s: %v\n", addr.String(), err)
			continue
		}

		// discovered image URL
		var imageURL string
		// HTML parsing hell
		var f func(node *html.Node)
		f = func(node *html.Node) {
			if node.Type == html.ElementNode && node.Data == "meta" {
				// check if we have "og:image" tag
				var isImage bool
				for _, attribute := range node.Attr {
					if attribute.Val == "og:image" {
						isImage = true
						break
					}
				}

				// is this an image?
				if isImage {
					// check for "content" attribute
					for _, attribute := range node.Attr {
						if attribute.Key == "content" {
							imageURL = attribute.Val
							break
						}
					}
				}
			}

			for c := node.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		_ = response.Body.Close()

		// get image content if we have a URL
		if imageURL != "" {
			// if the image URL doesn't start with http:// or https://, then this is likely a resource location
			// ie: /resources/image.png, resources/image.png etc
			//
			// if the content is the first scenario (/resources/image.png), then this
			// is an absolute resource location, ex: example.com/resources/image.png
			//
			// otherwise, the location is relative (resources/image.png).
			// this means the location depends on the URL, but is basically
			// the current URL + /location. ex:
			// example.com/wherever-we-are-currently/resources/image.png
			if !(strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://")) {
				if strings.HasPrefix(imageURL, "/") {
					// location is absolute
					imageURL = fmt.Sprintf("%s://%s%s", addr.Scheme, addr.Host, imageURL)
				} else {
					// location is relative
					imageURL = addr.String() + "/" + imageURL
				}
			}

			// fetch image
			if response, err = client.Get(imageURL); err != nil {
				// move on
				continue
			}

			// assume malicious if we get non-OK
			if response.StatusCode != http.StatusOK {
				markMalicious(fmt.Sprintf("Got non-OK status %s from %s", response.Status, imageURL))
				return
			}

			// read body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				continue
			}
			_ = response.Body.Close()

			// check if malicious
			if encodedChecksum, isMalicious := IsImageMalicious(body); isMalicious {
				markMalicious(fmt.Sprintf("Malicious image hash %s from website %s (from URL: %s)", encodedChecksum, addr.String(), imageURL))
				return
			}
		}
	}
}
