package main

import (
	"flag"
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
			results := MatchLines(tt.content, tt.cmdArgs, make(map[int]struct{}))
			if len(results) != len(tt.expected) {
				t.Errorf("预期 %d 行，实际 %d 行", len(tt.expected), len(results))
				return
			}
			for i := range results {
				if results[i] != tt.expected[i] {
					t.Errorf("第 %d 行期望 %q，实际 %q", i, tt.expected[i], results[i])
				}
			}
		})
	}
}
