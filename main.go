package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	cmdArgs, err := parseArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	var content = ""
	// 从输出流中读取
	if cmdArgs.filepath == "" {

	} else {
		contentByte, e := readFile(cmdArgs.filepath)
		if e != nil {
			fmt.Println(e)
			return
		}
		content = contentByte
	}
	searchFile(content, cmdArgs)
}

type cmdArgs struct {
	filepath            string
	searchText          string
	isIgnoreCase        bool
	isIncludeLineNumber bool
	afterLine           int
	beforeLine          int
	aroundLine          int
}

/**
 * 解析命令行参数
 */
func parseArgs(args []string) (cmdArgs, error) {
	// 选项参数用flag解析
	isIgnoreCase := flag.Bool("i", false, "ignore case")
	aroud := flag.Int("C", 0, "aroud line")
	befor := flag.Int("B", 0, "befor line")
	afert := flag.Int("A", 0, "afert line")
	isIncludeLineNumber := flag.Bool("n", false, "include line number")
	err := flag.CommandLine.Parse(args[1:])
	if err != nil {
		return cmdArgs{}, err
	}

	// 非选项参数用普通方式解析
	nonFlagArgs := flag.Args()
	// 只包含了搜索内容
	var searchText string
	var filepath string
	if len(nonFlagArgs) == 1 {
		filepath = ""
		searchText = nonFlagArgs[0]
	} else if len(nonFlagArgs) == 2 {
		filepath = nonFlagArgs[0]
		searchText = nonFlagArgs[1]
	} else {
		return cmdArgs{}, errors.New("参数错误,标准参数只允许有文件路径和搜索内容")
	}
	return cmdArgs{filepath, searchText, *isIgnoreCase, *isIncludeLineNumber, *afert, *befor, *aroud}, nil
}

/**
 * 搜索文件并打印内容
 */
func searchFile(content string, cmdArgs cmdArgs) {
	if content == "" {
		return
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		searchText := cmdArgs.searchText
		if cmdArgs.isIgnoreCase {
			searchText = strings.ToLower(cmdArgs.searchText)
			line = strings.ToLower(line)
		}
		if strings.Contains(line, searchText) {
			printLine(cmdArgs, lines, i)
		}
	}
}

func printLine(cmdArgs cmdArgs, lines []string, i int) {
	var a, b int
	if cmdArgs.aroundLine > 0 {
		a = cmdArgs.aroundLine
		b = cmdArgs.aroundLine
	} else {
		a = cmdArgs.afterLine
		b = cmdArgs.beforeLine
	}
	start := i - a
	if start < 0 {
		start = 0
	}
	end := i + b
	len := len(lines)
	if end >= len {
		end = len - 1
	}
	for index := start; index <= end; index++ {
		if cmdArgs.isIncludeLineNumber {
			fmt.Printf("%d:%s\n", index, lines[index])
		} else {
			fmt.Println(lines[index])
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
