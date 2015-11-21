package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kr/s3"
)

var (
	distURL      = os.Getenv("DISTURL")
	s3DistURL    = os.Getenv("S3DISTURL")
	s3PatchURL   = os.Getenv("S3PATCHURL")
	buildName    = os.Getenv("BUILDNAME")
	netrcPath    = filepath.Join(os.Getenv("HOME"), ".netrc")
	buildbranch  = os.Getenv("BUILDBRANCH")
	hkgenAppName = os.Getenv("HKGENAPPNAME")
	s3keys       = s3.Keys{
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
	}
)

type release struct {
	Plat   string `json:"platform"`
	Ver    string `json:"version"`
	Cmd    string `json:"cmd"`
	Sha256 []byte `json:"sha256"`
}

func (r release) Name() string {
	return r.Cmd + "/" + r.Ver + "/" + r.Plat
}

func (r release) Gzname() string {
	return r.Name() + ".gz"
}

var subcommands = map[string]func([]string){
	"gen":   gen,
	"build": build,
	"web":   web,
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: mcdist (web|gen|build [platforms])")
	os.Exit(2)
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		usage()
	} else if os.Args[1] == "web" && len(os.Args) != 2 {
		usage()
	} else if os.Args[1] == "gen" && len(os.Args) != 6 {
		usage()
	}
	f := subcommands[os.Args[1]]
	if f == nil {
		usage()
	}
	f(os.Args[2:])
}
