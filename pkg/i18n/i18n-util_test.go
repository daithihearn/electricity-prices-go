package i18n

import "testing"

func TestParseLanguage(t *testing.T) {
	testCases := []struct {
		name     string
		lang     string
		expected string
	}{
		{
			name:     "English",
			lang:     "en",
			expected: "en",
		},
		{
			name:     "English (US)",
			lang:     "en-US",
			expected: "en",
		},
		{
			name:     "English (GB)",
			lang:     "en-GB",
			expected: "en",
		},
		{
			name:     "Spanish",
			lang:     "es",
			expected: "es",
		},
		{
			name:     "Spanish (MX)",
			lang:     "es-MX",
			expected: "es",
		},
		{
			name:     "French",
			lang:     "fr",
			expected: "es",
		},
		{
			name:     "German",
			lang:     "de",
			expected: "es",
		},
		{
			name:     "Invalid",
			lang:     "invalid",
			expected: "es",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ParseLanguage(tc.lang)
			if actual.String() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			}
		})
	}
}
