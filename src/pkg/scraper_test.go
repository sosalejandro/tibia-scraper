package pkg

import (
	"runtime"
	"testing"
)

func TestScraper(t *testing.T) {
	globalCharactersRegistry, worldStatisticsRegistry := Scraper()
	if globalCharactersRegistry == nil {
		t.Error("globalCharactersRegistry is nil")
	}
	if worldStatisticsRegistry == nil {
		t.Error("worldStatisticsRegistry is nil")
	}
}

func BenchmarkScraper(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var memStatsStart, memStatsEnd runtime.MemStats
		runtime.ReadMemStats(&memStatsStart)

		Scraper()

		runtime.ReadMemStats(&memStatsEnd)
		alloc := memStatsEnd.TotalAlloc - memStatsStart.TotalAlloc
		b.Logf("Memory allocated: %d bytes", alloc)
	}
}
