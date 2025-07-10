package domain

type KillStatisticData struct {
	KilledPlayers   int32 `json:"killedPlayers"`
	KilledByPlayers int32 `json:"killedByPlayers"`
}

type LastDay KillStatisticData

type LastWeek KillStatisticData

type KillStatistic struct {
	Race     string   `json:"race"`
	LastDay  LastDay  `json:"lastDay"`
	LastWeek LastWeek `json:"lastWeek"`
}

type World string

type SiteUrl string

type URLS struct {
	WorldPlayers    *SiteUrl `json:"worldPlayers"`
	WorldStatistics *SiteUrl `json:"worldStatistics"`
}

type KillStatisticMessage struct {
	KS    KillStatistic
	World World
}

type Character struct {
	Name     string `json:"name"`
	World    World  `json:"world"`
	Level    int32  `json:"level"`
	Vocation string `json:"vocation"`
}

type CharacterMessage struct {
	Character Character `json:"character"`
}

// ConsumerChannels holds the channels for communication between producers and consumers
// in the domain model. It is used to send messages related to kill statistics and character data.
// This struct is designed to be used in the consumer part of the application, where it receives
// messages from producers that gather data from external sources like Tibia's website.
// The channels are defined as send-only (chan<-) to ensure that only the consumer can send messages,
// while the producer can only receive them. This helps maintain a clear separation of concerns
type WorldGeneratorBroadcaster struct {
	Channels []chan<- World
}

func (wgcc *WorldGeneratorBroadcaster) close() {
	for _, channel := range wgcc.Channels {
		close(channel)
	}
	wgcc.Channels = nil
}

func (wgcc *WorldGeneratorBroadcaster) Broacast(worldsChannel <-chan World) {
	defer wgcc.close()
	for world := range worldsChannel {
		for _, channel := range wgcc.Channels {
			channel <- world
		}
	}
}
