package services

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"

	"github.com/mkulish/mazes/app/models"
)

// ValidateMazeGrid validates maze grid
func ValidateMazeGrid(m *models.Maze, v *revel.Validation) {
	v.Check(m.Walls, revel.Required{})

	if ! v.HasErrors() {
		width, height := size(m)
		start := parseCell(m.Entrance)
		if start.x < 0 || start.x >= width || start.y < 0 || start.y >= height {
			v.Error("Incorrect entrance: %s", m.Entrance).Key("entrance")
		}

		walls := make(map[string]struct{}, len(m.Walls))
		for _, rawCell := range m.Walls {
			// validity and bounds check
			cell := parseCell(rawCell)
			if cell.x < 0 || cell.x >= width || cell.y < 0 || cell.y >= height {
				v.Error("Incorrect wall cell: %s", rawCell).Key("walls")
			}

			// entrance check
			if cell.x == start.x && cell.y == start.y {
				v.Error("Entrance is a wall cell: %s", rawCell).Key("entrance")
			}

			// duplicates check
			if _, found := walls[rawCell]; found {
				v.Error("Duplicate wall cell: %s", rawCell).Key("walls")
			}
			walls[rawCell] = struct{}{}
		}
	}
}

// SolveMaze returns maze solution path (if any)
// complexity: x * y * log(x * y)
func SolveMaze(m *models.Maze, min bool) ([]string, error) {
	width, heigth := size(m)

	// init walls and explored cells matrices
	walls, explored := make([][]bool, width), make([][]bool, width)
	for i := 0; i < width; i++ {
		walls[i], explored[i] = make([]bool, heigth), make([]bool, heigth)
	}
	for _, rawCell := range m.Walls {
		cell := parseCell(rawCell)
		walls[cell.x][cell.y] = true
	}

	// potential exit coordinate
	// max steps path lookup will find different solution before finish to detect multiple exits 
	exitX := -1
	var solution *cell

	start := parseCell(m.Entrance)
	h := cellHeap{[]*cell{&start}, min}
	heap.Init(&h)

	for h.Len() > 0 {
		// pop closest/farest cell from priority queue for shortest/longest path search
		next := heap.Pop(&h).(*cell)
		explored[next.x][next.y] = true
		if next.y == heigth - 1 {
			solution = next
			break
		}

		// add all adjacent cells to the priority queue
		if next.x > 0 && ! walls[next.x - 1][next.y] && ! explored[next.x - 1][next.y] {
			heap.Push(&h, &cell{x: next.x - 1, y: next.y, parent: next, steps: next.steps + 1})
		}
		if next.y > 0 && ! walls[next.x][next.y - 1] && ! explored[next.x][next.y - 1] {
			heap.Push(&h, &cell{x: next.x, y: next.y - 1, parent: next, steps: next.steps + 1})
		}
		if next.x < width - 1 && ! walls[next.x + 1][next.y] && ! explored[next.x + 1][next.y] {
			heap.Push(&h, &cell{x: next.x + 1, y: next.y, parent: next, steps: next.steps + 1})
		}
		if next.y < heigth - 1 && ! walls[next.x][next.y + 1] && ! explored[next.x][next.y + 1] {
			heap.Push(&h, &cell{x: next.x, y: next.y + 1, parent: next, steps: next.steps + 1})

			if next.y == heigth - 2 {
				// approaching exit row
				if exitX >= 0 && exitX != next.x {
					// another potential exit path was already found in max steps lookup
					return nil, fmt.Errorf("Maze has multiple exits")
				}
				exitX = next.x
			}
		}
	}

	if solution == nil {
		return nil, fmt.Errorf("Maze doesn't have a solution")
	}

	// trace back and reverse result path
	var res []string
	for solution != nil {
		res = append(res, encodeCell(*solution))
		solution = solution.parent
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res, nil
}

type cell struct {
	x, y	int
	steps	int
	parent	*cell
}
func parseCell(rawCell string) cell {
	x := int(rawCell[0] - 'A')
	y,_ := strconv.Atoi(rawCell[1:])
	// 0-based: A2 -> x:0, y:1
	return cell{x: x, y: y-1}
}
func encodeCell(cell cell) string {
	return fmt.Sprintf("%c%d", rune('A' + cell.x), cell.y+1)
}

type cellHeap struct {
	items []*cell
	min bool
}
func (h *cellHeap) Len() int {
	return len(h.items)
}
func (h *cellHeap) Less(i, j int) bool {
	var res bool
	if h.items[i].y == h.items[j].y {
		// cells with the same distance to exit are sorted by their path steps count
		res = h.items[i].steps > h.items[j].steps
	} else {
		// otherwise - by distance to exit
		res = h.items[i].y < h.items[j].y
	}

	if h.min {
		// reverse order to choose between shortest/longest path
		return ! res
	}
	return res
}
func (h cellHeap) Swap(i, j int){
	h.items[i], h.items[j] = h.items[j], h.items[i]
}
func (h *cellHeap) Push(x any) {
	h.items = append(h.items, x.(*cell))
}
func (h *cellHeap) Pop() any {
	n := len(h.items)
	x := h.items[n - 1]
	h.items[n - 1] = nil
	h.items = h.items[: n - 1]
	return x
}

func size(m *models.Maze) (int, int) {
	size := strings.Split(m.GridSize, "x")
	width, _ := strconv.Atoi(size[1])
	height, _ := strconv.Atoi(size[0])
	return width, height
}
