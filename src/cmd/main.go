package main

import (
	"fmt"
	"sync"

	"github.com/sosalejandro/tibia-scraper/domain"
	"github.com/sosalejandro/tibia-scraper/pkg"
)

func main() {
	iterateFor := 20
	killStatServer := domain.World("Esmera")
	_, killStatisticMap := pkg.Scraper()

	esmeraKillStatsAny, _ := killStatisticMap.Load(killStatServer)
	if esmeraKillStatsAny == nil {
		panic("Esmera kill stats not found")
	}

	esmeraKillStats := esmeraKillStatsAny.(*sync.Map)
	// Iterate over the esmeraKillStats map only for the amount specified by iterateFor
	count := 0
	esmeraKillStats.Range(func(key, value interface{}) bool {
		if count >= iterateFor {
			return false
		}
		v := value.(domain.KillStatistic)
		fmt.Printf("Race: %s, LastDay: {KilledPlayers: %d, KilledByPlayers: %d}, LastWeek: {KilledPlayers: %d, KilledByPlayers: %d}\n",
			v.Race,
			v.LastDay.KilledPlayers,
			v.LastDay.KilledByPlayers,
			v.LastWeek.KilledPlayers,
			v.LastWeek.KilledByPlayers)
		count++
		return true
	})
}
