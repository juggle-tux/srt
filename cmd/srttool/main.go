// srttool is a cli tool to adjust the times in a srt subtitle file
package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jawher/mow.cli"
	"github.com/juggle-tux/srt"
)

var app = cli.App(os.Args[0], "time adjustment tool for srt subtitle files")
var (
	clInfiles = app.Strings(cli.StringsArg{
		Name: "SRTFILES", Desc: "Input files", HideValue: true,
	})
	clOutfile = app.String(cli.StringOpt{
		Name: "f output", Desc: "Output file (default is stdout)",
	})
	clOffset = app.String(cli.StringOpt{
		Name: "o offset", Value: "0s", Desc: "Time offset",
	})
)

func main() {
	app.Spec = "[OPTIONS] SRTFILES..."
	app.Action = run
	app.Run(os.Args)
}

func run() {
	offset, err := time.ParseDuration(*clOffset)
	if err != nil {
		log.Fatal(err)
	}

	var outfile *os.File
	if *clOutfile == "" {
		outfile = os.Stdout
	} else {
		outfile, err = os.Create(*clOutfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer outfile.Close()

	enc := srt.NewEncoder(outfile, 0)
	defer enc.Flush()

	// Main loop
	for _, f := range *clInfiles {
		if err := func() error {
			f, err := os.Open(f)
			if err != nil {
				return err
			}
			defer f.Close()
			dec := srt.NewDecoder(f)

			for b, err := dec.Next(); err != io.EOF; b, err = dec.Next() {
				if err != nil {
					return err
				}
				b.Add(offset)
				if err := enc.Block(b); err != nil {
					return err
				}
			}
			return nil
		}(); err != nil {
			log.Fatal(err)
		}
	}
}
