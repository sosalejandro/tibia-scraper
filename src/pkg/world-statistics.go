package pkg

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/sosalejandro/tibia-scraper/domain"
)

func worldStatisticsProducer(worldsChannel <-chan domain.World, worldStatisticsChannel chan<- domain.KillStatisticMessage, maxConcurrency int) {
	producer := func(server domain.World, worldStatisticsChannel chan<- domain.KillStatisticMessage, wg *sync.WaitGroup, sem chan struct{}) {
		defer wg.Done()
		defer func() { <-sem }() // Release semaphore slot

		serverUrl := fmt.Sprintf("https://www.tibia.com/community/?subtopic=killstatistics&world=%s", server)
		c := colly.NewCollector(colly.AllowedDomains("www.tibia.com", "tibia.com"))

		// Scrape the kill statistics table for each world
		c.OnHTML("#KillStatisticsTable tr.Odd.TextRight.DataRow, #KillStatisticsTable tr.Even.TextRight.DataRow", func(e *colly.HTMLElement) {
			race := e.ChildText("td:nth-child(1)")
			lastDayKilled := parseInt32(e.ChildText("td:nth-child(2)"))
			lastDayKilledBy := parseInt32(e.ChildText("td:nth-child(3)"))
			lastWeekKilled := parseInt32(e.ChildText("td:nth-child(4)"))
			lastWeekKilledBy := parseInt32(e.ChildText("td:nth-child(5)"))

			worldStatisticsChannel <- domain.KillStatisticMessage{
				KS: domain.KillStatistic{
					Race: race,
					LastDay: domain.LastDay{
						KilledPlayers:   lastDayKilled,
						KilledByPlayers: lastDayKilledBy,
					},
					LastWeek: domain.LastWeek{
						KilledPlayers:   lastWeekKilled,
						KilledByPlayers: lastWeekKilledBy,
					},
				},
				World: server,
			}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Error:", err.Error())
		})

		err := c.Visit(serverUrl)
		if err != nil {
			fmt.Println("Failed to visit world URL:", err.Error())
			return
		}
	}

	defer close(worldStatisticsChannel)
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency) // Semaphore channel

	for world := range worldsChannel {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore slot
		go producer(world, worldStatisticsChannel, &wg, sem)
	}
	wg.Wait()
}

var (
	worldStatisticsRegistryInstance *sync.Map
	worldStatisticsRegistryOnce     sync.Once
)

func GetWorldStatisticsRegistry() *sync.Map {
	worldStatisticsRegistryOnce.Do(func() {
		worldStatisticsRegistryInstance = &sync.Map{}
	})
	return worldStatisticsRegistryInstance
}

func worldStatisticsConsumer(worldStatisticsChannel <-chan domain.KillStatisticMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	worldStatisticsRegistry := GetWorldStatisticsRegistry()
	for msg := range worldStatisticsChannel {
		worldStatisticAny, _ := worldStatisticsRegistry.LoadOrStore(msg.World, &sync.Map{})
		worldStatisticMap := worldStatisticAny.(*sync.Map)
		worldStatisticMap.Store(msg.KS.Race, msg.KS)
	}
}
