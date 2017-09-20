package rss

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/kapmahc/fly/plugins/nut/app"
	"golang.org/x/tools/blog/atom"
)

// Write write rss.atom to tmp
func Write(lang, title, user, email string) error {
	feed := atom.Feed{
		Title:   title,
		ID:      uuid.New().String(),
		Updated: atom.Time(time.Now()),
		Author: &atom.Person{
			Name:  user,
			Email: email,
		},
		Entry: make([]*atom.Entry, 0),
	}

	home := app.Home()

	for _, hnd := range handlers {
		items, err := hnd(lang)
		if err != nil {
			return err
		}
		for _, it := range items {
			for i := range it.Link {
				it.Link[i].Href = home + it.Link[i].Href
			}
			feed.Entry = append(feed.Entry, it)
		}
	}
	fn := path.Join("tmp", fmt.Sprintf("rss-%s.atom", lang))
	log.Printf("generate file %s", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	enc := xml.NewEncoder(fd)
	return enc.Encode(feed)
}
