package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type cmdArgs struct {
	filepath            string
	searchText          string
	isIgnoreCase        bool
	isIncludeLineNumber bool
	afterLine           int
	beforeLine          int
	aroundLine          int
	useRegex            bool
	compile             *regexp.Regexp
}

func main() {
	args := os.Args
	cmdArgs, err := parseArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	var input io.Reader
	// 从输入流中读取
	filepath := cmdArgs.filepath
	if filepath == "" {
		input = os.Stdin
	} else {
		file, e := os.Open(filepath)
		if e != nil {
			if os.IsNotExist(e) {
				fmt.Printf("找不到 %s 文件", filepath)
			} else if os.IsPermission(e) {
				fmt.Printf("权限被拒绝")
			} else {
				fmt.Printf("文件打开失败: %v", e)
			}
			return
		}
		input = file
		// defer 延迟关闭
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("file close file %v", err)
			}
		}(file)
	}
	var printedLine = make(map[int]struct{})

	// 行缓冲区，size为 1+before,最大只需当前行+指定的B参数
	size := 1 + cmdArgs.beforeLine
	buffer := make([]string, 0, size)
	reader := bufio.NewReader(input)
	var lineNum = 0
	scanner := bufio.NewScanner(reader)
	notPrintLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		buffer = addBuffer(line, buffer, size)
		if notPrintLine > 0 {
			notPrintLine--
			fmt.Println(formatLine(cmdArgs, lineNum, line))
		}
		if isMatch(cmdArgs, line) {
			result := getBeforeLine(cmdArgs, buffer, lineNum, printedLine)
			if len(result) > 0 {
				for _, line := range result {
					fmt.Println(line)
				}
			}
			// 不用考虑之前为打印完的行，因此再次匹配后打印的行一定是包含了上次未打印的行
			notPrintLine = cmdArgs.afterLine
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("扫描错误: %v\n", err)
	}
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
	if *around > 0 {
		after = around
		before = around
	}

	// 提前编译正则
	var compile *regexp.Regexp
	if *useRegex {
		c, err := regexp.Compile(searchText)
		if err != nil {
			return cmdArgs{}, fmt.Errorf("非法的正则表达式: %v\n", err)
		}
		compile = c
	}
	return cmdArgs{filepath,
		searchText,
		*isIgnoreCase,
		*isIncludeLineNumber,
		*after,
		*before,
		*around,
		*useRegex,
		compile}, nil
}

/**
 * 搜索文件并打印内容
 */
func isMatch(cmdArgs cmdArgs, line string) bool {
	// 使用正则模式
	if cmdArgs.useRegex {
		return cmdArgs.compile.MatchString(line)
	} else {
		searchText := cmdArgs.searchText
		if cmdArgs.isIgnoreCase {
			searchText = strings.ToLower(searchText)
			line = strings.ToLower(line)
		}
		return strings.Contains(line, searchText)
	}
}

func addBuffer(line string, buffer []string, size int) []string {
	if len(buffer) > size {
		buffer = buffer[1:]
	}
	return append(buffer, line)
}

/**
 * 获得匹配行之前的所有行（已格式化）
 */
func getBeforeLine(cmdArgs cmdArgs, buffedLine []string, currentLineNum int, printedLine map[int]struct{}) (result []string) {
	// 用于计算、定位的下标以buffer位置
	length := len(buffedLine)
	start := length - cmdArgs.beforeLine - 1
	if start < 0 {
		// 即从头打印
		start = 0
	}
	result = make([]string, 0)
	for index := start; index <= length-1; index++ {
		// index对应的当前真实行号 = 当前行号 - 偏移量（length-1-index）
		lineNum := currentLineNum - (length - 1 - index)
		// 如果打印过了就不再打印
		if _, exist := printedLine[lineNum]; exist {
			continue
		}
		result = append(result, formatLine(cmdArgs, lineNum, buffedLine[index]))
		printedLine[lineNum] = struct{}{}
	}
	return result
}

func formatLine(args cmdArgs, index int, line string) string {
	if args.isIncludeLineNumber {
		return fmt.Sprintf("%d:%s", index, line)
	} else {
		return line
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
