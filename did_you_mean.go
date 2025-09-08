package dymean

import (
	"sort"
)

// Suggestion represents a word suggestion with its similarity score
type Suggestion struct {
	Word       string
	Similarity float64
}

// DidYouMean is the main struct for the spell checker
type DidYouMean struct {
	bloomFilters map[Language]*BloomFilter // One Bloom filter per language
	candidates   *CandidateGenerator
	dictionaries map[Language]map[string]bool // One dictionary per language
	currentLang  Language
}

// NewDidYouMean creates a new DidYouMean instance
func NewDidYouMean(dictionarySize uint, numHashFuncs int) *DidYouMean {
	return &DidYouMean{
		bloomFilters: make(map[Language]*BloomFilter),
		candidates:   NewCandidateGenerator(),
		dictionaries: make(map[Language]map[string]bool),
		currentLang:  English, // Default to English
	}
}

// AddWords adds words to the dictionary for the current language
func (dym *DidYouMean) AddWords(words []string) {
	dym.AddWordsForLanguage(words, dym.currentLang)
}

// AddWordsForLanguage adds words to the dictionary for a specific language
func (dym *DidYouMean) AddWordsForLanguage(words []string, lang Language) {
	// Initialize Bloom filter and dictionary for this language if not exists
	if dym.bloomFilters[lang] == nil {
		dym.bloomFilters[lang] = NewBloomFilter(10000, 7)
		dym.dictionaries[lang] = make(map[string]bool)
	}

	langInfo := GetLanguageInfo(lang)

	for _, word := range words {
		normalized := langInfo.Normalizer(word)
		if IsValidWordForLanguage(normalized, lang) {
			dym.bloomFilters[lang].Add(normalized)
			dym.dictionaries[lang][normalized] = true
		}
	}
}

// SetLanguage sets the current language
func (dym *DidYouMean) SetLanguage(lang Language) {
	dym.currentLang = lang
}

// GetCurrentLanguage returns the current language
func (dym *DidYouMean) GetCurrentLanguage() Language {
	return dym.currentLang
}

// LoadDefaultDictionary loads the default dictionary for a language
func (dym *DidYouMean) LoadDefaultDictionary(lang Language) {
	words := GetWordsForLanguage(lang)
	dym.AddWordsForLanguage(words, lang)
}

// IsCorrect checks if a word is in the dictionary for the current language
func (dym *DidYouMean) IsCorrect(word string) bool {
	return dym.IsCorrectForLanguage(word, dym.currentLang)
}

// IsCorrectForLanguage checks if a word is in the dictionary for a specific language
func (dym *DidYouMean) IsCorrectForLanguage(word string, lang Language) bool {
	if dym.bloomFilters[lang] == nil || dym.dictionaries[lang] == nil {
		return false
	}

	langInfo := GetLanguageInfo(lang)
	normalized := langInfo.Normalizer(word)

	return dym.bloomFilters[lang].Contains(normalized) && dym.dictionaries[lang][normalized]
}

// GetSuggestions returns suggestions for a misspelled word in the current language
func (dym *DidYouMean) GetSuggestions(word string, maxSuggestions int, maxEditDistance int) []Suggestion {
	return dym.GetSuggestionsForLanguage(word, maxSuggestions, maxEditDistance, dym.currentLang)
}

// GetSuggestionsForLanguage returns suggestions for a misspelled word in a specific language
func (dym *DidYouMean) GetSuggestionsForLanguage(word string, maxSuggestions int, maxEditDistance int, lang Language) []Suggestion {
	if dym.bloomFilters[lang] == nil || dym.dictionaries[lang] == nil {
		return nil
	}

	langInfo := GetLanguageInfo(lang)
	normalized := langInfo.Normalizer(word)

	if !IsValidWordForLanguage(normalized, lang) {
		return nil
	}

	// If the word is correct, return it
	if dym.IsCorrectForLanguage(normalized, lang) {
		return []Suggestion{{Word: normalized, Similarity: 1.0}}
	}

	// Generate candidates
	candidates := dym.candidates.GenerateCandidates(normalized, maxEditDistance)

	// Also include common typo candidates
	typoCandidates := dym.candidates.GenerateCommonTypos(normalized)
	candidates = append(candidates, typoCandidates...)

	// Filter candidates that exist in the dictionary
	validCandidates := make([]string, 0)
	for _, candidate := range candidates {
		if dym.bloomFilters[lang].Contains(candidate) && dym.dictionaries[lang][candidate] {
			validCandidates = append(validCandidates, candidate)
		}
	}

	// Calculate similarity scores and create suggestions
	suggestions := make([]Suggestion, 0, len(validCandidates))
	for _, candidate := range validCandidates {
		similarity := CalculateSimilarity(normalized, candidate)
		suggestions = append(suggestions, Suggestion{
			Word:       candidate,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Similarity > suggestions[j].Similarity
	})

	// Return top suggestions
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return suggestions
}

// Suggest returns the best suggestion for a word in the current language
func (dym *DidYouMean) Suggest(word string) string {
	return dym.SuggestForLanguage(word, dym.currentLang)
}

// SuggestForLanguage returns the best suggestion for a word in a specific language
func (dym *DidYouMean) SuggestForLanguage(word string, lang Language) string {
	suggestions := dym.GetSuggestionsForLanguage(word, 1, 2, lang)
	if len(suggestions) > 0 {
		return suggestions[0].Word
	}
	return word // Return original if no suggestions found
}

// GetSuggestionsWithThreshold returns suggestions above a similarity threshold
func (dym *DidYouMean) GetSuggestionsWithThreshold(word string, threshold float64, maxSuggestions int) []Suggestion {
	allSuggestions := dym.GetSuggestions(word, maxSuggestions*2, 2) // Get more to filter
	filtered := make([]Suggestion, 0)

	for _, suggestion := range allSuggestions {
		if suggestion.Similarity >= threshold {
			filtered = append(filtered, suggestion)
		}
	}

	if len(filtered) > maxSuggestions {
		filtered = filtered[:maxSuggestions]
	}

	return filtered
}

// CheckAndSuggest checks a word and returns suggestions if it's misspelled
func (dym *DidYouMean) CheckAndSuggest(word string) (bool, []Suggestion) {
	return dym.CheckAndSuggestForLanguage(word, dym.currentLang)
}

// CheckAndSuggestForLanguage checks a word and returns suggestions if it's misspelled for a specific language
func (dym *DidYouMean) CheckAndSuggestForLanguage(word string, lang Language) (bool, []Suggestion) {
	if dym.IsCorrectForLanguage(word, lang) {
		return true, nil
	}

	suggestions := dym.GetSuggestionsForLanguage(word, 5, 2, lang)
	return false, suggestions
}

// AutoDetectAndSuggest automatically detects language and provides suggestions
func (dym *DidYouMean) AutoDetectAndSuggest(word string) (Language, bool, []Suggestion) {
	detectedLang := DetectLanguage(word)
	isCorrect, suggestions := dym.CheckAndSuggestForLanguage(word, detectedLang)
	return detectedLang, isCorrect, suggestions
}
