package dymean

import (
	"strings"
	"unicode"
)

// Language represents supported languages
type Language string

const (
	English  Language = "en"
	Persian  Language = "fa"
	Arabic   Language = "ar"
	French   Language = "fr"
	Spanish  Language = "es"
	German   Language = "de"
	Italian  Language = "it"
	Russian  Language = "ru"
	Chinese  Language = "zh"
	Japanese Language = "ja"
	Korean   Language = "ko"
)

// LanguageInfo contains information about a language
type LanguageInfo struct {
	Code       Language
	Name       string
	Direction  string // "ltr" or "rtl"
	Alphabet   string
	IsRTL      bool
	Normalizer func(string) string
}

// GetLanguageInfo returns information about a language
func GetLanguageInfo(lang Language) LanguageInfo {
	switch lang {
	case English:
		return LanguageInfo{
			Code:       English,
			Name:       "English",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyz",
			IsRTL:      false,
			Normalizer: normalizeEnglish,
		}
	case Persian:
		return LanguageInfo{
			Code:       Persian,
			Name:       "Persian",
			Direction:  "rtl",
			Alphabet:   "ابپتثجچحخدذرزژسشصضطظعغفقکگلمنوهی",
			IsRTL:      true,
			Normalizer: normalizePersian,
		}
	case Arabic:
		return LanguageInfo{
			Code:       Arabic,
			Name:       "Arabic",
			Direction:  "rtl",
			Alphabet:   "ابتثجحخدذرزسشصضطظعغفقكلمنهوي",
			IsRTL:      true,
			Normalizer: normalizeArabic,
		}
	case French:
		return LanguageInfo{
			Code:       French,
			Name:       "French",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyzàâäéèêëïîôöùûüÿç",
			IsRTL:      false,
			Normalizer: normalizeFrench,
		}
	case Spanish:
		return LanguageInfo{
			Code:       Spanish,
			Name:       "Spanish",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyzñáéíóúü",
			IsRTL:      false,
			Normalizer: normalizeSpanish,
		}
	case German:
		return LanguageInfo{
			Code:       German,
			Name:       "German",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyzäöüß",
			IsRTL:      false,
			Normalizer: normalizeGerman,
		}
	case Italian:
		return LanguageInfo{
			Code:       Italian,
			Name:       "Italian",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyzàèéìíîòóùú",
			IsRTL:      false,
			Normalizer: normalizeItalian,
		}
	case Russian:
		return LanguageInfo{
			Code:       Russian,
			Name:       "Russian",
			Direction:  "ltr",
			Alphabet:   "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
			IsRTL:      false,
			Normalizer: normalizeRussian,
		}
	case Chinese:
		return LanguageInfo{
			Code:       Chinese,
			Name:       "Chinese",
			Direction:  "ltr",
			Alphabet:   "", // Chinese doesn't use alphabet
			IsRTL:      false,
			Normalizer: normalizeChinese,
		}
	case Japanese:
		return LanguageInfo{
			Code:       Japanese,
			Name:       "Japanese",
			Direction:  "ltr",
			Alphabet:   "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん",
			IsRTL:      false,
			Normalizer: normalizeJapanese,
		}
	case Korean:
		return LanguageInfo{
			Code:       Korean,
			Name:       "Korean",
			Direction:  "ltr",
			Alphabet:   "ㄱㄴㄷㄹㅁㅂㅅㅇㅈㅊㅋㅌㅍㅎㅏㅑㅓㅕㅗㅛㅜㅠㅡㅣ",
			IsRTL:      false,
			Normalizer: normalizeKorean,
		}
	default:
		return LanguageInfo{
			Code:       English,
			Name:       "English",
			Direction:  "ltr",
			Alphabet:   "abcdefghijklmnopqrstuvwxyz",
			IsRTL:      false,
			Normalizer: normalizeEnglish,
		}
	}
}

// Normalization functions for different languages
func normalizeEnglish(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizePersian(word string) string {
	// Remove diacritics and normalize Persian text
	word = strings.TrimSpace(word)
	// Convert Arabic numerals to Persian numerals if needed
	word = strings.ReplaceAll(word, "0", "۰")
	word = strings.ReplaceAll(word, "1", "۱")
	word = strings.ReplaceAll(word, "2", "۲")
	word = strings.ReplaceAll(word, "3", "۳")
	word = strings.ReplaceAll(word, "4", "۴")
	word = strings.ReplaceAll(word, "5", "۵")
	word = strings.ReplaceAll(word, "6", "۶")
	word = strings.ReplaceAll(word, "7", "۷")
	word = strings.ReplaceAll(word, "8", "۸")
	word = strings.ReplaceAll(word, "9", "۹")
	return word
}

func normalizeArabic(word string) string {
	return strings.TrimSpace(word)
}

func normalizeFrench(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizeSpanish(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizeGerman(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizeItalian(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizeRussian(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func normalizeChinese(word string) string {
	return strings.TrimSpace(word)
}

func normalizeJapanese(word string) string {
	return strings.TrimSpace(word)
}

func normalizeKorean(word string) string {
	return strings.TrimSpace(word)
}

// DetectLanguage attempts to detect the language of a word
func DetectLanguage(word string) Language {
	if len(word) == 0 {
		return English
	}

	// Check for Persian/Arabic characters
	for _, r := range word {
		if r >= 0x0600 && r <= 0x06FF { // Arabic block
			return Persian // Default to Persian for Arabic script
		}
		if r >= 0x0750 && r <= 0x077F { // Arabic Supplement
			return Persian
		}
		if r >= 0x08A0 && r <= 0x08FF { // Arabic Extended-A
			return Persian
		}
		if r >= 0xFB50 && r <= 0xFDFF { // Arabic Presentation Forms-A
			return Persian
		}
		if r >= 0xFE70 && r <= 0xFEFF { // Arabic Presentation Forms-B
			return Persian
		}
	}

	// Check for Cyrillic characters (Russian)
	for _, r := range word {
		if r >= 0x0400 && r <= 0x04FF { // Cyrillic block
			return Russian
		}
	}

	// Check for Chinese characters
	for _, r := range word {
		if r >= 0x4E00 && r <= 0x9FFF { // CJK Unified Ideographs
			return Chinese
		}
	}

	// Check for Japanese Hiragana
	for _, r := range word {
		if r >= 0x3040 && r <= 0x309F { // Hiragana
			return Japanese
		}
	}

	// Check for Korean Hangul
	for _, r := range word {
		if r >= 0xAC00 && r <= 0xD7AF { // Hangul Syllables
			return Korean
		}
	}

	// Default to English for Latin script
	return English
}

// IsValidWordForLanguage checks if a word contains only valid characters for a language
func IsValidWordForLanguage(word string, lang Language) bool {
	if len(word) == 0 {
		return false
	}

	langInfo := GetLanguageInfo(lang)

	// For languages without alphabet (like Chinese), check for valid Unicode ranges
	if langInfo.Alphabet == "" {
		switch lang {
		case Chinese:
			for _, r := range word {
				if !unicode.Is(unicode.Han, r) && !unicode.IsLetter(r) {
					return false
				}
			}
			return true
		case Japanese:
			for _, r := range word {
				if !unicode.Is(unicode.Hiragana, r) && !unicode.Is(unicode.Katakana, r) &&
					!unicode.Is(unicode.Han, r) && !unicode.IsLetter(r) {
					return false
				}
			}
			return true
		case Korean:
			for _, r := range word {
				if !unicode.Is(unicode.Hangul, r) && !unicode.IsLetter(r) {
					return false
				}
			}
			return true
		}
	}

	// For languages with alphabet, check if all characters are in the alphabet
	for _, r := range word {
		if !strings.ContainsRune(langInfo.Alphabet, r) && !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// GetSupportedLanguages returns a list of all supported languages
func GetSupportedLanguages() []Language {
	return []Language{
		English, Persian, Arabic, French, Spanish, German,
		Italian, Russian, Chinese, Japanese, Korean,
	}
}
