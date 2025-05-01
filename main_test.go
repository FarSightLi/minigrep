package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected cmdArgs
		err      error
	}{{
		name: "test1",
		args: []string{"-i", "-a=2", "file.txt", "hello"},
		expected: cmdArgs{
			filepath:            "file.txt",
			searchText:          "hello",
			isIgnoreCase:        true,
			isIncludeLineNumber: false,
			aroundLine:          2,
			beforeLine:          0,
			afterLine:           0,
		},
		err: nil,
	}, {
		name: "test2",
		args: []string{"-n", "-B=3", "file.txt", "world"},
		expected: cmdArgs{
			filepath:            "file.txt",
			searchText:          "world",
			isIgnoreCase:        false,
			isIncludeLineNumber: true,
			aroundLine:          0,
			beforeLine:          3,
			afterLine:           0,
		},
		err: nil,
	}, {
		name: "test3",
		args: []string{"-A=1", "-B=2", "file.txt", "test"},
		expected: cmdArgs{
			filepath:            "file.txt",
			searchText:          "test",
			isIgnoreCase:        false,
			isIncludeLineNumber: false,
			aroundLine:          0,
			beforeLine:          2,
			afterLine:           1,
		},
		err: nil,
	}, {
		name: "test4",
		args: []string{"invalid_flag", "file.txt", "error"},
		expected: cmdArgs{
			filepath:   "file.txt",
			searchText: "error",
		},
		err: flag.ErrHelp,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
			result, err := parseArgs(test.args)
			if err != test.err {
				t.Errorf("expected error %v, got %v", test.err, err)
			}
			if result.filepath != test.expected.filepath ||
				result.searchText != test.expected.searchText ||
				result.isIgnoreCase != test.expected.isIgnoreCase ||
				result.isIncludeLineNumber != test.expected.isIncludeLineNumber ||
				result.aroundLine != test.expected.aroundLine ||
				result.beforeLine != test.expected.beforeLine ||
				result.afterLine != test.expected.afterLine {
				t.Errorf("expected %+v, got %+v", test.expected, result)
			}
		})
	}
}

func TestSearchFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		cmdArgs  cmdArgs
		expected []int
	}{{
		name:    "test1",
		content: "line1\nline2\nline3",
		cmdArgs: cmdArgs{
			searchText: "line",
		},
		expected: []int{0, 1, 2},
	}, {
		name:    "test2",
		content: "Hello\nWorld\nhello world",
		cmdArgs: cmdArgs{
			searchText:   "hello",
			isIgnoreCase: true,
		},
		expected: []int{0, 2},
	}, {
		name:    "test3",
		content: "apple\nbanana\ncherry\ndate\nelderberry",
		cmdArgs: cmdArgs{
			searchText:   "berry",
			isIgnoreCase: false,
		},
		expected: []int{4},
	}, {
		name:    "test4",
		content: "apple\nbanana\ncherry\ndate\nelderberry",
		cmdArgs: cmdArgs{
			searchText:   "peach",
			isIgnoreCase: false,
		},
		expected: []int{},
	}, {
		name:    "test5",
		content: "apple\nbanana\ncherry\ndate\nelderberry",
		cmdArgs: cmdArgs{
			searchText:          "NA",
			isIgnoreCase:        false,
			beforeLine:          1,
			afterLine:           1,
			isIncludeLineNumber: true,
		},
		expected: []int{1, 2, 3},
	}, {
		name:    "test6",
		content: "apple\nbanana\ncherry\ndate\nelderberry",
		cmdArgs: cmdArgs{
			searchText:          "apple",
			isIgnoreCase:        false,
			aroundLine:          1,
			isIncludeLineNumber: true,
		},
		expected: []int{0, 1, 2},
	}, {
		name:    "test7",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "grape",
			isIgnoreCase:        false,
			aroundLine:          2,
			isIncludeLineNumber: true,
		},
		expected: []int{5, 6},
	}, {
		name:    "test8",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "apple",
			isIgnoreCase:        false,
			aroundLine:          3,
			isIncludeLineNumber: true,
		},
		expected: []int{0, 1, 2, 3},
	}, {
		name:    "test9",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "fig",
			isIgnoreCase:        false,
			aroundLine:          3,
			isIncludeLineNumber: true,
		},
		expected: []int{3, 4, 5, 6},
	}, {
		name:    "test10",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "grape",
			isIgnoreCase:        false,
			aroundLine:          3,
			isIncludeLineNumber: true,
		},
		expected: []int{3, 4, 5, 6},
	}, {
		name:    "test11",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "banana",
			isIgnoreCase:        false,
			aroundLine:          3,
			isIncludeLineNumber: true,
		},
		expected: []int{0, 1, 2, 3},
	}, {
		name:    "test12",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "cherry",
			isIgnoreCase:        false,
			aroundLine:          2,
			isIncludeLineNumber: true,
		},
		expected: []int{0, 1, 2, 3, 4},
	}, {
		name:    "test13",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "date",
			isIgnoreCase:        false,
			aroundLine:          2,
			isIncludeLineNumber: true,
		},
		expected: []int{1, 2, 3, 4},
	}, {
		name:    "test14",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "elderberry",
			isIgnoreCase:        false,
			aroundLine:          2,
			isIncludeLineNumber: true,
		},
		expected: []int{2, 3, 4, 5},
	}, {
		name:    "test15",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "fig",
			isIgnoreCase:        false,
			aroundLine:          2,
			isIncludeLineNumber: true,
		},
		expected: []int{3, 4, 5, 6},
	}, {
		name:    "test16",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "date",
			isIgnoreCase:        false,
			beforeLine:          2,
			afterLine:           2,
			isIncludeLineNumber: true,
		},
		expected: []int{1, 2, 3, 4},
	}, {
		name:    "test17",
		content: "apple\nbanana\ncherry\ndate\nelderberry\nfig\ngrape",
		cmdArgs: cmdArgs{
			searchText:          "grape",
			isIgnoreCase:        false,
			beforeLine:          6,
			afterLine:           0,
			isIncludeLineNumber: true,
		},
		expected: []int{5, 6},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 捕获输出
			var output bytes.Buffer
			oldStdout := os.Stdout
			customWriter := bufio.NewWriter(&output)
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = oldStdout }()
			_, _ = r, w

			// 将内容分割成行
			lines := strings.Split(test.content, "\n")

			searchFile(test.content, test.cmdArgs)
			customWriter.Flush()

			// 验证输出
			outputLines := strings.Split(strings.TrimRight(output.String(), "\n"), "\n")
			for i, line := range outputLines {
				if i < len(test.expected) {
					expectedLine := test.expected[i]
					if test.cmdArgs.isIncludeLineNumber {
						expectedPrefix := fmt.Sprintf("%d:", expectedLine)
						if !strings.HasPrefix(line, expectedPrefix) {
							t.Errorf("expected line %d to start with '%s', got '%s'", i, expectedPrefix, line)
						}
						// 检查行号后的文本是否正确
						if expectedLine >= 0 && expectedLine < len(lines) && !strings.Contains(line, lines[expectedLine]) {
							t.Errorf("expected line %d to contain '%s', got '%s'", i, lines[expectedLine], line)
						}
					} else {
						// 不带行号的情况
						if expectedLine >= 0 && expectedLine < len(lines) && !strings.Contains(line, lines[expectedLine]) {
							t.Errorf("expected line %d to contain '%s', got '%s'", i, lines[expectedLine], line)
						}
					}
				} else if len(test.expected) > 0 {
					t.Errorf("extra line %d in output: %s", i, line)
				}
			}
			if len(outputLines) == 0 && len(test.expected) > 0 {
				t.Errorf("expected output but got none")
			}
		})
	}
}

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

// TestReadFile 测试读取文件功能
func TestReadFile(t *testing.T) {
	// 创建临时测试文件
	dir := t.TempDir()
	testFile := dir + "/test.txt"
	data := "hello world\nthis is a test\ngoodbye world"
	os.WriteFile(testFile, []byte(data), 0644)

	// 测试正常读取文件
	content, err := readFile(testFile)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if content != data {
		t.Errorf("Expected content '%s', but got '%s'", data, content)
	}

	// 测试读取不存在的文件
	_, err = readFile(dir + "/nonexistent.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, but got none")
	} else if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("Expected error message containing 'file not found', but got '%v'", err)
	}

	// 测试读取目录
	_, err = readFile(dir)
	if err == nil {
		t.Error("Expected error for reading directory, but got none")
	} else if !strings.Contains(err.Error(), "is a directory") {
		t.Errorf("Expected error message containing 'is a directory', but got '%v'", err)
	}
}
