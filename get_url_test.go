package main

import "testing"

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody	  string
		expected      []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one/two">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one/two", "https://other.com/path"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://other.com",
			inputBody: `
		<html>
			<body>
				<a href="https://blog.boot.dev/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			for i, _ := range tc.expected {
				if actual[i] != tc.expected[i] {
					t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected[i], actual[i])
				}
			}
		})
	}

}