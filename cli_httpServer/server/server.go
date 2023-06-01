package server

import (
	"fmt"
	"net/http"
)

func findLongestSubstring(s string) string {
	seen := make(map[rune]bool)
	var start, maxLength, maxStart, maxEnd int

	for i, c := range s {
		if seen[c] {
			for j := start; j < i; j++ {
				if []rune(s)[j] == c {
					start = j + 1
					break
				}
				delete(seen, rune(s[j]))
			}
		} else {
			seen[c] = true
		}

		if i-start > maxLength {
			maxLength = i - start
			maxStart = start
			maxEnd = i
		}
	}

	return s[maxStart : maxEnd+1]
}

func SubstringHandler(w http.ResponseWriter, r *http.Request) {
	str := r.FormValue("str")
	if len(str) == 0 {
		http.Error(w, "Empty string", http.StatusBadRequest)
		return
	}
	for i := 0; i < len([]rune(str)); i++ {
		if string(str[i]) == " " {
			http.Error(w, "Empty string", http.StatusBadRequest)
			return
		}
	}

	longestSubstring := findLongestSubstring(str)
	fmt.Fprint(w, longestSubstring)
}
