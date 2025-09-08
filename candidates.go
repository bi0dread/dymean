package dymean

import (
	"strings"
	"unicode"
)

// CandidateGenerator generates possible corrections for misspelled words
type CandidateGenerator struct {
	alphabet string
}

// NewCandidateGenerator creates a new candidate generator
func NewCandidateGenerator() *CandidateGenerator {
	return &CandidateGenerator{
		alphabet: "abcdefghijklmnopqrstuvwxyz",
	}
}

// GenerateCandidates generates possible corrections for a word
func (cg *CandidateGenerator) GenerateCandidates(word string, maxDistance int) []string {
	candidates := make(map[string]bool)
	word = strings.ToLower(word)

	// Generate candidates with different edit distances
	for distance := 1; distance <= maxDistance; distance++ {
		cg.generateCandidatesAtDistance(word, distance, candidates)
	}

	// Convert map to slice
	result := make([]string, 0, len(candidates))
	for candidate := range candidates {
		result = append(result, candidate)
	}

	return result
}

// generateCandidatesAtDistance generates candidates at a specific edit distance
func (cg *CandidateGenerator) generateCandidatesAtDistance(word string, distance int, candidates map[string]bool) {
	if distance == 0 {
		candidates[word] = true
		return
	}

	// Generate deletions
	for i := 0; i < len(word); i++ {
		deleted := word[:i] + word[i+1:]
		if distance == 1 {
			candidates[deleted] = true
		} else {
			cg.generateCandidatesAtDistance(deleted, distance-1, candidates)
		}
	}

	// Generate insertions
	for i := 0; i <= len(word); i++ {
		for _, char := range cg.alphabet {
			inserted := word[:i] + string(char) + word[i:]
			if distance == 1 {
				candidates[inserted] = true
			} else {
				cg.generateCandidatesAtDistance(inserted, distance-1, candidates)
			}
		}
	}

	// Generate substitutions
	for i := 0; i < len(word); i++ {
		for _, char := range cg.alphabet {
			if char != rune(word[i]) {
				substituted := word[:i] + string(char) + word[i+1:]
				if distance == 1 {
					candidates[substituted] = true
				} else {
					cg.generateCandidatesAtDistance(substituted, distance-1, candidates)
				}
			}
		}
	}

	// Generate transpositions (swapping adjacent characters)
	for i := 0; i < len(word)-1; i++ {
		transposed := word[:i] + string(word[i+1]) + string(word[i]) + word[i+2:]
		if distance == 1 {
			candidates[transposed] = true
		} else {
			cg.generateCandidatesAtDistance(transposed, distance-1, candidates)
		}
	}
}

// GenerateCommonTypos generates candidates based on common typing errors
func (cg *CandidateGenerator) GenerateCommonTypos(word string) []string {
	candidates := make(map[string]bool)
	word = strings.ToLower(word)

	// Common keyboard layout for QWERTY
	keyboard := map[rune][]rune{
		'q': {'w', 'a'}, 'w': {'q', 'e', 'a', 's'}, 'e': {'w', 'r', 's', 'd'},
		'r': {'e', 't', 'd', 'f'}, 't': {'r', 'y', 'f', 'g'}, 'y': {'t', 'u', 'g', 'h'},
		'u': {'y', 'i', 'h', 'j'}, 'i': {'u', 'o', 'j', 'k'}, 'o': {'i', 'p', 'k', 'l'},
		'p': {'o', 'l'}, 'a': {'q', 'w', 's', 'z'}, 's': {'a', 'w', 'e', 'd', 'x', 'z'},
		'd': {'s', 'e', 'r', 'f', 'c', 'x'}, 'f': {'d', 'r', 't', 'g', 'v', 'c'},
		'g': {'f', 't', 'y', 'h', 'b', 'v'}, 'h': {'g', 'y', 'u', 'j', 'n', 'b'},
		'j': {'h', 'u', 'i', 'k', 'm', 'n'}, 'k': {'j', 'i', 'o', 'l', 'm'},
		'l': {'k', 'o', 'p'}, 'z': {'a', 's', 'x'}, 'x': {'z', 's', 'd', 'c'},
		'c': {'x', 'd', 'f', 'v'}, 'v': {'c', 'f', 'g', 'b'}, 'b': {'v', 'g', 'h', 'n'},
		'n': {'b', 'h', 'j', 'm'}, 'm': {'n', 'j', 'k'},
	}

	// Generate candidates by replacing each character with adjacent keyboard characters
	for i, char := range word {
		if neighbors, exists := keyboard[char]; exists {
			for _, neighbor := range neighbors {
				candidate := word[:i] + string(neighbor) + word[i+1:]
				candidates[candidate] = true
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(candidates))
	for candidate := range candidates {
		result = append(result, candidate)
	}

	return result
}

// IsValidWord checks if a word contains only valid characters
func IsValidWord(word string) bool {
	if len(word) == 0 {
		return false
	}

	for _, char := range word {
		if !unicode.IsLetter(char) {
			return false
		}
	}

	return true
}
