package validation

import (
	"strings"
	"testing"
)

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		// Valid tags
		{name: "simple version", tag: "v1.0.0", wantErr: false},
		{name: "semantic version", tag: "v2.3.4", wantErr: false},
		{name: "date-based release", tag: "release-2025-11-15", wantErr: false},
		{name: "rc with slash", tag: "rc/1.2", wantErr: false},
		{name: "feature branch style", tag: "feature/x/y", wantErr: false},
		{name: "alpha tag", tag: "v1.0.0-alpha", wantErr: false},
		{name: "beta tag", tag: "v1.0.0-beta.1", wantErr: false},
		{name: "numeric only", tag: "123", wantErr: false},
		{name: "with underscore", tag: "v1_0_0", wantErr: false},
		{name: "with dash", tag: "v1-0-0", wantErr: false},
		{name: "mixed case", tag: "V1.0.0", wantErr: false},
		{name: "single char", tag: "v", wantErr: false},
		{name: "max length", tag: strings.Repeat("a", MaxTagLength), wantErr: false},

		// Invalid tags - injection attempts
		{name: "refs/branches prefix", tag: "refs/branches/master", wantErr: true},
		{name: "refs/heads prefix", tag: "refs/heads/main", wantErr: true},
		{name: "refs/tags prefix", tag: "refs/tags/v1.0.0", wantErr: true},
		{name: "parent directory", tag: "../master", wantErr: true},
		{name: "parent in middle", tag: "foo/../bar", wantErr: true},
		{name: "current directory", tag: "./master", wantErr: true},
		{name: "hidden file start", tag: ".hidden", wantErr: true},

		// Invalid tags - git rules violations
		{name: "double dot", tag: "v1..0", wantErr: true},
		{name: "at brace", tag: "bad@{1}", wantErr: true},
		{name: "trailing dot", tag: "trailing.", wantErr: true},
		{name: "trailing slash", tag: "trailing/", wantErr: true},
		{name: "trailing lock", tag: "foo.lock", wantErr: true},
		{name: "start with dash", tag: "-start", wantErr: true},
		{name: "start with slash", tag: "/start", wantErr: true},
		{name: "single at", tag: "@", wantErr: true},
		{name: "with space", tag: "space tag", wantErr: true},
		{name: "with tilde", tag: "v1~2", wantErr: true},
		{name: "with caret", tag: "v1^2", wantErr: true},
		{name: "with colon", tag: "foo:bar", wantErr: true},
		{name: "with question", tag: "foo?bar", wantErr: true},
		{name: "with asterisk", tag: "foo*bar", wantErr: true},
		{name: "with bracket", tag: "foo[bar", wantErr: true},
		{name: "with backslash", tag: "foo\\bar", wantErr: true},
		{name: "empty string", tag: "", wantErr: true},
		{name: "double slash", tag: "foo//bar", wantErr: true},
		{name: "segment ends dot", tag: "foo./bar", wantErr: true},
		{name: "segment starts dot", tag: "foo/.bar", wantErr: true},
		{name: "control character", tag: "foo\x00bar", wantErr: true},
		{name: "newline", tag: "foo\nbar", wantErr: true},
		{name: "tab", tag: "foo\tbar", wantErr: true},
		{name: "over max length", tag: strings.Repeat("a", MaxTagLength+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
			if err != nil && err != ErrInvalidTagFormat {
				t.Errorf("ValidateTag(%q) returned unexpected error type: %v", tt.tag, err)
			}
		})
	}
}

func TestValidateTagRealWorldExamples(t *testing.T) {
	// Real-world tag examples that should be valid
	validTags := []string{
		"v1.0.0",
		"v2.3.4-rc1",
		"release-1.2.3",
		"1.0.0",
		"v1.0.0-alpha.1",
		"v1.0.0-beta",
		"20251115",
		"v1.0.0+build.123",
		"chart-1.2.3",
		"operator-v1.0.0",
		"pkg/apis/monitoring/v0.65.2",
	}

	for _, tag := range validTags {
		t.Run("valid_"+tag, func(t *testing.T) {
			if err := ValidateTag(tag); err != nil {
				t.Errorf("ValidateTag(%q) should be valid, got error: %v", tag, err)
			}
		})
	}
}
