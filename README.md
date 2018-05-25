# duitdraw
shiny backend for duit

# About
[Duit](https://github.com/mjl-/duit) is a ui toolkit for [go](https://golang.org).
As a backend it uses plan9 port's /dev/draw emulation via a [go interface](https://github.com/9fans/go/tree/master/draw).

This has a some drawbacks:
- heavy run time dependency
- no easy windows support
- no simple deployment with a single static binary

As a first step `duitdraw` is a drop-in replacement for `github.com/9fans/go/draw`, using a backend based on [shiny](https://github.com/golang/exp/tree/master/shiny). This has the advantage, that no changes are needed for duit.

Once this becomes a valid alternative to the original drawing backend, duit could be changed to interface better with shiny.

The scope of the package is not a full implementation of `9fans/go/draw`. Everything that is not needed by duit in the initial release state is removed.


# Usage
To try the backend, copy the content of this repository to `$GOPATH/src/github.com/9fans/go/draw` and recompile duit.

# Current state
This is just a very basic first first release and tested only on windows.
Please test and comment.

- fonts (right now Go regular is embedded)
	- plan9 style or ttf path?
	- ttf: freetype or golang.org/x/image/font/sfnt?
- client windows cannot be closed
	- is it a bug in shiny (windows)?
	- how to propagate close requests? By a channel or with Release?
- drawing
	- only a simple line algorithm is implemented
	- which rasterization should be used, freetype or golang.org/x/image/vector?
	- general line rasterizer missing
	- Arc, FillArc missing (ellipse.go)
- clipboard
	- uses atotto's, is that ok?
- mouse movement
	- uses as/cursor, is that ok?
	- inner window offset is hard coded
- shiny
	- window flickers on resize, are we using shiny the wrong way?
- program crashes when the window is closed
