package nginx

import "regexp"

var (
	regexpHTTPBlock = regexp.MustCompile(`(?s)http\s*\{`)
	regexpTriple    = regexp.MustCompile(`(\d+)\s+(\d+)\s+(\d+)`)
	regexpIntCache  = map[string]*regexp.Regexp{}
)

func matchInt(content, pattern string) int64 {
	re, ok := regexpIntCache[pattern]
	if !ok {
		re = regexp.MustCompile(pattern)
		regexpIntCache[pattern] = re
	}
	match := re.FindStringSubmatch(content)
	if match == nil {
		return 0
	}
	return atoi(match[1])
}
