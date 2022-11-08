package l10n

import "fmt"

type Locale struct {
	code string
	vals map[string]string
}

var localeEn = &Locale{"en", localizationEn}
var localeRu = &Locale{"ru", localizationRu}

func TryLangs(langs []string) *Locale {
	for _, lang := range langs {
		switch lang {
		case "ru":
			return localeRu

		case "en":
			return localeEn
		}
	}

	return localeEn
}

func (l *Locale) Tr(key string, args ...any) string {
	if raw, ok := l.vals[key]; ok {
		return fmt.Sprintf(raw, args...)
	}
	return key
}

func (l *Locale) Code() string {
	return l.code
}
