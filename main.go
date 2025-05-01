package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	filepath, searchText, err := parseArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	content, e := os.ReadFile(filepath)
	if e != nil {
		fmt.Println(e)
		return
	}
	searchFile(string(content), searchText)
}

func parseArgs(args []string) (filepath string, searchText string, error error) {
	if len(args) != 3 {
		return "", "", fmt.Errorf("illegal arguments, should be filepath and searchText")
	}

	// 文件路径
	filepath = args[1]
	searchText = args[2]
	return filepath, searchText, nil
}
func searchFile(content string, searchText string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, searchText) {
			fmt.Println(line)
		}
	}
}
