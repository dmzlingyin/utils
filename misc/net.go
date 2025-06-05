package misc

import (
	"net"
	"net/url"
	"regexp"
	"strings"
)

func IsValidURL(input string) bool {
	// 尝试解析URL
	u, err := url.Parse(input)
	if err != nil {
		return false
	}

	// 检查必需组件
	if u.Scheme == "" || u.Host == "" {
		return false
	}

	// 验证协议是否为常见类型
	schemes := []string{"http", "https", "ftp", "sftp"}
	validScheme := false
	for _, s := range schemes {
		if u.Scheme == s {
			validScheme = true
			break
		}
	}
	if !validScheme {
		return false
	}

	// 从Host中分离主机名和端口号
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		// 如果分离失败，说明没有端口号，直接使用Host
		host = u.Host
	}

	// 验证主机名是否有效
	return isValidHost(host)
}

// isValidHost 验证主机名是否有效
func isValidHost(host string) bool {
	// 检查是否为空
	if host == "" {
		return false
	}

	// 检查是否为localhost
	if host == "localhost" {
		return true
	}

	// 检查是否为IPv4地址
	if net.ParseIP(host) != nil {
		return true
	}

	// 检查是否为IPv6地址（可能包含方括号）
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		ipv6 := host[1 : len(host)-1]
		if net.ParseIP(ipv6) != nil {
			return true
		}
	}

	// 验证域名格式
	return isValidDomain(host)
}

// isValidDomain 验证域名格式是否有效
func isValidDomain(domain string) bool {
	// 域名长度检查
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// 域名不能以点开始或结束
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	// 域名必须包含至少一个点（除非是localhost）
	if !strings.Contains(domain, ".") {
		return false
	}

	// 分割域名标签
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}

	// 验证每个标签
	labelRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if !labelRegex.MatchString(label) {
			return false
		}
	}

	return true
}
