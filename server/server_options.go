package server

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
