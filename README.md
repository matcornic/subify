# Subify
Subify is a tool to download subtitles for your favorite TV shows and movies.
It is directly able to open the video with your default player, once the subtitle is downloaded.

Subify uses [SubDB Web API](http://thesubdb.com/) to get subtitles. It also considers that you use a default player interpreting srt subtitles when the video file name is the same than the srt file (ex: [VLC](http://www.videolan.org/vlc/)).

## Installing
Download the [last version of Subify](https://github.com/matcornic/subify/releases), and that's it. No need to install something else. Works on Linux, Mac OS (Darwin) and Windows

If you use Golang, you can get Subify and its binary directly with :
```shell
go get -v github.com/matcornic/subify
```

On Mac OS, you can also create a Service Automator, in order to [add a "Subify" option in the Finder menu for your videos](https://github.com/matcornic/subify/wiki/Adding-a-Subify-option-in-the-Finder-menu-for-your-videos).

## Get started
Note : the binary is usable as is, if you want to run the command from anywhere on your OS, make sure to add Subify home installation to your PATH environment variable

```shell
# Download subtitle with default language
subify dl <path_to_your_video>
# Download subtitle with default language, and open video with your default player
subify dl <path_to_your_video> -o
# Download subtitle with french language, and open with your default player
subify dl <path_to_your_video> -o -l fr
```

## Documentation
### Global usage
```
Tool to handle subtitles for your best TV Shows and movies

Usage:
  subify [command]

Available Commands:
  dl          Download the subtitles for your video - 'subify dl --help'

Flags:
      --config string   Config file (default is $HOME/.subify.json). Build a file like this to change default behaviour
      --dev             Instanciate development sandbox instead of production variables
  -h, --help            help for subify
  -v, --verbose         Print more information while executing

Use "subify [command] --help" for more information about a command.
```

### Downloading command
```
Download the subtitles for your video (movie or TV Shows)
Give the path of your video as first parameter and let's go !

Usage:
  subify dl <video-path> [flags]

Aliases:
  dl, download


Flags:
  -l, --language string   Language of the subtitle (default "en")
  -o, --open              Once the subtitle is donwloaded, open the video with your default video player (OSX: "open", Windows: "start", Linux/Other: "xdg-open")

Global Flags:
      --config string   Config file (default is $HOME/.subify.json). Build a file like this to change default behaviour
      --dev             Instanciate development sandbox instead of production variables
  -v, --verbose         Print more information while executing
```

## Release Notes
* **0.1.0** Jan 15, 2016
  * Implement firs init

## Contributing

You've developed a cool feature or fixed a bug ?

1. Fork the repo
2. Create an issue  
3. Create your feature branch (git checkout -b my-new-feature)
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

## TODO
1. List of favorite languages (Downloads the first to match)
2. Auto update command
3. Upload command to contribute to SubDB database
4. Vendoring
5. Doc on default configuration (for example to change the default language for all downloads)
6. Use OpenSubtitles API and/or Addic7ed (better quality of translations, but no API)
7. Localization/Internationalization

## License
Subify is released under the Apache 2.0 license. See LICENSE.txt
