package minigrep

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestParseArgs 测试解析命令行参数的功能
func TestParseArgs(t *testing.T) {
	// 保存原始的命令行参数
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// 测试正常情况
	os.Args = []string{"minigrep", "test.txt", "hello"}
	filepath, searchText, err := parseArgs(os.Args)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if filepath != "test.txt" {
		t.Errorf("Expected filepath to be 'test.txt', but got %v", filepath)
	}
	if searchText != "hello" {
		t.Errorf("Expected searchText to be 'hello', but got %v", searchText)
	}

	// 测试参数数量错误的情况
	os.Args = []string{"minigrep", "test.txt"}
	_, _, err = parseArgs(os.Args)
	if err == nil {
		t.Error("Expected error for incorrect number of arguments, but got none")
	} else if !strings.Contains(err.Error(), "illegal arguments") {
		t.Errorf("Expected error message containing 'illegal arguments', but got '%v'", err)
	}
}

// TestSearchFile 测试在文件内容中搜索文本的功能
func TestSearchFile(t *testing.T) {
	// 测试包含搜索文本的情况
	content := "hello world\nthis is a test\ngoodbye world"
	searchText := "hello"
	expected := "hello world\n"

	// 重定向标准输出以便捕获
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	searchFile(content, searchText)
	w.Close()
	os.Stdout = stdout

	// 读取输出
	result, _ := io.ReadAll(r)
	if string(result) != expected {
		t.Errorf("Expected output '%s', but got '%s'", expected, string(result))
	}

	// 测试不包含搜索文本的情况
	content = "hello world\nthis is a test\ngoodbye world"
	searchText = "missing"
	r, w, _ = os.Pipe()
	stdout = os.Stdout
	os.Stdout = w

	searchFile(content, searchText)
	w.Close()
	os.Stdout = stdout

	// 读取输出
	result, _ = io.ReadAll(r)
	if len(result) > 0 {
		t.Errorf("Expected no output, but got '%s'", string(result))
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
