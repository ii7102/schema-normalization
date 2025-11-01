package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"runtime"
	"sync"
	"time"

	goNormalizer "diploma/go"
	jqNormalizer "diploma/jq"
	"diploma/rules"
)

var (
	normalizerType  = flag.String("normalizer", "go", "normalizer type: go or jq")
	concurrently    = flag.Bool("concurrently", false, "run the normalization concurrently")
	numOfIterations = flag.Int("iterations", 1, "number of iterations")
)

func normalizerOptions() []rules.NormalizerOption {
	return []rules.NormalizerOption{
		rules.WithBooleanFields("isActive", "isActive2", "isActive3", "isActive4"),
		rules.WithIntegerFields("age", "age2", "age3", "age4", "age5"),
		rules.WithStringFields("name", "name2", "name3", "name4", "name5"),
		rules.WithFloatFields("measure", "measure2", "measure3", "measure4", "measure5"),
		rules.WithEnumOfStringFields([]any{"Alice", "Bob", "Charlie"}, "enumString", "enumString2", "enumString3", "enumString4"),
		rules.WithArrayOfBooleanFields("booleanArray", "booleanArray2", "booleanArray3", "booleanArray4", "booleanArray5"),
		rules.WithArrayOfIntegerFields("integerArray", "integerArray2", "integerArray3", "integerArray4", "integerArray5"),
		rules.WithArrayOfStringFields("stringArray", "stringArray2", "stringArray3", "stringArray4", "stringArray5"),
		rules.WithArrayOfFloatFields("floatArray", "floatArray2", "floatArray3", "floatArray4", "floatArray5"),
		rules.WithArrayOfEnumOfStringFields([]any{"Alice", "Bob", "Charlie"}, "enumStringArray", "enumStringArray2", "enumStringArray3", "enumStringArray4"),
		rules.WithArrayOfEnumOfStringFields([]any{"Alice", "Bob", "Charlie"}, "enumStringArray", "enumStringArray2", "enumStringArray3", "enumStringArray4"),
	}
}

func main() {
	flag.Parse()

	input, err := os.ReadFile("test_data.json")
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
		os.Exit(1)
	}

	var data map[string]any
	if err := json.Unmarshal(input, &data); err != nil {
		log.Fatalf("failed to unmarshal input file: %v", err)
		os.Exit(1)
	}

	var normalizer rules.AbstractNormalizer
	switch *normalizerType {
	case "go":
		normalizer, err = goNormalizer.NewNormalizer(normalizerOptions()...)
	case "jq":
		normalizer, err = jqNormalizer.NewNormalizer(normalizerOptions()...)
	default:
		log.Fatalf("invalid normalizer type: %s", *normalizerType)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("failed to create normalizer: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Normalizer: %s\n", *normalizerType)
	fmt.Printf("Concurrently: %t\n", *concurrently)
	fmt.Printf("Number of iterations: %d\n", *numOfIterations)

	normalize(normalizer, data, *concurrently)
}

type MemStats struct {
	AllocMB      float64
	TotalAllocMB float64
	SysMB        float64
	NumGC        uint32
}

func getMemStats() MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemStats{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
	}
}

func printMemoryEfficiency(iterations int, memUsed, totalAllocated float64, duration time.Duration) {
	fmt.Printf("\n--- Memory Efficiency Analysis ---\n")
	fmt.Printf("Memory per iteration: %.2f KB\n", (memUsed*1024)/float64(iterations))
	fmt.Printf("Total allocations per iteration: %.2f KB\n", (totalAllocated*1024)/float64(iterations))
	fmt.Printf("Memory efficiency (MB/s): %.2f MB/s\n", totalAllocated/duration.Seconds())
	fmt.Printf("Processing rate: %.0f iterations/s\n", float64(iterations)/duration.Seconds())
}

func normalize(normalizer rules.AbstractNormalizer, data map[string]any, concurrently bool) {
	dataArray := make([]map[string]any, 0, *numOfIterations)
	for range *numOfIterations {
		dataArray = append(dataArray, maps.Clone(data))
	}

	runtime.GC()
	runtime.GC() // Call twice to ensure clean state

	afterPrepStats := getMemStats()

	timeNow := time.Now()
	if concurrently {
		normalizeConcurrently(normalizer, dataArray)
	} else {
		normalizeSequentially(normalizer, dataArray)
	}
	duration := time.Since(timeNow)

	// Get memory stats immediately after processing
	afterProcessStats := getMemStats()

	// Force GC to see actual memory freed
	runtime.GC()
	runtime.GC()

	// Calculate memory usage during processing
	memUsedDuringProcess := afterProcessStats.AllocMB - afterPrepStats.AllocMB
	totalMemAllocated := afterProcessStats.TotalAllocMB - afterPrepStats.TotalAllocMB

	printMemoryEfficiency(*numOfIterations, memUsedDuringProcess, totalMemAllocated, duration)

	fmt.Printf("\n time: %v\n", duration)
}

func normalizeSequentially(normalizer rules.AbstractNormalizer, dataArray []map[string]any) {
	var err error

	for _, data := range dataArray {
		if _, err = normalizer.Normalize(data); err != nil {
			log.Printf("Failed to normalize the data, err: %v\n", err)
		}
	}
}

func normalizeConcurrently(normalizer rules.AbstractNormalizer, dataArray []map[string]any) {
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(dataArray))

	for _, data := range dataArray {
		wg.Add(1)
		go func(dataElem map[string]any) {
			defer wg.Done()
			if _, err := normalizer.Normalize(dataElem); err != nil {
				errorsChan <- err
			}
		}(data)
	}

	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		log.Printf("Failed to normalize the data, err: %v\n", err)
	}
}
