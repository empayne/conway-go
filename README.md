Conway's Game of Life implemented using Golang.
Uses termbox-go UI library.

Sample usage:
* go run conway.go ui.go
* go run conway.go ui.go -size=20
* go run conway.go ui.go -input="glider.txt"

Press any key to exit.

Some possible optimizations include:
* Use the Hashlife memoized algorithm (instead of the current naive implementation).
* Use a single buffer to store the grid + list of changes (instead of the current double-buffered approach).

Some future release may contain these features:
* Select between toroidal/non-toroidal array.
* Select density of live cells when randomly generating the grid.
* Select speed of UI update (instead of current 100ms delay).
* Show current generation number on the UI.
