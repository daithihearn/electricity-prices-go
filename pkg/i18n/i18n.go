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

type File struct {
	Filename string
	Lang     language.Tag
}

func loadTranslations(file File) error {
	var translations Translations

	// Reading the file from the embedded filesystem
	data, err := translationsFS.ReadFile(file.Filename)
	if err != nil {
		return err
	}

	// Decoding the TOML data
	if _, err := toml.Decode(string(data), &translations); err != nil {
		return err
	}

	for key, value := range translations {
		err := message.SetString(file.Lang, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitialiseTranslations(files []File) error {
	log.Println("Loading translations")

	for _, file := range files {
		err := loadTranslations(file)
		if err != nil {
			return err
		}
	}

	return nil
}
