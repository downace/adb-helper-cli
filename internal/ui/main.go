package ui

import (
	"bufio"
	"fmt"
	"github.com/ttacon/chalk"
	"os"
	"strconv"
	"strings"
)

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(label)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func SelectPrompt[T any](label string, options []T) uint {
	for {
		fmt.Println(label)
		for n, option := range options {
			fmt.Printf(" %s. %s\n", chalk.Cyan.Color(fmt.Sprint(n)), option)
		}
		optionsRange := "0"
		if len(options) > 1 {
			optionsRange = fmt.Sprintf("%d-%d", 0, len(options)-1)
		}
		ans := StringPrompt("Enter number [" + chalk.Cyan.Color(optionsRange) + "] ")

		parsed, err := strconv.ParseInt(ans, 10, 0)

		if err == nil && parsed >= 0 && int(parsed) < len(options) {
			return uint(parsed)
		}
	}
}
