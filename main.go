package main

import (
	"devmem.de/srv/git/recipe-search/internal/h"
	"devmem.de/srv/git/recipe-search/internal/rs"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

// todo: func to find config file
// todo: message when search has no result: sorry and show new search dialog prompt, do not print helpline
// todo: ensure that only txt files are concatenated
// todo: what happens if a write protected recipe is edited? how to catch editor errors

type Config struct {
	Programs struct {
		Searcher     string `yaml:"searcher"`
		Editor       string `yaml:"editor"`
		TxtConverter string `yaml:"txtconverter"`
		PrintCmd     string `yaml:"printcmd"`
	} `yaml:"programs"`
	Args struct {
		SearcherArgs []string `yaml:"searcherargs"`
		LprArgs      string   `yaml:"lprargs"`
		TxtConvArgs  string   `yaml:"txtconvargs"`
		Printer      string   `yaml:"printer"`
		PrintDuplex  string   `yaml:"printduplex"`
		ColorPrint   string   `yaml:"colorprint"`
	} `yaml:"args"`
	Flags struct {
		Uri bool
	} `yaml:"flags"`
}

// FileClose closes a file and exits if error occurs
func fileClose(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// main is the text user interface of this program.
// It starts a loop for interacting with the user.
func main() {
	// C O N F I G U R A T I O N
	var cfg Config
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}

	defer fileClose(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		fmt.Println(err)
	}

	// put all configs for printing into a map
	prntcfg := map[string]string{
		"prntcmd":       cfg.Programs.PrintCmd,
		"converter":     cfg.Programs.TxtConverter,
		"prntcmdArgs":   cfg.Args.LprArgs,
		"converterArgs": cfg.Args.TxtConvArgs,
		"printer":       cfg.Args.Printer,
		"prntduplex":    cfg.Args.PrintDuplex,
		"prntcolor":     cfg.Args.ColorPrint,
	}
	// flag indicating if files are listed as URI or filename only
	// false means: as filename (default)
	Uri := cfg.Flags.Uri

	resultPathFile, resultFile := rs.Search(cfg.Programs.Searcher, cfg.Args.SearcherArgs)
	rs.ViewResult(resultPathFile, resultFile, Uri)
	helpLine := color.CyanString("h help | q quit | n new search | l long |" +
		" s short | p print | e edit | 1 = show file #1 | 2 = ...")

	for true {
		fmt.Println(helpLine)
		fmt.Printf("Enter a key: ")
		key := rs.Input()

		switch key {
		case "h":
			h.Help()
		case "q": // quit
			os.Exit(0)
		case "n": // new search
			// clear the list of command line args so that rs.args()
			// asks user for new search term
			myName := os.Args[0]
			os.Args = []string{myName}
			resultPathFile, resultFile = rs.Search(cfg.Programs.Searcher, cfg.Args.SearcherArgs)
			rs.ViewResult(resultPathFile, resultFile, Uri)
		case "l": // path and filename
			Uri = true
			rs.ViewResult(resultPathFile, resultFile, Uri)
		case "s": // only filename
			Uri = false
			rs.ViewResult(resultPathFile, resultFile, Uri)
		case "e":
			err = rs.EditFile(cfg.Programs.Editor, resultPathFile)
			rs.ViewResult(resultPathFile, resultFile, Uri)
			if err != nil {
				rs.PrtErr("Oops!", err)
			}
		case "p": // send file to printer
			err = rs.Print(prntcfg, resultPathFile)
			rs.ViewResult(resultPathFile, resultFile, Uri)
			if err != nil {
				rs.PrtErr("Oops!", err)
			}
		case "": // ENTER, print help line
			rs.ViewResult(resultPathFile, resultFile, Uri)
			color.Yellow("enter a valid key or a file number\n\n")
		default:
			i, err := strconv.Atoi(key)
			if err == nil {
				if i >= 1 && i <= len(resultPathFile) {
					err = rs.ConcatFile(resultPathFile, resultFile, i-1)
					rs.ViewResult(resultPathFile, resultFile, Uri)
					if err != nil {
						// fmt.Fprintf(os.Stderr, "%s\t", err)
						rs.PrtErr("Oops!", err)
					}
				} else {
					rs.ViewResult(resultPathFile, resultFile, Uri)
					color.Yellow("not a valid file number\n\n")
				}
			} else {
				rs.ViewResult(resultPathFile, resultFile, Uri)
				color.Yellow("invalid key\n\n")
			}
		}
	}
}
