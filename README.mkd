# Jekyll Serve Watcher

    $ jekyll serve --watch

`--watch` option waste many CPU time, because it is implemented by stat()
polling in each a second.

**jsw** replaces it by golang.

## How to compile

    $ go get github.com/howeyc/fsnotify
    $ go build jsw.go

Copy a file `jsw` (or `jsw.exe` for Windows) into one of your PATH.

## Execute

Just type in your jekyll project:

    $ jsw

instead of:

    $ jekyll serve --watch

## Requirements

*   go, of course. (1.1 or above)
*   jekyll (1.0 or above)
