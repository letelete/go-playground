// For online execution environment, see https://tour.golang.org/concurrency/8

package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

func SendTreeValuesToChannel(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	SendTreeValuesToChannel(t.Left, ch)
	ch <- t.Value
	SendTreeValuesToChannel(t.Right, ch)
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	SendTreeValuesToChannel(t, ch)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		ch1V, ch1Ok := <-ch1
		ch2V, ch2Ok := <-ch2
		if !ch1Ok && !ch2Ok {
			break
		} else if ch1Ok != ch2Ok || ch1V != ch2V {
			return false
		}
	}
	return true
}

// SameTrees determines whether all of the trees contain the same values
func SameTrees(trees []*tree.Tree) bool {
	channels := make([]chan int, 0, len(trees))
	for _, tr := range trees {
		channels = append(channels, make(chan int))
		go Walk(tr, channels[len(channels) - 1])
	}
	for {
		val, ok := <-channels[0]
		for i := 1; i < len(channels); i++ {
			nextVal, nextOk := <-channels[i]
			if nextVal != val || nextOk != ok {
				return false
			}
		}
		if !ok {
			break
		}
	}
	return true
}

func main() {
	equalTrees := []*(tree.Tree){tree.New(1), tree.New(1)}
	differentTrees := []*(tree.Tree){tree.New(1), tree.New(2)}

	// Same(*tree.Tree, *tree.Tree)
	fmt.Printf("Same(): Expected: true, Got: %t\n", Same(equalTrees[0], equalTrees[1]))
	fmt.Printf("Same(): Expected: false, Got: %t\n", Same(differentTrees[0], differentTrees[1]))

	// SameTrees([]*tree.Tree)
	fmt.Printf("SameTrees(): Expected: true, Got: %t\n", SameTrees(equalTrees))
	fmt.Printf("SameTrees(): Expected: false, Got: %t\n", SameTrees(differentTrees))
}
