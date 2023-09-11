package config

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Mode int64

const (
	Uknown Mode = iota
	Update
	Rebuild
	Purge
)

type Source int64

const (
	Collection Source = iota
	Meta
)

type MetaConfig struct {
	Path     string
	BySeries int
	ByAuthor int
	ByTags   int
	Tags     map[string]string
}

type Config struct {
	CalibreCollection string
	KindleCC          string
	Mode              Mode
	Source            Source
	Meta              MetaConfig
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		CalibreCollection: "/mnt/us/system/collections.json",
		KindleCC:          "/var/local/cc.db",
		Mode:              Uknown,
		Source:            Meta,
		Meta: MetaConfig{
			Path: "/mnt/us/metadata.calibre",
			Tags: make(map[string]string),
		},
	}

	meta, err := ini.Load("meta.ini")
	if err != nil {
		return nil, err
	}
	cfg.Meta.ByTags = meta.Section("general").Key("tags").MustInt()
	cfg.Meta.ByAuthor = meta.Section("general").Key("author").MustInt()
	cfg.Meta.BySeries = meta.Section("general").Key("series").MustInt()
	for _, kv := range meta.Section("tags").Keys() {
		k := strings.ToLower(kv.Name())
		v := kv.MustString("")
		if len(v) > 0 {
			cfg.Meta.Tags[k] = v
		}
	}

	buf := bytes.Buffer{}
	buf2str := func() string {
		return strings.TrimRight(buf.String(), "\n")
	}
	fs := flag.NewFlagSet("lsync", flag.ContinueOnError)
	fs.SetOutput(&buf)
	fs.Func("mode", "update | rebuild | purge", func(s string) error {
		mode := strings.ToLower(s)
		if mode == "update" {
			cfg.Mode = Update
		} else if mode == "rebuild" {
			cfg.Mode = Rebuild
		} else if mode == "purge" {
			cfg.Mode = Purge
		} else {
			return errors.New("valid choices [update | rebuild | purge]")
		}
		return nil
	})

	fs.Func("source", "collection | meta", func(s string) error {
		mode := strings.ToLower(s)
		if mode == "collection" {
			cfg.Source = Collection
		} else if mode == "meta" {
			cfg.Source = Meta
		} else {
			return errors.New("valid choices [collection | meta]")
		}
		return nil
	})

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf(buf2str())
	}

	if cfg.Mode == Uknown {
		fs.Usage()
		return nil, fmt.Errorf("missing required \"-mode\" flag\n%s", buf2str())
	}

	env_calibre_c := os.Getenv("CCS_C_COLLECTION")
	if len(env_calibre_c) > 0 {
		cfg.CalibreCollection = env_calibre_c
	}

	env_calibre_meta := os.Getenv("CCS_C_META")
	if len(env_calibre_meta) > 0 {
		cfg.Meta.Path = env_calibre_meta
	}

	env_kindle_c := os.Getenv("CCS_K_COLLECTION")
	if len(env_kindle_c) > 0 {
		cfg.KindleCC = env_kindle_c
	}

	os.Setenv("NO_PROXY", "127.0.0.1,localhost")

	return cfg, nil
}
