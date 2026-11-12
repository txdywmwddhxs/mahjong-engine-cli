package replay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// View uses snapshots only: old rules can be viewed without running old code.
// Frame 0 is the initial state. Frame numbers are stable across UI languages.
func View(f *File, in io.Reader, out io.Writer, step bool) error {
	show := func(index int) {
		state := f.Initial
		kind := "initial"
		if index > 0 {
			state = f.Frames[index-1].State
			kind = f.Frames[index-1].Kind
		}
		fmt.Fprintf(out, "Frame %d/%d | turn %d | %s\nHand: %s\nFixed: %s\nWall: %d | waiting: %t | result: %d\n", index, len(f.Frames), state.Turn, kind, state.Hand.Display(), state.Fixed.Display(), len(state.Wall), state.Waiting, state.Result)
		if index > 0 {
			frame := f.Frames[index-1]
			if frame.Request != nil {
				b, _ := json.Marshal(frame.Request)
				fmt.Fprintf(out, "Request: %s\n", b)
			}
			if frame.Action != nil {
				b, _ := json.Marshal(frame.Action)
				fmt.Fprintf(out, "Action: %s\n", b)
			}
			if frame.Tile != "" {
				fmt.Fprintf(out, "Tile: %s\n", frame.Tile)
			}
			if frame.Accepted != nil {
				fmt.Fprintf(out, "Accepted: %t\n", *frame.Accepted)
			}
			if frame.Settlement != nil {
				b, _ := json.Marshal(frame.Settlement)
				fmt.Fprintf(out, "Settlement: %s\n", b)
			}
		}
	}
	if !step {
		for i := 0; i <= len(f.Frames); i++ {
			show(i)
		}
		return nil
	}
	scanner := bufio.NewScanner(in)
	index := 0
	for {
		show(index)
		fmt.Fprint(out, "[Enter/n] next | p previous | g FRAME | t TURN | q quit > ")
		if !scanner.Scan() {
			return scanner.Err()
		}
		fields := strings.Fields(strings.ToLower(scanner.Text()))
		if len(fields) == 0 {
			if index < len(f.Frames) {
				index++
			}
			continue
		}
		switch fields[0] {
		case "q":
			return nil
		case "n":
			if index < len(f.Frames) {
				index++
			}
		case "p":
			if index > 0 {
				index--
			}
		case "g", "t":
			if len(fields) != 2 {
				fmt.Fprintln(out, "Expected an integer.")
				continue
			}
			n, err := strconv.Atoi(fields[1])
			if err != nil || n < 0 {
				fmt.Fprintln(out, "Expected a non-negative integer.")
				continue
			}
			if fields[0] == "g" {
				if n <= len(f.Frames) {
					index = n
				} else {
					fmt.Fprintln(out, "Frame not found.")
				}
			} else {
				found := false
				if n == 0 {
					index = 0
					found = true
				} else {
					for i, frame := range f.Frames {
						if frame.State.Turn == n {
							index = i + 1
							found = true
							break
						}
					}
				}
				if !found {
					fmt.Fprintln(out, "Turn not found.")
				}
			}
		default:
			fmt.Fprintln(out, "Unknown command.")
		}
	}
}
