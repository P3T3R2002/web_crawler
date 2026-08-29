package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev",
		},
		{
			name:     "remove scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme",
			inputURL: "http:/blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev//path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme",
			inputURL: "https://////blog.boot.dev/////path///",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path/subpath",
			expected: "blog.boot.dev/path/subpath",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := normalizeURL(tc.inputURL)
			
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}

}









