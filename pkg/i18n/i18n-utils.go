package i18n

import (
	"golang.org/x/text/language"
)

var supportedLanguages = []language.Tag{language.English, language.Spanish}

// ParseLanguage parses the language string and returns a language.Tag
// If the language is invalid or not supported the default language is returned (Spanish).
func ParseLanguage(lang string) language.Tag {
	l, err := language.Parse(lang)
	if err != nil {
		return language.Spanish
	}

	for _, supportedLang := range supportedLanguages {
		// If the first two letters match a supported language, return it
		if l.String()[:2] == supportedLang.String()[:2] {
			return supportedLang
		}
	}

	return language.Spanish
}
