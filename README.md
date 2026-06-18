# Dash

Dash is a high-performance shell autocomplete engine built in Go.

Instead of relying on simple linear scans or in-memory maps, Dash uses a custom B+ Tree indexing engine to provide fast prefix search, ranking, and command discovery across shell history, file paths, aliases, and developer workflows.

The project is designed as a practical exploration of search engines, indexing systems, and storage engine concepts through a real-world developer tool.

---

## Features

### Command Autocomplete

```bash
$ dash git ch
```

Suggestions:

```text
git checkout
git cherry-pick
```

---

### Prefix Search

Fast command lookup using B+ Tree range scans.

```bash
$ dash dock
```

Suggestions:

```text
docker ps
docker build
docker logs
docker compose up
```

---

### History-Based Ranking

Frequently used commands appear first.

```text
git checkout      score: 542
git cherry-pick   score: 12
```

---

### File Path Completion

```bash
$ dash src/
```

Suggestions:

```text
src/main.go
src/config.yaml
src/internal/
```

---

### Learning From Usage

Dash continuously updates command rankings from shell history.

Supported:

* Bash History
* Zsh History

Future:

* Fish Shell
* PowerShell

---

### Fuzzy Matching

Typo-tolerant search.

```bash
$ dash gco
```

Suggestions:

```text
git checkout
```

---

### Persistent Index

Indexes survive process restarts.

```text
index.db
```

No need to rebuild on every launch.

---

## Architecture

```text
             Shell History
                    │
                    ▼
              Command Parser
                    │
                    ▼
                Indexer
                    │
                    ▼
               B+ Tree Index
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
  Prefix Search           Range Search
        │                       │
        └───────────┬───────────┘
                    ▼
             Ranking Engine
                    │
                    ▼
              Suggestions
```

---

## Why B+ Trees?

Dash is built around a custom B+ Tree implementation.

Benefits:

* O(log n) lookup
* Efficient range queries
* Natural prefix search
* Ordered keys
* Storage-engine friendly
* Disk persistence support

Example:

```text
git add
git branch
git checkout
git commit
git merge
```

Query:

```text
git ch
```

Dash finds the first matching key and scans linked leaf nodes to retrieve all matching commands.

This is the same principle used by many database indexes.

---

## Data Model

### Command Entry

```go
type Command struct {
    Text       string
    Frequency  uint64
    LastUsed   time.Time
}
```

### B+ Tree Key

```text
Command String
```

### B+ Tree Value

```text
Command Metadata
```

---

## Project Structure

```text
dash/
├── cmd/
│   └── dash/
├── internal/
│   ├── bptree/
│   ├── indexer/
│   ├── parser/
│   ├── ranking/
│   ├── search/
│   ├── storage/
│   └── history/
├── data/
├── benchmarks/
├── docs/
└── README.md
```

---

## Search Flow

```text
User Input
     │
     ▼
Prefix Query
     │
     ▼
B+ Tree Range Scan
     │
     ▼
Candidate Commands
     │
     ▼
Ranking Engine
     │
     ▼
Suggestions
```

---

## Roadmap

### Phase 1

* B+ Tree implementation
* Search
* Insert
* Delete
* Range Scan

### Phase 2

* Command indexing
* History parsing
* Ranking engine

### Phase 3

* File path indexing
* Persistent storage
* Serialization

### Phase 4

* Fuzzy matching
* Context-aware suggestions
* Interactive shell mode

### Phase 5

* Shell integration
* Real-time indexing
* Plugin system

---

## Learning Goals

Dash explores concepts commonly found in:

* Search Engines
* Database Indexes
* Storage Engines
* Information Retrieval Systems
* Shell Tooling
* Systems Programming

Key topics:

* B+ Trees
* Prefix Search
* Range Queries
* Ranking Algorithms
* Disk Persistence
* Indexing Pipelines
* Query Processing

---

