package pkg

import (
	"sync"

	"github.com/sosalejandro/tibia-scraper/domain"
)

func Scraper() (*sync.Map, *sync.Map) {
	worldsChannel, staticticsWorldsChannel, playerWorldsChannel :=
		make(chan domain.World), make(chan domain.World), make(chan domain.World)

	worldPlayersChannel := make(chan domain.CharacterMessage)
	worldStatisticsChannel := make(chan domain.KillStatisticMessage)

	var wg sync.WaitGroup
	wg.Add(1)
	go worldGeneratorConsumer(worldsChannel, &domain.WorldGeneratorBroadcaster{Channels: []chan<- domain.World{
		staticticsWorldsChannel, playerWorldsChannel,
	}}, &wg)

	wg.Add(1)
	go worldPlayerConsumer(worldPlayersChannel, &wg)

	wg.Add(1)
	go worldStatisticsConsumer(worldStatisticsChannel, &wg)

	go worldGeneratorProducer(worldsChannel)
	go worldPlayerProducer(playerWorldsChannel, worldPlayersChannel, 10)
	go worldStatisticsProducer(staticticsWorldsChannel, worldStatisticsChannel, 10)
	wg.Wait()

	globalCharactersRegistry := GetGlobalCharacterRegistry()
	worldStatisticsRegistry := GetWorldStatisticsRegistry()
	return globalCharactersRegistry, worldStatisticsRegistry
}
