package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Mode int64

const (
	Uknown Mode = iota
	Update
	Rebuild
	Purge
)

type Config struct {
	CalibreCollection string
	CalibreMeta       string
	KindleCC          string
	Mode              Mode
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		CalibreCollection: "/mnt/us/system/collections.json",
		CalibreMeta:       "/mnt/us/metadata.calibre",
		KindleCC:          "/var/local/cc.db",
		Mode:              Uknown,
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
	err := fs.Parse(os.Args[1:])
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

	env_kindle_c := os.Getenv("CCS_K_COLLECTION")
	if len(env_kindle_c) > 0 {
		cfg.KindleCC = env_kindle_c
	}

	os.Setenv("NO_PROXY", "127.0.0.1,localhost")

	return cfg, nil
}
