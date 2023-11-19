package i18n

import "golang.org/x/text/language"

var supportedLanguages = []language.Tag{language.English, language.Spanish}

// ParseLanguage parses the language string and returns a language.Tag
// If the language is invalid or not supported the default language is returned (Spanish).
func ParseLanguage(lang string) language.Tag {
	l, err := language.Parse(lang)
	if err != nil {
		return language.Spanish
	}

	for _, supportedLang := range supportedLanguages {
		if l == supportedLang {
			return l
		}
	}

	return language.Spanish
}
