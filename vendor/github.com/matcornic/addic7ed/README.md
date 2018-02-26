# Addic7ed API

[![Build Status](https://travis-ci.org/matcornic/addic7ed.svg?branch=master)](https://travis-ci.org/matcornic/addic7ed)
[![Go Report Card](https://goreportcard.com/badge/github.com/matcornic/addic7ed)](https://goreportcard.com/report/github.com/matcornic/addic7ed)
[![Go Coverage](https://codecov.io/github/matcornic/addic7ed/coverage.svg)](https://codecov.io/github/matcornic/addic7ed/)
[![Godoc](https://godoc.org/github.com/matcornic/addic7ed?status.svg)](https://godoc.org/github.com/matcornic/addic7ed)

`addic7ed` is a Golang package to get subtitles from [Addic7ed](http://www.addic7ed.com/) website. As Addic7ed website does not provide a proper API yet, this package uses search feature of website and scraps HTML results to build data.

## Installation

As any golang package, just download it with `go get`.

```bash
go get -u github.com/matcornic/addic7ed
```

## Usage

### Searching all subtitles of a given TV show

```golang
c := addic7ed.New()
show, err := c.SearchAll("Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]") // Usually the name of the video file
if err != nil {
    panic(err)
}
fmt.Println(show.Name) // Output: Shameless (US) - 08x11 - A Gallagher Pedicure
fmt.Println(show.Subtitles) // Output: all subtitles with version, languages and download links
```

In order to find all the subtitles, this API:

1. Use `search.php` page of Addic7ed API
1. Parse the results

It means that if the tv show name is not precise enough, this API will not be able to find the exact TV show page.

### Searching the best subtitle of a given TV show

```golang
c := addic7ed.New()
showName, subtitle, err := c.SearchBest("Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]", "English")
if err != nil {
    panic(err)
}
fmt.Println(showName) // Output: Shameless (US) - 08x11 - A Gallagher Pedicure
fmt.Println(subtitle) // Output: the best suitable subtitle in English language
fmt.Println(subtitle.Version) // Output: BATV
fmt.Println(subtitle.Language) // Output: English

// Download the subtitle to a given file name
err := subtitle.DownloadTo("Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv].srt")
if err != nil {
    panic(err)
}
```

In order to search the best subtitle, this API:

1. Filters subtitles of the given language. Here: `English`
1. Scores similarities between the name of the show and available versions (combining Jaro-winkler distance and an internal weight)
    1. It means that the name of the show has to contain the `version`. Here: `BATV`
1. Choose the version with the best score
1. Choose the best subtitle of the chosen version (the most updated one)

### Helper functions

Some helper functions are provided to adapt `subtitles` structure to the context

```golang
c := addic7ed.New()
show, err := c.SearchAll("Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]") // Usually the name of the video file
if err != nil {
    panic(err)
}

// Filter subtitles to keep only english subtitles
subtitles = show.Subtitles.Filter(WithLanguage("English"))
// Group by version
subtitlesByVersion = subtitles.GroupByVersion()
fmt.Println(subtitlesByVersion["BATV"]) // Output: print all english subtitles of BATV version
```

Available filter functions:

- `WithLanguage`
- `WithVersion`
- `WithVersionRegexp`

Available groupBy functions:

- `GroupByVersion`
- `GroupByLanguage`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Licence

MIT. This package is not affiliated with Addic7ed website.