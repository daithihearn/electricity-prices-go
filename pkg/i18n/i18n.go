package i18n

import (
	"embed"
	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
)

// Embedding the TOML files.
// The path should be relative to this Go file.
// Adjust the path if your files are located elsewhere.
//
//go:embed *.toml
var translationsFS embed.FS

type Translations map[string]string

func loadTranslations(lang language.Tag, filename string) error {
	var translations Translations

	// Reading the file from the embedded filesystem
	data, err := translationsFS.ReadFile(filename)
	if err != nil {
		return err
	}

	// Decoding the TOML data
	if _, err := toml.Decode(string(data), &translations); err != nil {
		return err
	}

	for key, value := range translations {
		err := message.SetString(lang, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitialiseTranslations() error {
	log.Println("Loading translations")
	err := loadTranslations(language.English, "en.toml")
	if err != nil {
		return err
	}
	err = loadTranslations(language.Spanish, "es.toml")
	if err != nil {
		return err
	}

	return nil
}
