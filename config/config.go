package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// todo: func to find config file

type Config struct {
	Programs struct {
		Searcher string `yaml:"searcher"`
		Editor string `yaml:"editor"`
		TxtConverter string `yaml:"txtconverter"`
		PrintCmd string `yaml:"printcmd"`
	} `yaml:"programs"`
	Args struct {
		SearcherArgs []string `yaml:"searcherargs"`
		LprArgs string `yaml:"lprargs"`
		TxtConvArgs string `yaml:"txtconvargs"`
		Printer string `yaml:"printer"`
		PrintDuplex string `yaml:"printduplex"`
		ColorPrint string `yaml:"colorprint"`
	} `yaml:"args"`
}


// FileClose closes a file and exits if error occurs
func fileClose(f *os.File)  {
	err := f.Close()
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

func Conf() Config {
	var cfg Config
	f, err := os.Open("config/config.yml")
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
	return cfg
}
