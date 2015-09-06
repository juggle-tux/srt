// srttool is a cli tool to adjust the times in a srt subtitle file
package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jawher/mow.cli"
	"juggle.tux/srt"
)

var app = cli.App("srttool-go", "time adjustment tool for srt subtitle files")
var (
	ClInfiles = app.Strings(cli.StringsArg{Name: "SRTFILES", Desc: "Input files", HideValue: true})
	ClOutfile = app.String(cli.StringOpt{Name: "f output", Desc: "Output file"})
	ClOffset  = app.String(cli.StringOpt{Name: "o offset", Value: "0", Desc: "Time offset"})
)

func init() {
	app.Spec = "[OPTIONS] SRTFILES..."
	app.Action = run
}

func run() {
	offset, err := time.ParseDuration(*ClOffset)
	if err != nil {
		log.Fatal(err)
	}

	var outfile *os.File
	if *ClOutfile == "" {
		outfile = os.Stdout
	} else {
		outfile, err = os.Create(*ClOutfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer outfile.Close()
	enc := srt.NewEncoder(outfile, 0)
	defer enc.Flush()

	// Main loop
	for _, f := range *ClInfiles {
		bs, err := getFile(f)
		if err != nil {
			log.Print(err)
			return
		}
		for _, b := range bs {
			b.Add(offset)
			if err := enc.Block(b); err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func main() {
	app.Run(os.Args)
}

func getFile(file string) ([]srt.Block, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := srt.NewDecoder(f)
	bs := make([]srt.Block, 0, 100)
	for b, err := dec.Next(); err != io.EOF; b, err = dec.Next() {
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}
