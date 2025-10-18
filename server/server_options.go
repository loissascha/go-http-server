package server

import "os"

type ServerOptionName string

const (
	TRANSLATIONS_ENABLED ServerOptionName = "translations_enabled"
	TRANSLATIONS_ADD     ServerOptionName = "translations_add"
	TRANSLATION_DEFAULT  ServerOptionName = "translation_default"
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

func (s *Server) initServerOptions() {
	for _, option := range s.Options {
		switch option.Name {
		case TRANSLATIONS_ENABLED:
			s.TranslationsEnabled = true
		case TRANSLATIONS_ADD:
			readTranslationFile(option.Filename)
			break
		case TRANSLATION_DEFAULT:
			s.DefaultLanguage = option.Value
		}
	}
}

func readTranslationFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}
