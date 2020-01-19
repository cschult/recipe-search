package h

import (
	"fmt"
)

func Help()  {
	fmt.Println()
	fmt.Println("" +
		"Help\n" +
		"\n" +
		"h: help - this help\n" +
		"q: quit - quit program\n" +
		"n: new search - enter new search term\n" +
		"l: long - list files with URN\n" +
		"s: short - list only filenames\n" +
		"p: print - print file to default printer\n" +
					"\t\tprinter must be set in config file\n" +
					"\t\tor in environment variable PRINTER\n" +
		"e: edit - edit recipe with configured editor\n" +
		"1|2|3|...: enter number of file to view file\n" +
		"\n" +
		"Search Patterns\n" +
		"\n" +
		"word               = all files having word\n" +
		"word*              = all files having word beginning with word\n" +
		"word1 word2        = all files having both words\n" +
		"word1 AND word2    = all files having both words\n" +
		"word1 OR word2     = all files having one of both words\n" +
		"word1 -word2       = all files having word1 but not word2\n" +
		"\n" +
		"you can combine above search patterns\n" +
		"\n" +
		"ENTER              = all files")
	fmt.Println()
}
