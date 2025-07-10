package pkg

import (
	"fmt"
	"strings"
	"sync"

	"strconv"

	"github.com/gocolly/colly"
	"github.com/sosalejandro/tibia-scraper/domain"
)

func parseInt32(s string) int32 {
	s = strings.ReplaceAll(s, ",", "")
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func worldPlayerProducer(worldsChannel <-chan domain.World, worldPlayersChannel chan<- domain.CharacterMessage, maxConcurrency int) {
	producer := func(server domain.World, worldPlayersChannel chan<- domain.CharacterMessage, wg *sync.WaitGroup, sem chan struct{}) {
		defer wg.Done()
		defer func() { <-sem }() // Release semaphore slot

		serverUrl := fmt.Sprintf("https://www.tibia.com/community/?subtopic=worlds&world=%s", server)
		c := colly.NewCollector(colly.AllowedDomains("www.tibia.com", "tibia.com"))

		c.OnHTML("tr.Odd, tr.Even", func(e *colly.HTMLElement) {
			name := e.ChildText("td:nth-child(1)")
			level := e.ChildText("td:nth-child(2)")
			vocation := e.ChildText("td:nth-child(3)")

			// Clean up name: trim spaces and replace non-breaking spaces
			cleanName := strings.TrimSpace(strings.ReplaceAll(name, "\u00A0", " "))

			// Parse level from string to int32
			levelInt := parseInt32(level)

			characterMsg := domain.CharacterMessage{
				Character: domain.Character{
					Name:     cleanName,
					Level:    levelInt,
					Vocation: vocation,
					World:    server, // Add the world name to the character
				},
			}
			worldPlayersChannel <- characterMsg
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

	defer close(worldPlayersChannel)
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency) // Semaphore channel

	for world := range worldsChannel {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore slot
		go producer(world, worldPlayersChannel, &wg, sem)
	}
	wg.Wait()
}

var (
	globalCharacterRegistryInstance *sync.Map
	globalCharacterRegistryOnce     sync.Once

	worldCharacterRegistryInstance *sync.Map
	worldCharacterRegistryOnce     sync.Once
)

func GetGlobalCharacterRegistry() *sync.Map {
	globalCharacterRegistryOnce.Do(func() {
		globalCharacterRegistryInstance = &sync.Map{}
	})
	return globalCharacterRegistryInstance
}

func GetWorldCharacterRegistry() *sync.Map {
	worldCharacterRegistryOnce.Do(func() {
		worldCharacterRegistryInstance = &sync.Map{}
	})
	return worldCharacterRegistryInstance
}

func worldPlayerConsumer(worldPlayersChannel <-chan domain.CharacterMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	globalCharacterRegistry := GetGlobalCharacterRegistry()
	worldCharacterRegistry := GetWorldCharacterRegistry()
	for characterMsg := range worldPlayersChannel {
		character := characterMsg.Character
		globalCharacterRegistry.Store(character.Name, character)
		worldMapAny, _ := worldCharacterRegistry.LoadOrStore(character.World, &sync.Map{})
		worldMap := worldMapAny.(*sync.Map)
		worldMap.Store(character.Name, character)
	}
}
