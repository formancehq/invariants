package domain

import "strings"

// isIdentifierChar returns true if the rune is allowed in a Formance
// identifier (ledger name, metadata key). The charset is intentionally
// narrow (identifier-style): letters, digits, colon, dot, underscore,
// hyphen. Matches the spirit of Kubernetes labels, OpenTelemetry attribute
// keys, and gRPC metadata header names without committing to any single
// platform's exact rule.
//
// Account addresses use a slightly tighter charset (no dot) because the
// colon is the conventional hierarchy separator there.
func isIdentifierChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ':' || r == '.' || r == '_' || r == '-'
}

// ValidateMetadataKey checks that a metadata key matches the identifier
// charset [a-zA-Z0-9._:-]. The strict charset excludes whitespace, free
// punctuation, and UTF-8 multibyte sequences; this also rules out null
// bytes that would corrupt null-terminated key encodings in storage
// backends, and keeps keys safe to render in CLI output, logs, URLs and
// HTTP/2 metadata trailers.
func ValidateMetadataKey(key string) error {
	if key == "" {
		return ErrMetadataKeyEmpty
	}

	for _, r := range key {
		if !isIdentifierChar(r) {
			return ErrMetadataKeyInvalidChar
		}
	}

	return nil
}

// ValidateMetadataString checks that a metadata string payload is safe for
// null-terminated key encodings used by the metadata read index. Values are
// free-form UTF-8 (multilingual text, descriptions, etc.) — only embedded
// null bytes are rejected.
func ValidateMetadataString(value string) error {
	if strings.ContainsRune(value, 0) {
		return ErrMetadataValueContainsNullByte
	}

	return nil
}
