package utils

import (
	"strings"

	"github.com/eduardolat/goeasyi18n"
)

type Translator interface {
	GT(key string) string
	TOnly(key string) string
}

type I18nWithEscape struct {
	*goeasyi18n.I18n
}

func NewI18nWithEscape(i18nInstance *goeasyi18n.I18n) *I18nWithEscape {
	return &I18nWithEscape{I18n: i18nInstance}
}

func (i *I18nWithEscape) GT(key string) string {
	return Escape(i.I18n.Translate("it", key), []rune{'.', '!', '=', '-'})
}

func (i *I18nWithEscape) TOnly(key string) string {
	return i.I18n.Translate("it", key)
}

// Escape escapes specific chars for MarkdownV2 without touching formatting.
func Escape(text string, charsToEscape []rune) string {
	escapeMap := make(map[rune]bool)
	for _, char := range charsToEscape {
		escapeMap[char] = true
	}
	var b strings.Builder
	for _, char := range text {
		if escapeMap[char] {
			b.WriteRune('\\')
		}
		b.WriteRune(char)
	}
	return b.String()
}
