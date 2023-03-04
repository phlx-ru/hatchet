package sanitize

import (
	"regexp"
	"strings"
)

const (
	countryCode      = `7`
	nonDigitsPattern = `\D+`
)

// Phone from `+7 900 900 90 90` returns `9009009090`
func Phone(phone string) string {
	if phone == `` {
		return ``
	}
	sanitized := regexp.MustCompile(nonDigitsPattern).ReplaceAllString(phone, ``)
	if strings.HasPrefix(sanitized, `8`) || strings.HasPrefix(sanitized, `7`) {
		sanitized = sanitized[1:]
	}
	return sanitized
}

// PhoneWithCountryCode from `900 900 90 90` returns `79009009090`
func PhoneWithCountryCode(phone string) string {
	if phone == `` {
		return ``
	}
	sanitized := Phone(phone)
	return countryCode + sanitized
}
