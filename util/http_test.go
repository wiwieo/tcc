package util

import "testing"

func TestURLRewrite(t *testing.T) {
	if URLRewrite("/api/(.*)/actions/(.*)", "/api/1234/actions/hello", "/api/v1/$1/actions/$2") != "/api/v1/1234/actions/hello" {
		t.Fail()
	}
}
