package main

import (
	"flag"
	"fmt"
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
		args: []string{"main.go", "-i", "-C=2", "file.txt", "hello"},
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
		args: []string{"main.go", "-n", "-B=3", "file.txt", "world"},
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
		args: []string{"main.go", "-A=1", "-B=2", "file.txt", "test"},
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
		args: []string{"main.go", "invalid_flag", "file.txt", "error"},
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
			if test.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			} else if test.err != nil && err == nil {
				// 针对Test4
				t.Errorf("expected error %v, but got nil", test.err)
			}
			if err != nil {
				return
			}
			if result.filepath != test.expected.filepath ||
				result.searchText != test.expected.searchText ||
				result.isIgnoreCase != test.expected.isIgnoreCase ||
				result.isIncludeLineNumber != test.expected.isIncludeLineNumber ||
				result.aroundLine != test.expected.aroundLine ||
				result.beforeLine != test.expected.beforeLine ||
				result.afterLine != test.expected.afterLine {
				t.Errorf("expected:\n%+v\n\ngot:\n%+v", test.expected, result)
			}
		})
	}
}

// 定义测试用的参数结构体
type testCmdArgs struct {
	afterLine  int
	beforeLine int
	pattern    string
	matchFunc  func(string) bool
}

// 测试场景数据结构
type testScenario struct {
	name           string
	input          string
	cmdArgs        testCmdArgs
	expectedOutput []string
}

func TestMatchLines(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		cmdArgs  cmdArgs
		expected []string
	}{
		{
			name:    "普通匹配",
			content: "apple banana\ncherry apple\ndate",
			cmdArgs: cmdArgs{
				searchText: "apple",
			},
			expected: []string{"apple banana", "cherry apple"},
		},
		{
			name:    "忽略大小写",
			content: "Apple banana\nCHERRY APPLE",
			cmdArgs: cmdArgs{
				searchText:   "apple",
				isIgnoreCase: true,
			},
			expected: []string{"Apple banana", "CHERRY APPLE"},
		},
		{
			name:    "正则匹配以 a 开头",
			content: "apple banana\ngrape apple\nant",
			cmdArgs: cmdArgs{
				searchText: "^[aA]",
				useRegex:   true,
			},
			expected: []string{"apple banana", "ant"},
		},
		{
			name:    "上下文两行",
			content: "line1\nline2\napple\nline4\nline5",
			cmdArgs: cmdArgs{
				searchText: "apple",
				aroundLine: 2,
			},
			expected: []string{"line1", "line2", "apple", "line4", "line5"},
		},
		{
			name:    "显示行号",
			content: "hello\nworld\nhello again",
			cmdArgs: cmdArgs{
				searchText:          "hello",
				isIncludeLineNumber: true,
			},
			expected: []string{"0:hello", "2:hello again"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO 实现自动测试
			//reader := strings.NewReader(tt.content)
			//scanner := bufio.NewScanner(reader)
			//var results []string
			//for scanner.Scan() {
			//	text := scanner.Text()
			//	results = append(results, MatchLines(text, tt.cmdArgs, make(map[int]struct{}))...)
			//}
			//
			//if len(results) != len(tt.expected) {
			//	t.Errorf("预期 %d 行，实际 %d 行", len(tt.expected), len(results))
			//	return
			//}
			//for i := range results {
			//	if results[i] != tt.expected[i] {
			//		t.Errorf("第 %d 行期望 %q，实际 %q", i, tt.expected[i], results[i])
			//	}
			//}
		})
	}
}

func TestSet(t *testing.T) {
	// 创建一个集合
	printedLine := make(map[int]struct{})

	// 添加元素
	lineNum := 42
	if _, exists := printedLine[lineNum]; !exists {
		printedLine[lineNum] = struct{}{}
		fmt.Printf("添加了 lineNum: %d\n", lineNum)
	} else {
		fmt.Printf("lineNum: %d 已存在\n", lineNum)
	}

	// 再次尝试添加同一个元素
	if _, exists := printedLine[lineNum]; !exists {
		printedLine[lineNum] = struct{}{}
		fmt.Printf("添加了 lineNum: %d\n", lineNum)
	} else {
		fmt.Printf("lineNum: %d 已存在\n", lineNum)
	}

	// 删除元素
	delete(printedLine, lineNum)
	fmt.Println("已删除 lineNum:", lineNum)

	// 再次检查是否还存在
	if _, exists := printedLine[lineNum]; !exists {
		fmt.Printf("lineNum: %d 已被删除\n", lineNum)
	}
}

func TestAppendList(t *testing.T) {
	result := make([]string, 0)
	result = append(result, "hello")
	fmt.Printf("结果: %v\n", result)
}
