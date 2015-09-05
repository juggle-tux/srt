package main // srttool is a cli tool to adjust the times in a srt subtitle file

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"juggle.tux/srt"
)

var ( // FLags
	offset  = flag.Duration("offset", 0, "")
	outfile = flag.String("o", "", "output file default stdout")
	infiles []string

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	infiles = flag.Args()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	var ofile *os.File
	var err error
	if *outfile == "" {
		ofile = os.Stdout
	} else {
		ofile, err = os.Create(*outfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer ofile.Close()
	enc := srt.NewEncoder(ofile, 0)
	defer enc.Flush()

	// Main loop
	for _, f := range infiles {
		bs, err := getFile(f)
		if err != nil {
			log.Print(err)
			return
		}
		for _, b := range bs {
			b.Add(*offset)
			if err := enc.Block(b); err != nil {
				log.Print(err)
				return
			}
		}
	}
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
