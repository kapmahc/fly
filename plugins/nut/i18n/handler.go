package i18n

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const (
	// LOCALE key
	LOCALE = "locale"
)

// DetectLocale detect locale from request
func DetectLocale(c *gin.Context) {
	var written bool
	// 1. Check URL arguments.
	lang := c.Query(LOCALE)

	// 2. Get language information from cookies.
	if lang == "" {
		if ck, er := c.Cookie(LOCALE); er == nil {
			lang = ck
		}
	} else {
		written = true
	}

	if lang == "" {
		// 3. Get language information from 'Accept-Language'.
		if al := c.GetHeader("Accept-Language"); len(al) > 4 {
			lang = al[:5] // Only compare first 5 letters.
		}
	}

	tag, err := language.Parse(lang)
	if err != nil {
		tag = language.AmericanEnglish
	}
	tag, _, _ = language.NewMatcher(_languages).Match(tag)
	if lang != tag.String() {
		lang = tag.String()
		written = true
	}

	if written {
		c.SetCookie(LOCALE, lang, 1<<32-1, "", "", false, false)
	}
	c.Set(LOCALE, lang)
	c.Set("languages", _languages)
}
