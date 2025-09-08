# Dymean - "Did You Mean" Package with Bloom Filter

A Go package that implements a "Did you mean" spell checker using a Bloom filter for efficient dictionary storage and Levenshtein distance for similarity calculations.

## Features

- **Bloom Filter**: Space-efficient probabilistic data structure for fast dictionary lookups
- **Levenshtein Distance**: Calculates edit distance between words for similarity scoring
- **Candidate Generation**: Generates possible corrections using edit operations (insertion, deletion, substitution, transposition)
- **Keyboard Layout Awareness**: Considers common typing errors based on QWERTY keyboard layout
- **Configurable Parameters**: Adjustable Bloom filter size, hash functions, and suggestion limits

## How It Works

1. **Dictionary Storage**: Words are stored in a Bloom filter for fast membership testing
2. **Misspelling Detection**: The Bloom filter quickly identifies if a word is likely misspelled
3. **Candidate Generation**: For misspelled words, the system generates possible corrections using:
   - Edit distance algorithms (insertions, deletions, substitutions, transpositions)
   - Common keyboard layout errors
4. **Similarity Scoring**: Candidates are ranked using Levenshtein distance similarity
5. **Suggestion Ranking**: Results are sorted by similarity score and returned

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
    
    // Add words to the dictionary
    words := []string{
        "hello", "world", "golang", "programming", "computer", "science",
        "algorithm", "data", "structure", "bloom", "filter", "spell",
        "checker", "dictionary", "suggestion", "similarity", "distance",
    }
    dym.AddWords(words)
    
    // Check if a word is correct
    fmt.Println("Is 'hello' correct?", dym.IsCorrect("hello")) // true
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

### Advanced Usage

```go
// Get suggestions with a similarity threshold
suggestions := dym.GetSuggestionsWithThreshold("helo", 0.5, 5)

// Check multiple words
words := []string{"helo", "wrld", "progamming"}
for _, word := range words {
    isCorrect, suggestions := dym.CheckAndSuggest(word)
    if !isCorrect {
        fmt.Printf("'%s' -> Did you mean: %s?\n", word, suggestions[0].Word)
    }
}
```

## API Reference

### Types

```go
type Suggestion struct {
    Word       string  // The suggested word
    Similarity float64 // Similarity score (0.0 to 1.0)
}

type DidYouMean struct {
    // Main spell checker instance
}
```

### Functions

```go
// Create a new DidYouMean instance
func NewDidYouMean(dictionarySize uint, numHashFuncs int) *DidYouMean

// Add words to the dictionary
func (dym *DidYouMean) AddWords(words []string)

// Check if a word is correct
func (dym *DidYouMean) IsCorrect(word string) bool

// Get suggestions for a misspelled word
func (dym *DidYouMean) GetSuggestions(word string, maxSuggestions int, maxEditDistance int) []Suggestion

// Get the best suggestion
func (dym *DidYouMean) Suggest(word string) string

// Get suggestions above a similarity threshold
func (dym *DidYouMean) GetSuggestionsWithThreshold(word string, threshold float64, maxSuggestions int) []Suggestion

// Check and suggest in one call
func (dym *DidYouMean) CheckAndSuggest(word string) (bool, []Suggestion)
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

## Configuration

### Bloom Filter Parameters

- **Dictionary Size**: Should be larger than the number of words you plan to store
- **Hash Functions**: More hash functions reduce false positives but increase computation time
- **Recommended**: For 1000 words, use size 10000 and 7 hash functions

### Suggestion Parameters

- **Max Suggestions**: Maximum number of suggestions to return
- **Max Edit Distance**: Maximum edit distance for candidate generation (1-3 recommended)
- **Similarity Threshold**: Minimum similarity score for suggestions (0.0-1.0)

## Performance

The Bloom filter provides O(k) lookup time where k is the number of hash functions, making it very fast for dictionary lookups. The candidate generation and similarity calculation are the main performance bottlenecks, but the system is optimized for typical use cases.

## Limitations

1. **False Positives**: Bloom filters can have false positives (saying a word exists when it doesn't)
2. **No False Negatives**: If the Bloom filter says a word doesn't exist, it definitely doesn't
3. **Memory vs Accuracy**: Larger Bloom filters reduce false positives but use more memory
4. **Edit Distance**: Very long words with large edit distances can be slow to process

## Testing

Run the tests:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
