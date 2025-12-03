package main

import (
	"encoding/json"
	"flag"
	"log"
	"maps"
	"os"
	"runtime"
	"sync"
	"time"

	gonormalizer "github.com/ii7102/schema-normalization/go_normalizer"
	jqnormalizer "github.com/ii7102/schema-normalization/jq_normalizer"
	"github.com/ii7102/schema-normalization/schema"
)

func normalizerOptions() []schema.NormalizerOption {
	return []schema.NormalizerOption{
		schema.WithBooleanFields("isActive", "isActive2", "isActive3", "isActive4"),
		schema.WithIntegerFields("age", "age2", "age3", "age4", "age5"),
		schema.WithStringFields("name", "name2", "name3", "name4", "name5"),
		schema.WithFloatFields("measure", "measure2", "measure3", "measure4", "measure5"),
		schema.WithEnumOfStringFields(
			[]any{"Alice", "Bob", "Charlie"},
			"enumString", "enumString2", "enumString3", "enumString4",
		),
		schema.WithArrayOfBooleanFields("booleanArray", "booleanArray2", "booleanArray3", "booleanArray4", "booleanArray5"),
		schema.WithArrayOfIntegerFields("integerArray", "integerArray2", "integerArray3", "integerArray4", "integerArray5"),
		schema.WithArrayOfStringFields("stringArray", "stringArray2", "stringArray3", "stringArray4", "stringArray5"),
		schema.WithArrayOfFloatFields("floatArray", "floatArray2", "floatArray3", "floatArray4", "floatArray5"),
		schema.WithArrayOfEnumOfStringFields(
			[]any{"Alice", "Bob", "Charlie"},
			"enumStringArray", "enumStringArray2", "enumStringArray3", "enumStringArray4",
		),
	}
}

func main() {
	normalizerType := flag.String("normalizer", "go", "normalizer type: go or jq")
	concurrently := flag.Bool("concurrently", false, "run the normalization concurrently")
	numOfIterations := flag.Int("iterations", 1, "number of iterations")
	batch := flag.Bool("batch", false, "run the normalization in batch mode")

	flag.Parse()

	if *numOfIterations <= 0 {
		log.Fatal("number of iterations must be greater than 0")
	}

	if *batch && *concurrently {
		log.Fatal("batch and concurrently cannot be used together")
	}

	var (
		data       map[string]any
		normalizer schema.AbstractNormalizer
		err        error
	)

	switch *normalizerType {
	case "go":
		normalizer, err = gonormalizer.NewNormalizer(normalizerOptions()...)
	case "jq":
		normalizer, err = jqnormalizer.NewNormalizer(normalizerOptions()...)
	default:
		log.Fatalf("invalid normalizer type: %s", *normalizerType)
	}

	if err != nil {
		log.Fatalf("failed to create normalizer: %v", err)
	}

	input, err := os.ReadFile("test_data.json")
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}

	if err = json.Unmarshal(input, &data); err != nil {
		log.Fatalf("failed to unmarshal input file: %v", err)
	}

	log.Println("--------------------------------")
	log.Println("--- Benchmark Configuration ---")
	log.Printf("Normalizer: %s\n", *normalizerType)
	log.Printf("Concurrently: %t\n", *concurrently)
	log.Printf("Number of iterations: %d\n", *numOfIterations)
	log.Printf("Batch: %t\n", *batch)
	log.Println("--------------------------------")

	normalize(normalizer, data, *concurrently, *batch, *numOfIterations)
}

type MemStats struct {
	AllocMB      float64
	TotalAllocMB float64
	SysMB        float64
	NumGC        uint32
}

func getMemStats() MemStats {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	return MemStats{
		AllocMB:      float64(memStats.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(memStats.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(memStats.Sys) / 1024 / 1024,
		NumGC:        memStats.NumGC,
	}
}

func printMemoryEfficiency(iterations int, memUsed, totalAllocated float64, duration time.Duration) {
	memUsedKB := memUsed * 1024
	totalAllocatedKB := totalAllocated * 1024

	log.Println("--- Memory Efficiency Analysis ---")
	log.Printf("Memory per iteration: %.5f KB\n", memUsedKB/float64(iterations))
	log.Printf("Total allocations per iteration: %.3f KB\n", totalAllocatedKB/float64(iterations))
	log.Printf("Memory efficiency (MB/s): %.3f MB/s\n", totalAllocated/duration.Seconds())
	log.Println("--------------------------------")
}

func normalize(normalizer schema.AbstractNormalizer, data map[string]any, concurrently, batch bool, numIterations int) {
	dataArray := make([]map[string]any, 0, numIterations)

	dataAnyArray := make([]any, 0, numIterations)
	for range numIterations {
		dataClone := maps.Clone(data)
		dataArray = append(dataArray, dataClone)
		dataAnyArray = append(dataAnyArray, dataClone)
	}

	runtime.GC()
	runtime.GC() // Call twice to ensure clean state

	afterPrepStats := getMemStats()
	timeNow := time.Now()

	switch {
	case batch:
		normalizeBatch(normalizer, dataAnyArray)
	case concurrently:
		normalizeConcurrently(normalizer, dataArray)
	default:
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

	printMemoryEfficiency(numIterations, memUsedDuringProcess, totalMemAllocated, duration)

	log.Println("--- Time Efficiency Analysis ---")
	log.Printf("Time per iteration: %s\n", duration/time.Duration(numIterations))
	log.Printf("Processing rate: %.1f iterations/s\n", float64(numIterations)/duration.Seconds())
	log.Printf("Total time: %.3fs\n", duration.Seconds())
	log.Println("--------------------------------")
}

func normalizeSequentially(normalizer schema.AbstractNormalizer, dataArray []map[string]any) {
	for _, data := range dataArray {
		if _, err := normalizer.Normalize(data); err != nil {
			log.Printf("Failed to normalize the data, err: %v\n", err)
		}
	}
}

func normalizeBatch(normalizer schema.AbstractNormalizer, dataArray []any) {
	if _, err := normalizer.NormalizeBatch(dataArray); err != nil {
		log.Printf("Failed to normalize the data batch, err: %v\n", err)
	}
}

func normalizeConcurrently(normalizer schema.AbstractNormalizer, dataArray []map[string]any) {
	var (
		errorsChan = make(chan error, len(dataArray))
		waitGroup  sync.WaitGroup
	)

	waitGroup.Add(len(dataArray))

	for _, data := range dataArray {
		go func(dataElem map[string]any) {
			defer waitGroup.Done()

			if _, err := normalizer.Normalize(dataElem); err != nil {
				errorsChan <- err
			}
		}(data)
	}

	waitGroup.Wait()
	close(errorsChan)

	for err := range errorsChan {
		log.Printf("Failed to normalize the data, err: %v\n", err)
	}
}
