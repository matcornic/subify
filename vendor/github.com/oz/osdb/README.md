[![Build Status](https://travis-ci.org/oz/osdb.svg?branch=master)](https://travis-ci.org/oz/osdb)

This is a Go client, and library for
[OpenSubtitles](http://opensubtitles.org/).

This project has not reached version `0.1` yet, and it will change in many
breaking ways. But of course, you are welcome to check it out, and participate:
report issues when you find those, or check the "TODO" section down there.


# Install

If you are only interested in (yet another) CLI interface for OpenSubtitles,
then install the latest [Go](http://golang.org) release, and run:

```
go get github.com/oz/osdb && go install github.com/oz/osdb/cmd/osdb
```

Provided that you setup your Go environment correctly, you now have a basic
`osdb` command to interact with OpenSubtitles' API.

```
$ osdb --help
Usage:
    osdb get [--lang=<lang>] <file>
    osdb (put|upload) <movie_file> <sub_file>
    osdb imdb <query>...
    osdb imdb show <movie id>
    osdb -h | --help
    osdb --version
```

Hence, to try to download french or english subtitles for a sample file:

```
$ osdb get --lang fra,eng sample.avi
- Getting fra subtitles for file: sample.avi
- No subtitles found!
- Getting eng subtitles for file: sample.avi
- Downloading to: sample.srt
```


# Hack...

The generated documentation for this package is available at:
http://godoc.org/github.com/oz/osdb

To get started...

 * Install with `go get -d github.com/oz/osdb`,
 * and import `"github.com/oz/osdb"` in your Go code,
 * or try some of the examples in the README.

To access the OpenSubtitles' XML-RPC server you first need to allocate a
client, and then use it to login (even anonymously) in order to receive a
session token. With that, you are finally be allowed to talk. Here is a short
example:

```go
package main

import "github.com/oz/osdb"

func main() {
	c, err := osdb.NewClient()
	if err != nil {
		// ...
	}

	// Anonymous login will set c.Token when successful
	if err = c.LogIn("", "", ""); err != nil {
		// ...
	}

	// etc.
}

```

# Basic examples

## Getting a user session token

Although this library tries to be simple, to use OpenSubtitles' API you need to
login first so as to receive a session token: without it you will not be able
to call any API method.

```go
c, err := osdb.NewClient()
if err != nil {
	// ...
}

err := c.LogIn("user", "password", "language")
if err != nil {
	// ...
}
// c.Token is now set, and subsequent API calls will not be refused.
```

However, you do not need to register a user. To login anonymously, just leave
the `user` and `password` parameters blank:

```go
c.LogIn("", "", "")
```

## Searching subtitles

Subtitle search can be done in a number of ways: using special file-hashes,
IMDB movie IDs, or even using full-text queries. Hash-based search will
generally yield the best results, so this is what is demoed next. However, in
order to search with this method, you *must* have a movie file to hash.

```go
path := "/path/to/movie.avi"
languages := []string{"eng"}

// Hash movie file, and search...
res, err := client.FileSearch(path, languages)
if err != nil {
	// ...
}

for _, sub := range res {
	fmt.Printf("Found %s subtitles file \"%s\" at %s\n",
		sub.LanguageName, sub.SubFileName, sub.ZipDownloadLink)
}
```

## Downloading subtitles

Let's say you have just made a search, for example using `FileSearch()`, and as
the API provided a few results, you decide to pick one for download:

```go
subs, err := c.FileSearch(...)

// Download subtitle file, and write to disk using subs[0].SubFileName
if err := c.Download(&subs[0]); err != nil {
	// ...
}

// Alternatively, use the filename of your choice:
if err := c.DownloadTo(&subs[0], "safer-name.srt"); err != nil {
	// ...
}
```

## Checking if a subtitle exists

Before trying to upload an allegedly "new" subtitles file to OSDB, you should
always check whether they already have it.

As some movies fit on more than one "CD" (remember those?), you will need to
use the `Subtitles` type (note the *s*?), one per subtitle file:

```go
subs := osdb.Subtitles{
		{
			SubHash:       subHash,       // md5 hash of subtitle file
			SubFileName:   subFileName,
			MovieHash:     movieHash,     // see osdb.Hash()
			MovieByteSize: movieByteSize, // careful, it's a string...
			MovieFileName: movieFileName,
		},
}
```

Then simply feed that to `HasSubtitles`, and you will be done.

```go
found, err := c.HasSubtitles(subs)
if err != nil {
	// ...
}
```

## Hashing a file

OSDB uses a custom checksum-hash to identify movie files. If you ever need
these:

```go
hash, err := osdb.Hash("somefile.avi")
if err != nil {
	// ...
}
fmt.Println("hash: %x\n", hash)
```


# On user agents...

If you have read OSDB's [developer documentation][osdb], you should notice that
you need to register an "official" user agent in order to use their API.

By default this library will present itself with the "osdb-go" agent, which is
fine for me. However, if you need to change this, simply set the client's
`UserAgent` with:

```go
c, err := osdb.NewClient()
if err != nil {
	// ...
}
c.UserAgent = "My custom user agent"
```

# TODO

 - [ ] Move docs from README to godoc.org
 - [ ] [Full API coverage][apidocs]:
  - [x] LogIn
  - [x] LogOut
  - [x] NoOperation
  - [x] SearchSubtitles by hash
  - [x] SearchSubtitles by IMDB IDs
  - [ ] SearchToMail
  - [x] DownloadSubtitles
  - [x] TryUploadSubtitles
  - [ ] UploadSubtitles
  - [x] SearchMoviesOnIMDB
  - [x] GetIMDBMovieDetails
  - [ ] InsertMovie
  - [ ] ServerInfo
  - [ ] ReportWrongMovieHash
  - [ ] SubtitlesVote
  - [ ] AddComment
  - [ ] GetSubLanguages
  - [ ] DetectLanguage
  - [ ] GetAvailableTranslations
  - [ ] GetTranslation
  - [ ] AutoUpdate
  - [ ] CheckMovieHash
  - [ ] CheckSubHash


# License

BSD, see the LICENSE file.

[osdb]: http://trac.opensubtitles.org/projects/opensubtitles
[apidocs]: http://trac.opensubtitles.org/projects/opensubtitles/wiki/XmlRpcIntro#XML-RPCmethods
