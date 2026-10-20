package nginx

import (
	"regexp"
	"strings"
)

// LuaDirective is one OpenResty Lua directive found in a config block.
type LuaDirective struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ParsedLocation is a parsed nginx location block.
type ParsedLocation struct {
	Path          string         `json:"path"`
	ProxyPass     string         `json:"proxyPass,omitempty"`
	Root          string         `json:"root,omitempty"`
	Return        string         `json:"return,omitempty"`
	LuaDirectives []LuaDirective `json:"luaDirectives"`
}

// ParsedSite is a parsed nginx server block.
type ParsedSite struct {
	Name          string           `json:"name"`
	FilePath      string           `json:"filePath"`
	ServerName    []string         `json:"serverName"`
	Listen        []string         `json:"listen"`
	SSL           bool             `json:"ssl"`
	HTTP2         bool             `json:"http2"`
	Locations     []ParsedLocation `json:"locations"`
	AccessLog     string           `json:"accessLog,omitempty"`
	ErrorLog      string           `json:"errorLog,omitempty"`
	Root          string           `json:"root,omitempty"`
	ProxyPass     string           `json:"proxyPass,omitempty"`
	LuaDirectives []LuaDirective   `json:"luaDirectives"`
	Content       string           `json:"content"`
}

var (
	serverBlockPattern   = regexp.MustCompile(`\bserver\b[^{}]*\{`)
	locationBlockPattern = regexp.MustCompile(`\blocation\b[^{}]*\{`)
	serverNamePattern    = regexp.MustCompile(`server_name\s+([^;]+);`)
	listenPattern        = regexp.MustCompile(`listen\s+([^;]+);`)
	accessLogPattern     = regexp.MustCompile(`access_log\s+([^;]+);`)
	errorLogPattern      = regexp.MustCompile(`error_log\s+([^;]+);`)
	rootPattern          = regexp.MustCompile(`root\s+([^;]+);`)
	proxyPassPattern     = regexp.MustCompile(`proxy_pass\s+([^;]+);`)
	returnPattern        = regexp.MustCompile(`return\s+([^;]+);`)
)

// LuaDirectives lists the OpenResty lifecycle directives the manager understands.
var LuaDirectives = []string{
	"init_by_lua", "init_by_lua_block", "init_by_lua_file",
	"init_worker_by_lua", "init_worker_by_lua_block", "init_worker_by_lua_file",
	"content_by_lua", "content_by_lua_block", "content_by_lua_file",
	"access_by_lua", "access_by_lua_block", "access_by_lua_file",
	"rewrite_by_lua", "rewrite_by_lua_block", "rewrite_by_lua_file",
	"log_by_lua", "log_by_lua_block", "log_by_lua_file",
	"header_filter_by_lua", "header_filter_by_lua_block", "header_filter_by_lua_file",
	"body_filter_by_lua", "body_filter_by_lua_block", "body_filter_by_lua_file",
	"set_by_lua", "set_by_lua_block", "set_by_lua_file",
	"balancer_by_lua", "balancer_by_lua_block", "balancer_by_lua_file",
	"ssl_certificate_by_lua", "ssl_certificate_by_lua_block", "ssl_certificate_by_lua_file",
	"server_rewrite_by_lua",
}

// ParseSites extracts every server block from a config file.
func ParseSites(content, filePath string) []ParsedSite {
	blocks := extractBlocks(content, serverBlockPattern)
	sites := make([]ParsedSite, 0, len(blocks))
	for _, block := range blocks {
		body := block.body
		serverNames := extractValues(body, serverNamePattern)
		listens := extractValues(body, listenPattern)

		ssl := strings.Contains(body, "ssl on")
		http2 := false
		for _, listen := range listens {
			if strings.Contains(listen, "ssl") {
				ssl = true
			}
			if strings.Contains(listen, "http2") {
				http2 = true
			}
		}

		name := "unnamed"
		if len(serverNames) > 0 {
			name = serverNames[0]
		}

		sites = append(sites, ParsedSite{
			Name:          name,
			FilePath:      filePath,
			ServerName:    serverNames,
			Listen:        listens,
			SSL:           ssl,
			HTTP2:         http2,
			Locations:     parseLocations(body),
			AccessLog:     extractSingle(body, accessLogPattern),
			ErrorLog:      extractSingle(body, errorLogPattern),
			Root:          extractSingle(body, rootPattern),
			ProxyPass:     extractSingle(body, proxyPassPattern),
			LuaDirectives: ExtractLuaDirectives(body),
			Content:       body,
		})
	}
	return sites
}

type rawBlock struct {
	header string
	body   string
}

// extractBlocks collects `<keyword ...> { ... }` blocks, honouring nested braces.
func extractBlocks(content string, pattern *regexp.Regexp) []rawBlock {
	blocks := make([]rawBlock, 0)
	offset := 0
	for offset < len(content) {
		match := pattern.FindStringIndex(content[offset:])
		if match == nil {
			break
		}
		openIndex := offset + match[1] - 1
		closeIndex := findMatchingBrace(content, openIndex)
		if closeIndex == -1 {
			break
		}
		blocks = append(blocks, rawBlock{
			header: strings.TrimSpace(content[offset+match[0] : openIndex]),
			body:   content[openIndex+1 : closeIndex],
		})
		offset = closeIndex + 1
	}
	return blocks
}

func findMatchingBrace(content string, openIndex int) int {
	depth := 0
	for i := openIndex; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func extractValues(content string, pattern *regexp.Regexp) []string {
	match := pattern.FindStringSubmatch(content)
	if match == nil {
		return []string{}
	}
	return strings.Fields(match[1])
}

func extractSingle(content string, pattern *regexp.Regexp) string {
	match := pattern.FindStringSubmatch(content)
	if match == nil {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func parseLocations(serverBlock string) []ParsedLocation {
	blocks := extractBlocks(serverBlock, locationBlockPattern)
	locations := make([]ParsedLocation, 0, len(blocks))
	locationPrefix := regexp.MustCompile(`^location\s+`)
	for _, block := range blocks {
		path := strings.TrimSpace(locationPrefix.ReplaceAllString(block.header, ""))
		if path == "" {
			path = "/"
		}
		locations = append(locations, ParsedLocation{
			Path:          path,
			ProxyPass:     extractSingle(block.body, proxyPassPattern),
			Root:          extractSingle(block.body, rootPattern),
			Return:        extractSingle(block.body, returnPattern),
			LuaDirectives: ExtractLuaDirectives(block.body),
		})
	}
	return locations
}

// ExtractLuaDirectives finds the OpenResty Lua directives inside a block.
func ExtractLuaDirectives(content string) []LuaDirective {
	directives := make([]LuaDirective, 0)
	for _, name := range LuaDirectives {
		blockPattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(name) + `_block\s*((?:[^}]|\}[^}])*)\}`)
		if match := blockPattern.FindStringSubmatch(content); match != nil {
			directives = append(directives, LuaDirective{Type: name + "_block", Value: strings.TrimSpace(match[1])})
			continue
		}
		filePattern := regexp.MustCompile(regexp.QuoteMeta(name) + `_file\s+([^;]+);`)
		if match := filePattern.FindStringSubmatch(content); match != nil {
			directives = append(directives, LuaDirective{Type: name + "_file", Value: strings.TrimSpace(match[1])})
		}
	}
	return directives
}
