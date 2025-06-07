package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
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
	// 从输入流中读取
	if cmdArgs.filepath == "" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("读取输入流出错")
			return
		}
		content = string(bytes)
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
	useRegex            bool
}

/**
 * 解析命令行参数
 */
func parseArgs(args []string) (cmdArgs, error) {
	// 选项参数用flag解析
	isIgnoreCase := flag.Bool("i", false, "ignore case")
	around := flag.Int("a", 0, "around line")
	before := flag.Int("B", 0, "before line")
	after := flag.Int("A", 0, "after line")
	isIncludeLineNumber := flag.Bool("n", false, "include line number")
	useRegex := flag.Bool("e", false, "use regex module")
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
	return cmdArgs{filepath,
		searchText,
		*isIgnoreCase,
		*isIncludeLineNumber,
		*after,
		*before,
		*around,
		*useRegex}, nil
}

/**
 * 搜索文件并打印内容
 */
func searchFile(content string, cmdArgs cmdArgs) {
	if content == "" {
		return
	}
	searchText := cmdArgs.searchText
	// 提前编译正则
	compile, err := regexp.Compile(searchText)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "非法的正则表达式: %v\n", err)
		return
	}
	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		// 使用正则模式
		if cmdArgs.useRegex {
			if compile.MatchString(line) {
				printLine(cmdArgs, lines, lineNum)
			}
		} else {
			if cmdArgs.isIgnoreCase {
				searchText = strings.ToLower(cmdArgs.searchText)
				line = strings.ToLower(line)
			}
			if strings.Contains(line, searchText) {
				printLine(cmdArgs, lines, lineNum)
			}
		}
	}
}

/**
 * 根据参数打印匹配上的行
 */
func printLine(cmdArgs cmdArgs, lines []string, lineNum int) {
	var a, b int
	if cmdArgs.aroundLine > 0 {
		a = cmdArgs.aroundLine
		b = cmdArgs.aroundLine
	} else {
		a = cmdArgs.afterLine
		b = cmdArgs.beforeLine
	}
	start := lineNum - a
	if start < 0 {
		start = 0
	}
	end := lineNum + b
	length := len(lines)
	if end >= length {
		end = length - 1
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
		return "", fmt.Errorf("找不到 %s 文件", filepath)
	}

	// 检查是否是目录
	if err == nil && info.IsDir() {
		return "", fmt.Errorf("%s 是一个目录而非文件", filepath)
	}

	// 读取文件内容
	content, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsPermission(err) {
			return "", errors.New("权限被拒绝")
		}
		return "", err
	}

	return string(content), nil
}
