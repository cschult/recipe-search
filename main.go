package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"devmem.de/srv/git/recipe-search/internal/h"
	"devmem.de/srv/git/recipe-search/internal/rs"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

// todo: debug error handling of Print() and EditFile()
// todo: set editor via config (done) OR via environment variable
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

// check if env var is set and not empty
func isEnv(s string) (string, bool) {
	env, ok := os.LookupEnv(s)
	if len(env) == 0 {
		return env, false
	}
	return env, ok
}

// fileExist tests if a file exists
func fileExists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return false
	}
	return true
}

// findConfig looks for a config file
func findConfig() (string, error) {
	// var path string
	// var ok bool
	var f string
	path, ok := isEnv("XDG_CONFIG_HOME")
	if ok {
		f = fmt.Sprintf("%s/recipe-search/config.yml", path)
	} else {
		path, ok = isEnv("HOME")
		if ok {
			f = fmt.Sprintf("%s/.config/recipe-search/config.yml", path)
		} else {
			// 	$HOME not set
			return "", errors.New("$HOME not set")
		}
	}
	if fileExists(f) {
		return f, nil
	}
	return "", errors.New("no config file")
}

// printer set in config: error when printing
// if not in cfg: setPrinter (read env PRINTER)
//   if PRINTER not set: disable printing
//   if PRINTER is set, test if ok, else disable printing
func setPrinter() (string, error) {
	printer, ok := isEnv("PRINTER")
	if !ok {
		// $PRINTER is not set or empty
		err := errors.New("environment variable PRINTER not set or empty")
		return "error", err
	} else {
		// check if $PRINTER is a valid printer
		cmd := exec.Command("lpq", "-P", printer)
		err := cmd.Start()
		if err != nil {
			// debug: rs.PrtErr("error starting lpq", err)
			return "error starting lpq", err
		}
		err = cmd.Wait()
		if err != nil {
			// debug: rs.PrtErr("error waiting lpq", err)
			return "error waiting lpq", err
		}
	}
	return printer, nil
}

// setEditor looks for env var 'EDITOR'
func setEditor() (string, error) {
	editor, ok := isEnv("EDITOR")
	if !ok {
		// EDITOR not set
		err := errors.New("environment variable EDITOR not set or empty")
		return "error", err
	}
	return editor, nil
}

// main is the text user interface of this program.
// It starts a loop for interacting with the user.
func main() {
	// C O N F I G U R A T I O N
	cf, err := findConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	f, err := os.Open(cf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer fileClose(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		fmt.Println(err)
	}

	// look for printer with $PRINTER, then with lpstat -d if not in config.yml
	if cfg.Args.Printer == "" {
		s, err := setPrinter()
		if err != nil {
			// color.Yellow("no printer, printing disabled")
			cfg.Args.Printer = ""
		} else {
			cfg.Args.Printer = s
		}
	}

	if cfg.Programs.Editor == "" {
		editor, err := setEditor()
		if err != nil {
			cfg.Programs.Editor = ""
		} else {
			cfg.Programs.Editor = editor
		}
	}

	// put all print configs into a map
	prntcfg := map[string]string{
		"prntcmd":       cfg.Programs.PrintCmd,
		"converter":     cfg.Programs.TxtConverter,
		"prntcmdArgs":   cfg.Args.LprArgs,
		"converterArgs": cfg.Args.TxtConvArgs,
		"printer":       cfg.Args.Printer,
		"prntduplex":    cfg.Args.PrintDuplex,
		"prntcolor":     cfg.Args.ColorPrint,
	}
	// flag indicating that search results are listed as URI or filename only
	// false means: as filename (default)
	Uri := cfg.Flags.Uri


	// user interface starting here
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
			fmt.Println("\ngood bye")
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
			if cfg.Args.Printer != "" {
				err = rs.Print(prntcfg, resultPathFile)
				rs.ViewResult(resultPathFile, resultFile, Uri)
				if err != nil {
					rs.PrtErr("Oops!", err)
				}
			} else {
				// no printer configured and no system default printer
				color.Yellow("printing disabled, found no printer")
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
