# Dymean - Multi-Language "Did You Mean" Package with Bloom Filter

A comprehensive Go package that implements a multi-language "Did you mean" spell checker using Bloom filters for efficient dictionary storage and Levenshtein distance for similarity calculations.

## 🌍 Features

- **Multi-Language Support**: 11 languages including English, Persian, Arabic, French, Spanish, German, Italian, Russian, Chinese, Japanese, and Korean
- **Automatic Language Detection**: Detects language based on character sets and Unicode ranges
- **RTL Support**: Full support for Right-to-Left languages (Arabic, Persian)
- **Bloom Filter**: Space-efficient probabilistic data structure for fast dictionary lookups
- **Levenshtein Distance**: Calculates edit distance between words for similarity scoring
- **Candidate Generation**: Generates possible corrections using edit operations (insertion, deletion, substitution, transposition)
- **Keyboard Layout Awareness**: Considers common typing errors based on QWERTY keyboard layout
- **Language-Specific Normalization**: Proper text normalization for each language
- **Configurable Parameters**: Adjustable Bloom filter size, hash functions, and suggestion limits

## How It Works

1. **Language Detection**: Automatically detects the language of input text based on Unicode character ranges
2. **Dictionary Storage**: Words are stored in separate Bloom filters for each language
3. **Misspelling Detection**: The appropriate Bloom filter quickly identifies if a word is likely misspelled
4. **Candidate Generation**: For misspelled words, the system generates possible corrections using:
   - Edit distance algorithms (insertions, deletions, substitutions, transpositions)
   - Common keyboard layout errors
   - Language-specific character sets
5. **Similarity Scoring**: Candidates are ranked using Levenshtein distance similarity
6. **Suggestion Ranking**: Results are sorted by similarity score and returned

## Installation

```bash
go get github.com/bi0dread/dymean
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/bi0dread/dymean"
)

func main() {
    // Create a new DidYouMean instance
    dym := dymean.NewDidYouMean(10000, 7) // dictionary size, hash functions
    
    // Load default dictionaries for multiple languages
    dym.LoadDefaultDictionary(dymean.English)
    dym.LoadDefaultDictionary(dymean.Persian)
    
    // Check if a word is correct (auto-detects language)
    fmt.Println("Is 'hello' correct?", dym.IsCorrect("hello")) // true
    fmt.Println("Is 'سلام' correct?", dym.IsCorrect("سلام"))   // true
    fmt.Println("Is 'helo' correct?", dym.IsCorrect("helo"))   // false
    
    // Get suggestions for a misspelled word
    suggestions := dym.GetSuggestions("helo", 3, 2)
    for _, suggestion := range suggestions {
        fmt.Printf("%s (similarity: %.2f)\n", suggestion.Word, suggestion.Similarity)
    }
    
    // Get the best suggestion
    best := dym.Suggest("progamming")
    fmt.Println("Best suggestion for 'progamming':", best) // "programming"
    
    // Check and suggest in one call
    isCorrect, suggestions := dym.CheckAndSuggest("algoritm")
    if !isCorrect && len(suggestions) > 0 {
        fmt.Printf("Did you mean: %s?\n", suggestions[0].Word) // "algorithm"
    }
}
```

### Multi-Language Example

```go
package main

import (
    "fmt"
    "github.com/bi0dread/dymean"
)

func main() {
    dym := dymean.NewDidYouMean(10000, 7)
    
    // Load multiple language dictionaries
    dym.LoadDefaultDictionary(dymean.English)
    dym.LoadDefaultDictionary(dymean.Persian)
    dym.LoadDefaultDictionary(dymean.Arabic)
    
    // Test mixed language text
    words := []string{"hello", "سلام", "مرحبا", "world", "دنیا", "عالم"}
    
    for _, word := range words {
        // Auto-detect language and get suggestions
        detectedLang, isCorrect, suggestions := dym.AutoDetectAndSuggest(word)
        
        fmt.Printf("'%s' -> Language: %s, Correct: %t", word, detectedLang, isCorrect)
        
        if !isCorrect && len(suggestions) > 0 {
            fmt.Printf(", Suggestions: %s", suggestions[0].Word)
        }
        fmt.Println()
    }
}
```

### Language-Specific Operations

```go
// Set specific language
dym.SetLanguage(dymean.Persian)

// Check word in specific language
isCorrect := dym.IsCorrectForLanguage("سلام", dymean.Persian)

// Get suggestions for specific language
suggestions := dym.GetSuggestionsForLanguage("برنام", 3, 2, dymean.Persian)

// Add words to specific language
persianWords := []string{"سلام", "دنیا", "برنامه"}
dym.AddWordsForLanguage(persianWords, dymean.Persian)
```

### Advanced Usage

```go
// Get suggestions with a similarity threshold
suggestions := dym.GetSuggestionsWithThreshold("helo", 0.5, 5)

// Check multiple words in mixed languages
words := []string{"helo", "برنام", "wrld", "دنی", "progamming"}
for _, word := range words {
    detectedLang, isCorrect, suggestions := dym.AutoDetectAndSuggest(word)
    if !isCorrect {
        fmt.Printf("'%s' (%s) -> Did you mean: %s?\n", word, detectedLang, suggestions[0].Word)
    }
}

// Language detection examples
fmt.Println(dymean.DetectLanguage("hello"))    // English
fmt.Println(dymean.DetectLanguage("سلام"))      // Persian
fmt.Println(dymean.DetectLanguage("مرحبا"))     // Arabic
fmt.Println(dymean.DetectLanguage("привет"))    // Russian
fmt.Println(dymean.DetectLanguage("你好"))       // Chinese
```

## API Reference

### Types

```go
type Suggestion struct {
    Word       string  // The suggested word
    Similarity float64 // Similarity score (0.0 to 1.0)
}

type Language string // Language code (e.g., "en", "fa", "ar")

type LanguageInfo struct {
    Code        Language
    Name        string
    Direction   string // "ltr" or "rtl"
    Alphabet    string
    IsRTL       bool
    Normalizer  func(string) string
}

type DidYouMean struct {
    // Main spell checker instance with multi-language support
}
```

### Core Functions

```go
// Create a new DidYouMean instance
func NewDidYouMean(dictionarySize uint, numHashFuncs int) *DidYouMean

// Add words to the current language dictionary
func (dym *DidYouMean) AddWords(words []string)

// Add words to a specific language dictionary
func (dym *DidYouMean) AddWordsForLanguage(words []string, lang Language)

// Load default dictionary for a language
func (dym *DidYouMean) LoadDefaultDictionary(lang Language)

// Set current language
func (dym *DidYouMean) SetLanguage(lang Language)

// Get current language
func (dym *DidYouMean) GetCurrentLanguage() Language
```

### Spell Checking Functions

```go
// Check if a word is correct (uses current language)
func (dym *DidYouMean) IsCorrect(word string) bool

// Check if a word is correct in specific language
func (dym *DidYouMean) IsCorrectForLanguage(word string, lang Language) bool

// Get suggestions for a misspelled word (uses current language)
func (dym *DidYouMean) GetSuggestions(word string, maxSuggestions int, maxEditDistance int) []Suggestion

// Get suggestions for a misspelled word in specific language
func (dym *DidYouMean) GetSuggestionsForLanguage(word string, maxSuggestions int, maxEditDistance int, lang Language) []Suggestion

// Get the best suggestion (uses current language)
func (dym *DidYouMean) Suggest(word string) string

// Get the best suggestion for specific language
func (dym *DidYouMean) SuggestForLanguage(word string, lang Language) string

// Get suggestions above a similarity threshold
func (dym *DidYouMean) GetSuggestionsWithThreshold(word string, threshold float64, maxSuggestions int) []Suggestion

// Check and suggest in one call (uses current language)
func (dym *DidYouMean) CheckAndSuggest(word string) (bool, []Suggestion)

// Check and suggest in one call for specific language
func (dym *DidYouMean) CheckAndSuggestForLanguage(word string, lang Language) (bool, []Suggestion)

// Auto-detect language and provide suggestions
func (dym *DidYouMean) AutoDetectAndSuggest(word string) (Language, bool, []Suggestion)
```

### Language Functions

```go
// Detect language of a word
func DetectLanguage(word string) Language

// Get language information
func GetLanguageInfo(lang Language) LanguageInfo

// Get all supported languages
func GetSupportedLanguages() []Language

// Check if a word is valid for a specific language
func IsValidWordForLanguage(word string, lang Language) bool
```

### Utility Functions

```go
// Calculate Levenshtein distance between two strings
func LevenshteinDistance(s1, s2 string) int

// Calculate similarity score between two strings
func CalculateSimilarity(s1, s2 string) float64

// Check if a word contains only valid characters
func IsValidWord(word string) bool
```

### Supported Languages

```go
const (
    English Language = "en"
    Persian Language = "fa"
    Arabic  Language = "ar"
    French  Language = "fr"
    Spanish Language = "es"
    German  Language = "de"
    Italian Language = "it"
    Russian Language = "ru"
    Chinese Language = "zh"
    Japanese Language = "ja"
    Korean  Language = "ko"
)
```

## Configuration

### Bloom Filter Parameters

- **Dictionary Size**: Should be larger than the number of words you plan to store
- **Hash Functions**: More hash functions reduce false positives but increase computation time
- **Recommended**: 
  - For 1,000 words: size 10,000 and 7 hash functions
  - For 1,000,000 words: size 15,000,000 and 10 hash functions
  - For maximum accuracy: size 20,000,000 and 14 hash functions

### Suggestion Parameters

- **Max Suggestions**: Maximum number of suggestions to return
- **Max Edit Distance**: Maximum edit distance for candidate generation (1-3 recommended)
- **Similarity Threshold**: Minimum similarity score for suggestions (0.0-1.0)

### Language Configuration

- **Auto-Detection**: Automatically detects language based on Unicode character ranges
- **RTL Support**: Full support for Right-to-Left languages (Arabic, Persian)
- **Mixed Language**: Handles text with multiple languages seamlessly

## Performance

The Bloom filter provides O(k) lookup time where k is the number of hash functions, making it very fast for dictionary lookups. The candidate generation and similarity calculation are the main performance bottlenecks, but the system is optimized for typical use cases.

### Performance Benchmarks

- **Dictionary Lookup**: ~5.5ms per suggestion on Apple M2 Pro
- **Memory Usage**: 1.2MB for 1M words (7 hash functions) to 2.4MB (14 hash functions)
- **False Positive Rate**: 0.07% to 0.007% depending on configuration
- **Language Detection**: O(1) constant time based on character analysis

## Limitations

1. **False Positives**: Bloom filters can have false positives (saying a word exists when it doesn't)
2. **No False Negatives**: If the Bloom filter says a word doesn't exist, it definitely doesn't
3. **Memory vs Accuracy**: Larger Bloom filters reduce false positives but use more memory
4. **Edit Distance**: Very long words with large edit distances can be slow to process
5. **Language Detection**: Based on character sets, may not be 100% accurate for mixed scripts
6. **Dictionary Size**: Default dictionaries are limited; custom dictionaries recommended for production

## Testing

Run all tests:

```bash
go test -v
```

Run Persian-specific tests:

```bash
go test -v -run TestPersian
```

Run benchmarks:

```bash
go test -bench=.
```

Run with coverage:

```bash
go test -v -cover
```

### Test Coverage

The package includes comprehensive tests for:
- ✅ Basic spell checking functionality
- ✅ Multi-language support (English, Persian, Arabic, etc.)
- ✅ Language detection and switching
- ✅ RTL language support
- ✅ Mixed language text handling
- ✅ Bloom filter operations
- ✅ Levenshtein distance calculations
- ✅ Similarity scoring
- ✅ Performance benchmarks

## Examples

### Real-World Usage

```go
// Create a spell checker for a multi-language application
dym := dymean.NewDidYouMean(15000000, 10) // Optimized for 1M words

// Load dictionaries for your target languages
dym.LoadDefaultDictionary(dymean.English)
dym.LoadDefaultDictionary(dymean.Persian)
dym.LoadDefaultDictionary(dymean.Arabic)

// Process user input with automatic language detection
userInput := "helo wrld سلام دنیا مرحبا"
words := strings.Fields(userInput)

for _, word := range words {
    detectedLang, isCorrect, suggestions := dym.AutoDetectAndSuggest(word)
    
    if !isCorrect {
        fmt.Printf("'%s' (%s) -> Did you mean: %s?\n", 
            word, detectedLang, suggestions[0].Word)
    }
}
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Adding New Languages

To add support for a new language:

1. Add the language constant to `language.go`
2. Implement the normalization function
3. Add language detection logic
4. Create a dictionary with common words
5. Add comprehensive tests
6. Update documentation
