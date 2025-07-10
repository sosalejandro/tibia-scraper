## Documentation: Usage of CSP (Communicating Sequential Processes)

### Overview

This project utilizes the Communicating Sequential Processes (CSP) model, as implemented in Go, to orchestrate concurrent web scraping and data processing tasks. CSP enables safe, efficient, and scalable communication between independent components (goroutines) using channels.

### How CSP is Used

- **Producers and Consumers:**  
  Each major data type (worlds, players, statistics) has a dedicated producer goroutine that scrapes data and sends it through a channel. Consumer goroutines receive data from these channels and process or store it.
- **Channels:**  
  Channels are the primary means of communication between goroutines. For example, `worldsChannel` is used to send world names, while `worldPlayersChannel` and `worldStatisticsChannel` are used for player and statistics data.
- **Broadcasting:**  
  The `WorldGeneratorBroadcaster` struct allows broadcasting world names to multiple consumers, ensuring both player and statistics scrapers receive the same data.
- **Synchronization:**  
  The `sync.WaitGroup` is used to wait for all goroutines to finish. Semaphores (buffered channels) are used to limit concurrency and avoid overwhelming the target website.
- **Thread-Safe Registries:**  
  Data is stored in `sync.Map` registries, ensuring safe concurrent access.

### Perks of Using CSP

- **Safe Concurrency:**  
  CSP avoids race conditions by sharing data only through channels, not shared memory.
- **Scalability:**  
  The pattern allows easy scaling by adjusting the number of concurrent goroutines.
- **Separation of Concerns:**  
  Producers and consumers are decoupled, making the codebase modular and easier to maintain.
- **Error Isolation:**  
  Each goroutine can handle its own errors, improving robustness.

### Example Flow

1. The world generator scrapes all world names and sends them through `worldsChannel`.
2. The broadcaster sends each world name to both the player and statistics producers.
3. Each producer scrapes its respective data and sends results through dedicated channels.
4. Consumers receive data, process it, and store it in thread-safe registries.

---

## ADR: Use of CSP (Communicating Sequential Processes)

### Context

The Tibia Scraper project requires concurrent scraping and processing of data from multiple pages on the Tibia website. The solution must be efficient, safe, and maintainable.

### Decision

We chose to use the Communicating Sequential Processes (CSP) model, as implemented in Go, as the core architectural pattern. This involves using goroutines for concurrent tasks and channels for communication and synchronization.

### Status

Accepted

### Consequences

**Pros:**
- Enables efficient, parallel data scraping and processing.
- Reduces risk of data races and synchronization bugs.
- Makes the codebase modular and testable.
- Easily extensible for new data types or consumers.

**Cons:**
- Requires careful channel management to avoid deadlocks.
- Debugging concurrent code can be more complex.

### Usage

- All data flows between producers and consumers are handled via channels.
- WaitGroups and semaphores are used for synchronization and concurrency control.
- The broadcaster pattern is used to distribute data to multiple consumers.

### Alternatives Considered

- Shared memory with mutexes (rejected due to higher risk of race conditions and complexity).
- Sequential scraping (rejected due to inefficiency).

