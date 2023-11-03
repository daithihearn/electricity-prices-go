package i18n

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
)

type Translations map[string]string

func loadTranslations(lang language.Tag, filename string) {
	var translations Translations
	if _, err := toml.DecodeFile(filename, &translations); err != nil {
		log.Fatalf("Failed to load translations for %s: %v", lang, err)
	}
	for key, value := range translations {
		message.SetString(lang, key, value)
	}
}

func InitialiseTranslations() {
	log.Println("Loading translations")
	loadTranslations(language.English, "pkg/i18n/en.toml")
	loadTranslations(language.Spanish, "pkg/i18n/es.toml")
}
