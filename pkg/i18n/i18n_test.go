package i18n

import "testing"

func TestInitialiseTranslations(t *testing.T) {
	err := InitialiseTranslations()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
