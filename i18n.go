package main

import (
	"flag"
	"strings"

	"github.com/jeandeaual/go-locale"
)

type Locale string

const (
	LangRU Locale = "ru"
	LangEN Locale = "en"
)

var currentLocale Locale = LangEN

var translations = map[Locale]map[string]string{
	LangRU: {
		"History":                  "История",
		"EmptyHistory":             "История пуста",
		"Clear":                    "Очистить",
		"Paste":                    "Вставить",
		"Copied":                   "Скопировано",
		"MemoryCleared":            "Память очищена",
		"AddedToMemory":            "Добавлено в память",
		"SubtractedFromMemory":     "Вычтено из памяти",
		"CannotDivideByZero":       "Нельзя делить на ноль",
		"Error":                    "Ошибка",
		"WindowTitle":              "Material Design Calculator",
	},
	LangEN: {
		"History":                  "History",
		"EmptyHistory":             "History is empty",
		"Clear":                    "Clear",
		"Paste":                    "Paste",
		"Copied":                   "Copied",
		"MemoryCleared":            "Memory cleared",
		"AddedToMemory":            "Added to memory",
		"SubtractedFromMemory":     "Subtracted from memory",
		"CannotDivideByZero":       "Cannot divide by zero",
		"Error":                    "Error",
		"WindowTitle":              "Material Design Calculator",
	},
}

func initLocale() {
	localeFlag := flag.String("locale", "", "Force locale: ru or en")
	flag.Parse()

	if *localeFlag == "ru" {
		currentLocale = LangRU
		return
	}
	if *localeFlag == "en" {
		currentLocale = LangEN
		return
	}

	userLocale, err := locale.GetLocale()
	if err != nil {
		currentLocale = LangEN
		return
	}

	lang := strings.ToLower(userLocale)
	if strings.HasPrefix(lang, "ru") || strings.HasPrefix(lang, "be") || strings.HasPrefix(lang, "uk") {
		currentLocale = LangRU
	} else {
		currentLocale = LangEN
	}
}

func T(key string) string {
	if trans, ok := translations[currentLocale][key]; ok {
		return trans
	}
	return translations[LangEN][key]
}