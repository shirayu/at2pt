package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/shirayu/at2pt"
)

func getFile(ifname string, ofname string) (inf, outf *os.File, err error) {

	if ifname == "-" {
		inf = os.Stdin
	} else {
		inf, err = os.Open(ifname)
		if err != nil {
			return nil, nil, err
		}
	}

	if ofname == "-" {
		outf = os.Stdout
	} else {
		outf, err = os.Create(ofname)
		if err != nil {
			return inf, nil, err
		}
	}

	return inf, outf, nil
}

type cmdOptions struct {
	Input  string `short:"i" long:"input" description:"Input file name. - or no designation means STDIN" default:"-"`
	Output string `short:"o" long:"output" description:"Output file name. - or no designation means STDOUT" default:"-"`
	//     Log      bool   `long:"log" description:"Enable logging" default:"false"`
	Mode  at2pt.Mode `short:"m" long:"mode" description:"Mode {0:PLAIN, 1:TOKENIZED, 2:TOKENIZEDwPRED}" default:"0"`
	Style string     `short:"s" long:"style" description:"Input file style {KNP, MeCab, CaboCha}" default:"KNP"`

	Version bool `short:"v" long:"version" description:"Show version"`
}

//Version of this program
var Version = "Unknown version"

//VersionDate is the commit date of this version
var VersionDate = ""

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = "at2pt"
	optparser.Usage = "-i input -o output [OPTIONS]"
	_, err := optparser.Parse()

	//show help
	if len(os.Args) == 1 {
		optparser.WriteHelp(os.Stdout)
		os.Exit(0)
	} else if err != nil {
		if e, ok := err.(*flags.Error); !ok {
			log.Fatalf("Expected flags.Error, but got %T", err)
		} else if e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("%s\n", err)
	} else if opts.Version {
		if len(VersionDate) != 0 {
			fmt.Fprintf(os.Stderr, "%s: %s on %s\n", optparser.Name, Version, VersionDate)
		} else {
			fmt.Fprintf(os.Stderr, "%s: %s\n", optparser.Name, Version)
		}
		os.Exit(0)
	}

	inf, outf, err := getFile(opts.Input, opts.Output)
	if err == nil {
		defer inf.Close()
		defer outf.Close()
		if strings.ToLower(opts.Style) == "knp" {
			err = at2pt.ConvertKNP(inf, outf, opts.Mode)
		} else if strings.ToLower(opts.Style) == "cabocha" {
			err = at2pt.ConvertCaboCha(inf, outf, opts.Mode)
		} else if strings.ToLower(opts.Style) == "mecab" {
			err = at2pt.ConvertCaboCha(inf, outf, opts.Mode)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown style: %s\n", opts.Style)
			os.Exit(1)
		}
	} else {
		if inf == nil {
			fmt.Fprintf(os.Stderr, "Error Input: %s in %s\n", err, opts.Input)
		}
		if outf == nil {
			fmt.Fprintf(os.Stderr, "Error Output: %s in %s\n", err, opts.Output)
		}
		os.Exit(1)
	}
}
