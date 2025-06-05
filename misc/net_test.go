package misc

import (
	"strings"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 有效的 URL 测试用例
		{
			name:     "有效的 HTTP URL",
			input:    "http://www.example.com",
			expected: true,
		},
		{
			name:     "有效的 HTTPS URL",
			input:    "https://www.example.com",
			expected: true,
		},
		{
			name:     "有效的 FTP URL",
			input:    "ftp://ftp.example.com",
			expected: true,
		},
		{
			name:     "有效的 SFTP URL",
			input:    "sftp://sftp.example.com",
			expected: true,
		},
		{
			name:     "包含路径的有效 URL",
			input:    "https://www.example.com/path/to/resource",
			expected: true,
		},
		{
			name:     "包含查询参数的有效 URL",
			input:    "https://www.example.com/search?q=test&page=1",
			expected: true,
		},
		{
			name:     "包含端口号的有效 URL",
			input:    "https://www.example.com:8080",
			expected: true,
		},
		{
			name:     "localhost URL",
			input:    "http://localhost",
			expected: true,
		},
		{
			name:     "带端口的 localhost URL",
			input:    "http://localhost:3000",
			expected: true,
		},
		{
			name:     "IP 地址 URL",
			input:    "http://192.168.1.1",
			expected: true,
		},
		{
			name:     "子域名 URL",
			input:    "https://api.example.com",
			expected: true,
		},

		// 无效的 URL 测试用例
		{
			name:     "缺少协议",
			input:    "www.example.com",
			expected: false,
		},
		{
			name:     "缺少主机名",
			input:    "http://",
			expected: false,
		},
		{
			name:     "无效的协议",
			input:    "invalid://www.example.com",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
		{
			name:     "只有协议",
			input:    "http://",
			expected: false,
		},
		{
			name:     "无效的主机名格式",
			input:    "http://invalidhost",
			expected: false,
		},
		{
			name:     "包含空格的 URL",
			input:    "http://www.exa mple.com",
			expected: false,
		},
		{
			name:     "格式错误的 URL",
			input:    "not-a-url",
			expected: false,
		},
		{
			name:     "只有路径",
			input:    "/path/to/resource",
			expected: false,
		},
		{
			name:     "javascript 协议",
			input:    "javascript:alert('xss')",
			expected: false,
		},
		{
			name:     "file 协议",
			input:    "file:///etc/passwd",
			expected: false,
		},
		{
			name:     "mailto 协议",
			input:    "mailto:test@example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkIsValidURL(b *testing.B) {
	testURL := "https://www.example.com/path?query=value"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsValidURL(testURL)
	}
}

// 测试边界情况
func TestIsValidURL_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "非常长的 URL",
			input:    "https://www.example.com/" + strings.Repeat("a", 1000),
			expected: true,
		},
		{
			name:     "包含特殊字符的域名",
			input:    "https://sub-domain.example-site.com",
			expected: true,
		},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}
