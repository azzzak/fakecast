package feed

import (
	"encoding/xml"
	"strings"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
)

// RSS entity
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Itunes  string   `xml:"xmlns:itunes,attr"`
	Content string   `xml:"xmlns:content,attr"`

	Channel Channel `xml:"channel"`
}

// Channel entity
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link,omitempty"`
	Copyright   string `xml:"copyright,omitempty"`
	Author      string `xml:"itunes:author,omitempty"`
	Description string `xml:"description,omitempty"`
	Type        string `xml:"itunes:type,omitempty"`
	Image       struct {
		Href string `xml:"href,attr,omitempty"`
	} `xml:"itunes:image"`

	Items []Item `xml:"item"`
}

// Item entity
type Item struct {
	Title     string    `xml:"title"`
	Enclosure Enclosure `xml:"enclosure"`

	GUID        string `xml:"guid,omitempty"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description,omitempty"`
	Duration    int    `xml:"itunes:duration"`
	Link        string `xml:"link,omitempty"`
	Explicit    bool   `xml:"itunes:explicit,omitempty"`
	Season      int    `xml:"itunes:season,omitempty"`
	Episode     int    `xml:"itunes:episode,omitempty"`
}

// Enclosure entity
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
}

// GenerateFeed with content of channel
func GenerateFeed(channel *store.Channel, podcasts []store.Podcast, host string) RSS {
	feed := RSS{
		Version: "2.0",
		Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Content: "http://purl.org/rss/1.0/modules/content/",
	}

	feed.Channel = Channel{
		Title:       channel.Title,
		Description: channel.Description,
		Author:      channel.Author,
	}

	feed.Channel.Image.Href = channel.Cover

	var items []Item
	for _, p := range podcasts {
		var explicit bool
		if p.Explicit == 1 {
			explicit = true
		}

		// Types
		// mp3 (*.mp3): audio/mpeg
		// aac (*.m4a): audio/x-m4a

		cType := "audio/mpeg"
		if _, ext := fs.NameAndExtFrom(p.Filename); ext == "m4a" {
			cType = "audio/x-m4a"
		}

		items = append(items, Item{
			Title: p.Title,
			Enclosure: Enclosure{
				URL:    strings.Join([]string{host, "files", channel.Alias, p.Filename}, "/"),
				Length: p.Length,
				Type:   cType,
			},
			GUID:        p.GUID,
			PubDate:     p.PubDate,
			Description: p.Description,
			Duration:    p.Duration,
			Explicit:    explicit,
			Season:      p.Season,
			Episode:     p.Episode,
		})
	}

	feed.Channel.Items = items
	return feed
}
