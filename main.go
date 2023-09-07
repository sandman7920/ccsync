package main

import (
	"fmt"
	"os"

	"github.com/sandman7920/ccsync/calibre"
	"github.com/sandman7920/ccsync/kindle"
)

func run() error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}
	updater := &kindle.Commander{}

	ke, err := kindle.NewEntries(cfg.KindleCC)
	if err != nil {
		return err
	}

	if cfg.Mode == Rebuild || cfg.Mode == Purge {
		err := updater.Purge(ke.Collection)
		if err != nil {
			return err
		}

		if cfg.Mode == Purge {
			return nil
		}

		ke.Collection = nil
	}

	if err != nil {
		return err
	}

	cc, err := calibre.NewCollection(cfg.CalibreCollection)
	if err != nil {
		return err
	}

	unique := map[string]int{}
	update := kindle.Collection{}

	for _, c := range cc {
		var uuid string
		cde_books := ke.Books.BooksByCDEKeys(c.CDEKeys())
		if idx := ke.Collection.IdxByTitle(c.Title); idx != -1 {
			books := ke.Collection[idx].Books

			if len(c.Entities) == len(books) && len(c.Entities) == len(cde_books) {
				for _, b := range cde_books {
					unique[b.UUID] += 1
				}
				continue
			}
			uuid = ke.Collection[idx].UUID
		}
		entry := &kindle.CollEntry{
			UUID:  uuid,
			Title: c.Title,
			Books: cde_books,
		}

		for _, b := range entry.Books {
			unique[b.UUID] += 1
		}
		update = append(update, entry)
	}

	if len(update) > 0 {
		return updater.Update(update, ke.IsCcAware, unique)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}