package dymean_test

import (
	"fmt"
	"github.com/bi0dread/dymean"
	"math"
	"testing"
)

// ExampleDidYouMean demonstrates basic usage of the DidYouMean package
func ExampleDidYouMean() {
	// Create a new DidYouMean instance
	// Parameters: dictionary size (10000), number of hash functions (7)
	dym := dymean.NewDidYouMean(10000, 7)

	// Add some words to the dictionary
	words := []string{
		"hello", "world", "golang", "programming", "computer", "science",
		"algorithm", "data", "structure", "bloom", "filter", "spell",
		"checker", "dictionary", "suggestion", "similarity", "distance",
		"edit", "levenshtein", "candidate", "generation", "typo", "keyboard",
	}

	dym.AddWords(words)

	// Test with a correct word
	fmt.Printf("Is 'hello' correct? %t\n", dym.IsCorrect("hello"))

	// Test with a misspelled word
	fmt.Printf("Is 'helo' correct? %t\n", dym.IsCorrect("helo"))

	// Get suggestions for a misspelled word
	suggestions := dym.GetSuggestions("helo", 3, 2)
	fmt.Printf("Suggestions for 'helo': ")
	for i, suggestion := range suggestions {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s (%.2f)", suggestion.Word, suggestion.Similarity)
	}
	fmt.Println()

	// Get the best suggestion
	bestSuggestion := dym.Suggest("progamming")
	fmt.Printf("Best suggestion for 'progamming': %s\n", bestSuggestion)

	// Check and suggest
	isCorrect, suggestions := dym.CheckAndSuggest("algoritm")
	fmt.Printf("Is 'algoritm' correct? %t\n", isCorrect)
	if !isCorrect && len(suggestions) > 0 {
		fmt.Printf("Did you mean: %s?\n", suggestions[0].Word)
	}

	// Output:
	// Is 'hello' correct? true
	// Is 'helo' correct? false
	// Suggestions for 'helo': hello (0.80)
	// Best suggestion for 'progamming': programming
	// Is 'algoritm' correct? false
	// Did you mean: algorithm?
}

// TestDidYouMeanBasic tests basic functionality
func TestDidYouMeanBasic(t *testing.T) {
	dym := dymean.NewDidYouMean(1000, 5)

	words := []string{"hello", "world", "test", "go", "programming"}
	dym.AddWords(words)

	// Test correct words
	if !dym.IsCorrect("hello") {
		t.Error("Expected 'hello' to be correct")
	}

	if !dym.IsCorrect("world") {
		t.Error("Expected 'world' to be correct")
	}

	// Test incorrect words
	if dym.IsCorrect("helo") {
		t.Error("Expected 'helo' to be incorrect")
	}

	if dym.IsCorrect("wrld") {
		t.Error("Expected 'wrld' to be incorrect")
	}
}

// TestSuggestions tests suggestion functionality
func TestSuggestions(t *testing.T) {
	dym := dymean.NewDidYouMean(1000, 5)

	words := []string{"hello", "help", "hell", "world", "word", "work"}
	dym.AddWords(words)

	suggestions := dym.GetSuggestions("helo", 3, 2)

	if len(suggestions) == 0 {
		t.Error("Expected at least one suggestion for 'helo'")
	}

	// Check that suggestions are sorted by similarity
	for i := 1; i < len(suggestions); i++ {
		if suggestions[i-1].Similarity < suggestions[i].Similarity {
			t.Error("Suggestions should be sorted by similarity (descending)")
		}
	}

	// Check that the best suggestion is reasonable
	bestSuggestion := dym.Suggest("helo")
	if bestSuggestion == "" {
		t.Error("Expected a non-empty suggestion")
	}
}

// TestLevenshteinDistance tests the Levenshtein distance calculation
func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"hello", "hello", 0},
		{"hello", "helo", 1},
		{"hello", "world", 4},
		{"kitten", "sitting", 3},
		{"saturday", "sunday", 3},
	}

	for _, test := range tests {
		result := dymean.LevenshteinDistance(test.s1, test.s2)
		if result != test.expected {
			t.Errorf("dymean.LevenshteinDistance(%q, %q) = %d, expected %d",
				test.s1, test.s2, result, test.expected)
		}
	}
}

// TestSimilarity tests the similarity calculation
func TestSimilarity(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected float64
	}{
		{"hello", "hello", 1.0},
		{"hello", "helo", 0.8},
		{"hello", "world", 0.2},
		{"", "", 1.0},
		{"a", "b", 0.0},
	}

	for _, test := range tests {
		result := dymean.CalculateSimilarity(test.s1, test.s2)
		// Use a small tolerance for floating point comparison
		if math.Abs(result-test.expected) > 0.001 {
			t.Errorf("dymean.CalculateSimilarity(%q, %q) = %.2f, expected %.2f",
				test.s1, test.s2, result, test.expected)
		}
	}
}

// TestBloomFilter tests the Bloom filter functionality
func TestBloomFilter(t *testing.T) {
	bf := dymean.NewBloomFilter(1000, 5)

	words := []string{"hello", "world", "test", "go"}
	bf.AddWords(words)

	// Test that added words are found
	for _, word := range words {
		if !bf.Contains(word) {
			t.Errorf("Expected Bloom filter to contain %q", word)
		}
	}

	// Test that non-added words might not be found (false negatives are impossible)
	nonWords := []string{"xyz", "abc", "def"}
	for _, word := range nonWords {
		// Note: Bloom filter might have false positives, so we can't test for definite absence
		// We can only test that it doesn't crash
		bf.Contains(word)
	}
}

// BenchmarkDidYouMean benchmarks the suggestion performance
func BenchmarkDidYouMean(b *testing.B) {
	dym := dymean.NewDidYouMean(10000, 7)

	// Add a larger dictionary
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%d", i)
	}
	dym.AddWords(words)

	// Add some common English words
	commonWords := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "i",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "she", "or", "an",
		"will", "my", "one", "all", "would", "there", "their", "what",
		"so", "up", "out", "if", "about", "who", "get", "which", "go",
		"me", "when", "make", "can", "like", "time", "no", "just", "him",
		"know", "take", "people", "into", "year", "your", "good", "some",
		"could", "them", "see", "other", "than", "then", "now", "look",
		"only", "come", "its", "over", "think", "also", "back", "after",
		"use", "two", "how", "our", "work", "first", "well", "way", "even",
		"new", "want", "because", "any", "these", "give", "day", "most",
		"us", "is", "was", "are", "been", "has", "had", "were", "said",
		"each", "which", "their", "said", "if", "will", "up", "other",
		"about", "out", "many", "then", "them", "these", "so", "some",
		"her", "would", "make", "like", "into", "him", "time", "has",
		"two", "more", "go", "no", "way", "could", "my", "than", "first",
		"been", "call", "who", "its", "now", "find", "long", "down",
		"day", "did", "get", "come", "made", "may", "part",
	}
	dym.AddWords(commonWords)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dym.GetSuggestions("helo", 5, 2)
	}
}
