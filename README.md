# jazzdata

**jazzdata** is a Go package for retrieving album artwork from multiple music providers. It is primarily a library package that allows you to extract album metadata from your audio files or directories and then fetch high–quality cover art from sources such as Apple Music, Deezer, Discogs, and iTunes.

An example CLI application is included to demonstrate how to use the package.

## Features

- **Retrieve Album Artwork:** Get cover images using various online providers.
- **Metadata Extraction:** Extract artist and album metadata from your audio files or directories.
- **Modular Design:** Shared types (e.g. `Cover`, `Metadata`, `Source`) are defined in the `models` package.
- **CLI Example:** An example CLI application illustrates how to integrate the package functionality into your projects.

## Installation

Make sure you have Go (version 1.22.2 or later) installed. Then, install the package using:

```bash
go get github.com/rickcollette/jazzdata
```

Or clone the repository:

```bash
git clone https://github.com/rickcollette/jazzdata.git
cd jazzdata
```

## Using the Package in Your Code

Import the package and its providers into your project:

```go
import (
    "fmt"
    "github.com/rickcollette/jazzdata/providers"
)

func main() {
    // Example: Retrieve cover art using the iTunes provider.
    artist := "Adele"
    album := "25"
    itunesProvider := providers.ITunesProvider{}
    covers, err := itunesProvider.GetCovers(artist, album)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    if len(covers) > 0 {
        fmt.Println("Cover URL:", covers[0].CoverURL)
    } else {
        fmt.Println("No cover found!")
    }
}
```

This snippet shows how to use the package to fetch artwork for a given artist and album.

## Example CLI Application

A simple CLI example is included in the `cli` directory. To build and run the CLI:

1. Change to the `cli` directory and build the executable:

    ```bash
    cd cli
    go build -o albumart-cli
    ```

2. Run the CLI with a sample command:

    ```bash
    ./albumart-cli -provider applemusic,deezer /path/to/album
    ```

### CLI Options

- **-provider**  
  Comma–separated list of providers (options: `applemusic`, `deezer`, `discogs`, `itunes`).  
  _Default:_ `"applemusic,deezer"`

- **-cover-name**  
  Name for the downloaded cover image file.  
  _Default:_ `"cover"`

- **-cache**  
  Path to a cache file to avoid processing albums multiple times (used with upgrade).

- **-tag**  
  Comma–separated list of tags. _(Not actively used in this example.)_

- **-recursive**  
  Enable recursive scan in directories.

- **-upgrade**  
  Upgrade existing cover art. _(Upgrade functionality is provided as an example and may need further implementation.)_

- **Other flags** for maximum sizes, strict mode, and warning options are available as well.

## Project Structure

```text
jazzdata/
├── cli/                  # Example CLI application demonstrating package usage.
├── models/               # Shared types (Cover, Metadata, Source, etc.).
├── providers/            # Provider implementations (Apple Music, Deezer, Discogs, iTunes).
├── cache.go              # Cache implementation for upgrade scenarios.
├── metadata.go           # Audio metadata extraction functions.
├── options.go            # Options structure for configuration.
├── utils.go              # Utility functions (downloading images, file operations, etc.).
├── go.mod
├── go.sum
```

## Dependencies

- **[goimagehash](https://github.com/corona10/goimagehash):** For image hashing and similarity checks.
- **[tag](https://github.com/dhowden/tag):** For reading metadata from audio files.
- Standard library packages such as `net/http`, `os`, `flag`, etc.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---
