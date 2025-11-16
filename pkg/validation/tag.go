package validation

import (
	"errors"
	"strings"
)

const (
	// MaxTagLength is the maximum allowed length for a git tag.
	MaxTagLength = 128
)

var (
	// ErrInvalidTagFormat is returned when a tag doesn't meet git ref requirements.
	ErrInvalidTagFormat = errors.New("invalid tag format")
)

// ValidateTag validates that a tag name conforms to git ref naming rules
// and security requirements. It allows slashes in tag names (e.g., release/v1.0).
//
// Rules enforced:
// - Length must be <= MaxTagLength
// - No control characters (< 0x20 or 0x7F)
// - No space, ~, ^, :, ?, *, [, \
// - No consecutive dots (..)
// - No @{ sequence
// - Cannot start with - or /
// - Cannot end with /, ., or .lock
// - Cannot be a single @ character
// - Path components cannot start or end with .
func ValidateTag(tag string) error {
	if tag == "" {
		return ErrInvalidTagFormat
	}
	if len(tag) > MaxTagLength {
		return ErrInvalidTagFormat
	}
	if tag == "@" {
		return ErrInvalidTagFormat
	}
	if tag[0] == '-' || tag[0] == '/' {
		return ErrInvalidTagFormat
	}
	if tag[len(tag)-1] == '/' || tag[len(tag)-1] == '.' {
		return ErrInvalidTagFormat
	}
	if strings.HasPrefix(tag, "refs/") {
		return ErrInvalidTagFormat
	}
	if strings.HasSuffix(tag, ".lock") {
		return ErrInvalidTagFormat
	}
	if strings.Contains(tag, "..") {
		return ErrInvalidTagFormat
	}
	if strings.Contains(tag, "@{") {
		return ErrInvalidTagFormat
	}

	segments := strings.Split(tag, "/")
	for _, segment := range segments {
		if segment == "" {
			return ErrInvalidTagFormat
		}
		if segment[0] == '.' || segment[len(segment)-1] == '.' {
			return ErrInvalidTagFormat
		}

		for i := 0; i < len(segment); i++ {
			c := segment[i]
			if c < 0x20 || c == 0x7F {
				return ErrInvalidTagFormat
			}
			switch c {
			case ' ', '~', '^', ':', '?', '*', '[', '\\':
				return ErrInvalidTagFormat
			}
		}
	}

	return nil
}
