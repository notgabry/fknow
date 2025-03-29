package utils

import (
	"strings"

	"github.com/eduardolat/goeasyi18n"
)

type Translator interface {
	GT(key string) string
}

type I18nWithEscape struct {
	*goeasyi18n.I18n
}

func NewI18nWithEscape(i18nInstance *goeasyi18n.I18n) *I18nWithEscape {
	return &I18nWithEscape{
		I18n: i18nInstance,
	}
}

// Get Escaped
func (i *I18nWithEscape) GT(key string) string {
	return Escape(i.I18n.Translate("it", key), []rune{'.', '!', '=', '-'})
}

// dont use tgbotapi#Escapetext or the formatting will be escaped too
func Escape(text string, charsToEscape []rune) string {
	escapeMap := make(map[rune]bool)
	for _, char := range charsToEscape {
		escapeMap[char] = true
	}

	var builder strings.Builder
	for _, char := range text {
		if escapeMap[char] {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}

	return builder.String()
}
