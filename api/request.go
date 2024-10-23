package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type secType int

const (
	secTypeNone   secType = iota // authentication is not required
	secTypeSigned                // authentication is required
)

type params map[string]any

type request struct {
	method   string
	endpoint string
	fullURL  string
	secType  secType
	query    url.Values
	form     url.Values
	header   http.Header
	body     io.Reader
}

// addParam add param with key/value to query string
// func (r *request) addQueryParam(key string, value any) *request {
// 	if r.query == nil {
// 		r.query = make(url.Values)
// 	}
// 	r.query.Add(key, fmt.Sprintf("%v", value))
// 	return r
// }

// setParam set param with key/value to query string
func (r *request) setQueryParam(key string, value any) error {
	if r.query == nil {
		r.query = make(url.Values)
	}

	switch v := value.(type) {
	case []string:
		// Add each element of the string slice as a separate key-value pair
		for _, item := range v {
			r.query.Add(key, item)
		}
	case []int:
		// Add each element of the int slice as a separate key-value pair
		for _, item := range v {
			r.query.Add(key, fmt.Sprintf("%d", item))
		}
	case []any:
		// Convert arbitrary slice to JSON and set as a single parameter value
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("error marshalling slice value for key '%s': %w", key, err)
		}
		r.query.Set(key, string(jsonValue))
	default:
		// Handle single values by converting them to string
		r.query.Set(key, fmt.Sprintf("%v", value))
	}

	return nil
}

// setParams set params with key/values to query string
func (r *request) setQueryParams(m params) *request {
	for k, v := range m {
		r.setQueryParam(k, v)
	}
	return r
}

// setFormParam set param with key/value to request form body
// func (r *request) setFormParam(key string, value any) *request {
// 	if r.form == nil {
// 		r.form = make(url.Values)
// 	}
// 	r.form.Set(key, fmt.Sprintf("%v", value))
// 	return r
// }

// setFormParams set params with key/values to request form body
// func (r *request) setFormParams(m params) *request {
// 	for k, v := range m {
// 		r.setFormParam(k, v)
// 	}
// 	return r
// }

// ensureInitialized sets up default values for query and form parameters
// if they are nil. This ensures that the request object is in a valid state
// for subsequent operations.
func (r *request) ensureInitialized() {
	// Initialize query and form to empty Values if they are nil
	if r.query == nil {
		r.query = make(url.Values)
	}

	if r.form == nil {
		r.form = make(url.Values)
	}
}

// RequestOption define option type for request
type RequestOption func(*request)

// WithHeader set or add a header value to the request
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *request) {
		if r.header == nil {
			r.header = make(http.Header)
		}
		if replace {
			r.header.Set(key, value)
		} else {
			r.header.Add(key, value)
		}
	}
}

// WithHeaders set or replace the headers of the request
func WithHeaders(header http.Header) RequestOption {
	return func(r *request) {
		r.header = header.Clone()
	}
}
