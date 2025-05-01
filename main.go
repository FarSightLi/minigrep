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
	content, err := readFile(filepath)
	if err != nil {
		fmt.Println(err)
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
	if content == "" {
		return
	}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, searchText) {
			fmt.Println(line)
		}
	}
}

func readFile(filepath string) (string, error) {
	// 检查文件是否存在
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	// 检查是否是目录
	if err == nil && info.IsDir() {
		return "", fmt.Errorf("is a directory, not a file")
	}

	// 读取文件内容
	content, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsPermission(err) {
			return "", fmt.Errorf("permission denied")
		}
		return "", err
	}

	return string(content), nil
}
