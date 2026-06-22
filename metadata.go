package domain

import "strings"

// ValidateMetadataKey checks that a metadata key is safe for null-terminated
// key encodings. Null bytes would corrupt the canonical key layouts used by
// the read index.
func ValidateMetadataKey(key string) error {
	if key == "" {
		return ErrMetadataKeyEmpty
	}

	if strings.ContainsRune(key, 0) {
		return ErrMetadataKeyContainsNullByte
	}

	return nil
}

// ValidateMetadataString checks that a metadata string payload is safe for
// null-terminated key encodings used by the metadata read index. Empty
// strings are accepted — only embedded null bytes are rejected.
func ValidateMetadataString(value string) error {
	if strings.ContainsRune(value, 0) {
		return ErrMetadataValueContainsNullByte
	}

	return nil
}
