package utils

import (
	"testing"
)

func TestDomainMatch(t *testing.T) {
	patterns := []string{"*.cloudfront.net", "baidu.com"}
	for _, d := range []string{"xxx.cloudfront.net", "baidu.com"} {
		if !DomainMatch(d, patterns) {
			t.Error("domain match test failed for", d)
		}
	}
	for _, d := range []string{"x.baidu.com", "xbaidu.com"} {
		if DomainMatch(d, patterns) {
			t.Error("domain shouldn't match for", d)
		}
	}
}
