package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type ServerOptionName string

const (
	TRANSLATIONS_ENABLED          ServerOptionName = "translations_enabled"
	TRANSLATIONS_ADD              ServerOptionName = "translations_add"
	TRANSLATION_DEFAULT           ServerOptionName = "translation_default"
	TRANSLATIONS_AUTO_DETECT_LANG ServerOptionName = "translations_auto_detect_language"
	EXPORT_TYPE                   ServerOptionName = "exprt_type"
	EXPORT_TYPE_LOCATION          ServerOptionName = "exprt_type_location"
	READ_HEADER_TIMEOUT           ServerOptionName = "read_header_timeout"
	READ_TIMEOUT                  ServerOptionName = "read_timeout"
	WRITE_TIMEOUT                 ServerOptionName = "write_timeout"
	IDLE_TIMEOUT                  ServerOptionName = "idle_timeout"
	MAX_HEADER_BYTES              ServerOptionName = "max_header_bytes"
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

func SetExportTypesLocation(path string) ServerOption {
	return ServerOption{
		Name:  EXPORT_TYPE_LOCATION,
		Value: path,
	}
}

func EnableExportTypes(enable bool) ServerOption {
	v := "disable"
	if enable {
		v = "enable"
	}
	return ServerOption{
		Name:  EXPORT_TYPE,
		Value: v,
	}
}

func SetReadHeaderTimeout(d time.Duration) ServerOption {
	return ServerOption{
		Name:  READ_HEADER_TIMEOUT,
		Value: d.String(),
	}
}

func SetReadTimeout(d time.Duration) ServerOption {
	return ServerOption{
		Name:  READ_TIMEOUT,
		Value: d.String(),
	}
}

func SetWriteTimeout(d time.Duration) ServerOption {
	return ServerOption{
		Name:  WRITE_TIMEOUT,
		Value: d.String(),
	}
}

func SetIdleTimeout(d time.Duration) ServerOption {
	return ServerOption{
		Name:  IDLE_TIMEOUT,
		Value: d.String(),
	}
}

func SetMaxHeaderBytes(n int) ServerOption {
	return ServerOption{
		Name:  MAX_HEADER_BYTES,
		Value: strconv.Itoa(n),
	}
}

func (s *Server) initServerOptions() error {
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
		case EXPORT_TYPE:
			if option.Value == "enable" {
				s.ExportTypes = true
			}
			if option.Value == "disable" {
				s.ExportTypes = false
			}
		case EXPORT_TYPE_LOCATION:
			s.ExportTypesLocation = option.Value
		case READ_HEADER_TIMEOUT:
			d, err := time.ParseDuration(option.Value)
			if err != nil {
				return fmt.Errorf("invalid read header timeout %q: %w", option.Value, err)
			}
			s.ReadHeaderTimeout = d
		case READ_TIMEOUT:
			d, err := time.ParseDuration(option.Value)
			if err != nil {
				return fmt.Errorf("invalid read timeout %q: %w", option.Value, err)
			}
			s.ReadTimeout = d
		case WRITE_TIMEOUT:
			d, err := time.ParseDuration(option.Value)
			if err != nil {
				return fmt.Errorf("invalid write timeout %q: %w", option.Value, err)
			}
			s.WriteTimeout = d
		case IDLE_TIMEOUT:
			d, err := time.ParseDuration(option.Value)
			if err != nil {
				return fmt.Errorf("invalid idle timeout %q: %w", option.Value, err)
			}
			s.IdleTimeout = d
		case MAX_HEADER_BYTES:
			n, err := strconv.Atoi(option.Value)
			if err != nil {
				return fmt.Errorf("invalid max header bytes %q: %w", option.Value, err)
			}
			s.MaxHeaderBytes = n
		}
	}
	return nil
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
