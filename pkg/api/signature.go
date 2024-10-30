package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// Sign generates an HMAC SHA512 signature from the given secret and payload.
func Sign(secret string, payload params) (string, error) {
	sig, err := signPayload([]byte(secret), payload)
	if err != nil {
		return "", fmt.Errorf("failed to sign payload: %w", err)
	}
	// fmt.Printf("payload:[%+v]\n", payload)
	// fmt.Printf("signature:[%v]\n", hex.EncodeToString(sig))
	return hex.EncodeToString(sig), nil
}

// signPayload computes the HMAC SHA512 signature for the given parameters.
func signPayload(secret []byte, p params) ([]byte, error) {
	qs := queryString(p)
	mac := hmac.New(sha512.New, secret)
	_, err := mac.Write([]byte(qs))
	if err != nil {
		return nil, fmt.Errorf("failed to write query string to HMAC: %w", err)
	}
	return mac.Sum(nil), nil
}

// Verify checks if the provided signature is valid for the given parameters.
func Verify(secret []byte, p params, sig []byte) bool {
	// Capture both return values from signPayload
	calculatedSig, err := signPayload(secret, p)
	if err != nil {
		return false // Return false if there is an error during signing
	}
	return hmac.Equal(calculatedSig, sig)
}

// queryString creates a URL-encoded query string from the given parameters.
func queryString(p params) string {
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	// or sort.Strings(keys)
	slices.Sort(keys)

	var builder strings.Builder
	for i, k := range keys {
		v := p[k]

		appendKeyValue(&builder, k, v)
		if i != len(keys)-1 {
			builder.WriteString("&")
		}
	}
	return builder.String()
}

// appendKeyValue appends a key-value pair to the builder, handling nested maps and arrays.
func appendKeyValue(builder *strings.Builder, key string, value any) {
	switch v := value.(type) {
	case map[string]any:
		for nestedKey, nestedValue := range v {
			appendKeyValue(builder, fmt.Sprintf("%s[%s]", key, nestedKey), nestedValue)
		}
	case []any:
		for i, item := range v {
			appendKeyValue(builder, fmt.Sprintf("%s[%d]", key, i), item)
		}
	default:
		writeKeyValue(builder, key, value)
	}
}

// writeKeyValue writes a key-value pair to the builder.
func writeKeyValue(builder *strings.Builder, key string, value any) {
	switch v := value.(type) {
	case float64:
		fmt.Fprintf(builder, "%s=%.f", key, v)
	default:
		fmt.Fprintf(builder, "%s=%v", key, v)
	}
}

func convertURLValuesToParams(values url.Values) params {
	result := make(params)

	for key, vals := range values {
		switch len(vals) {
		case 0:
			// Key exists but has no values; decide on your handling
			result[key] = nil // or result[key] = ""
		case 1:
			// Single value, use it directly
			result[key] = vals[0]
		default:
			// Multiple values, store as slice
			result[key] = vals
		}
	}

	return result
}
