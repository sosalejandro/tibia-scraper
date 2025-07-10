package pkg

import (
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/sosalejandro/tibia-scraper/domain"
)

func worldGeneratorProducer(worldsChannel chan<- domain.World) {
	defer close(worldsChannel)
	c := colly.NewCollector(colly.AllowedDomains("www.tibia.com", "tibia.com"))

	// Let's say I want to visit all game worlds. The urls are found at tr.Odd and tr.Even, td:nth-child(1) found at the href link
	c.OnHTML("tr.Odd, tr.Even", func(e *colly.HTMLElement) {
		worldName := e.ChildText("td:nth-child(1)")

		// Clean up world name: trim spaces and replace non-breaking spaces
		cleanWorldName := strings.TrimSpace(strings.ReplaceAll(worldName, "\u00A0", " "))

		worldsChannel <- domain.World(cleanWorldName)
	})

	c.OnError(func(r *colly.Response, err error) {
		println("Error:", err.Error())
	})

	err := c.Visit("https://www.tibia.com/community/?subtopic=worlds")
	if err != nil {
		println("Failed to visit URL:", err.Error())
		return
	}
}

func worldGeneratorConsumer(worldsChannel <-chan domain.World, wgcc *domain.WorldGeneratorBroadcaster, wg *sync.WaitGroup) {
	defer wg.Done()
	wgcc.Broacast(worldsChannel)
}
