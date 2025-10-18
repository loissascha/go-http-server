package server

import (
	"encoding/json"
	"io"
	"os"
)

type ServerOptionName string

const (
	TRANSLATIONS_ENABLED          ServerOptionName = "translations_enabled"
	TRANSLATIONS_ADD              ServerOptionName = "translations_add"
	TRANSLATION_DEFAULT           ServerOptionName = "translation_default"
	TRANSLATIONS_AUTO_DETECT_LANG ServerOptionName = "translations_auto_detect_language"
)

type ServerOption struct {
	Name     ServerOptionName
	Value    string
	Filename string
}

func EnableTranslations() ServerOption {
	return ServerOption{
		Name: TRANSLATIONS_ENABLED,
	}
}

func AddTranslationFile(short, filepath string) ServerOption {
	return ServerOption{
		Name:     TRANSLATIONS_ADD,
		Value:    short,
		Filename: filepath,
	}
}

func SetDefaultLanguage(short string) ServerOption {
	return ServerOption{
		Name:  TRANSLATION_DEFAULT,
		Value: short,
	}
}

func EnableAutoDetectLanguage() ServerOption {
	return ServerOption{
		Name: TRANSLATIONS_AUTO_DETECT_LANG,
	}
}

func (s *Server) initServerOptions() {
	for _, option := range s.Options {
		switch option.Name {
		case TRANSLATIONS_ENABLED:
			s.TranslationsEnabled = true
		case TRANSLATIONS_ADD:
			data := readTranslationFile(option.Filename)
			s.Languages[option.Value] = data
		case TRANSLATION_DEFAULT:
			s.DefaultLanguage = option.Value
		case TRANSLATIONS_AUTO_DETECT_LANG:
			s.AutoDetectLanguageEnabled = true
		}
	}
}

func readTranslationFile(filepath string) map[string]string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var data map[string]string
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	return data
}
