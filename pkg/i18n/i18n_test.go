package i18n

import (
	"golang.org/x/text/language"
	"testing"
)

func TestInitialiseTranslations(t *testing.T) {
	testCases := []struct {
		name      string
		files     []File
		expectErr bool
	}{
		{
			name: "Valid",
			files: []File{
				{
					Filename: "en.toml",
					Lang:     language.English,
				},
				{
					Filename: "es.toml",
					Lang:     language.Spanish,
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid",
			files: []File{
				{
					Filename: "invalid.toml",
					Lang:     language.English,
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := InitialiseTranslations(tc.files)
			if tc.expectErr && actual == nil {
				t.Error("expected error, got nil")
			}
			if !tc.expectErr && actual != nil {
				t.Errorf("expected nil, got %s", actual)
			}
		})
	}
}
