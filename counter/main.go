package counter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func countWords(filename string) (map[string]int, error) {
	counts := make(map[string]int)

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	words := strings.Fields(string(file))

	for _, word := range words {
		counts[strings.ToLower(word)]++
	}

	return counts, nil
}

type result struct {
	filename string
	counts   map[string]int
}

func Run() {
	files := []string{"assets/texts/file1.txt", "assets/texts/file2.txt"}
	results := make(chan result)
	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(filename string) {
			defer wg.Done()
			count, err := countWords(filename)
			if err != nil {
				fmt.Println("Error counting words in", filename, err)
				return
			}
			results <- result{filename, count}
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	merged := make(map[string]int)
	for result := range results {
		for word, count := range result.counts {
			merged[word] += count
		}
	}

	if err := writeJSON("assets/outputs/counts.json", merged); err != nil {
		fmt.Println("Error writing JSON file:", err)
	}

	// to get the highest count
	highestCount := 0
	highestWord := ""
	for word, count := range merged {
		if count > highestCount {
			highestCount = count
			highestWord = word
		}
	}
	fmt.Println("Highest count:", highestCount, "for word:", highestWord)
}

func writeJSON(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoding := json.NewEncoder(file)
	encoding.SetIndent("", "  ")
	return encoding.Encode(data)
}
