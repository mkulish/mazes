package models

import (
	"regexp"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"

)

// Maze represents maze object with grid
// swagger:model Maze
type Maze struct {
	// swagger:ignore
	ID int64 `json:"id"`

	// swagger:ignore
	OwnerID int64 `json:"-"`

	// Entrance cell on the grid
	// required: true
	// type: string
	// pattern: ^[A-Z][1-9][0-9]?$
	// example: A1
	Entrance string `json:"entrance"`

	// Grid size (cols x rows, up to 27x99)
	// required: true
	// pattern: ^([1-9]|1[0-9]|2[0-7])x[1-9][0-9]?$
	// example: 4x3
	GridSize string `json:"gridSize"`

	// Array of wall cells
	// required: true
	// example: ["B2", "B4", "C4"]
	// items.pattern: ^[A-Z][1-9][0-9]?$
	Walls []string `json:"walls"`

	// swagger:ignore
	// temporary fix for storing slice in sqlite
	WallsStr string `json:"-"`

	// precalculated solutions
	// swagger:ignore
	MinPathStr string `json:"-"`
	// swagger:ignore
	MaxPathStr string `json:"-"`
}
// PostGet hook is executed after reading maze from sqlite
func (m *Maze) PostGet(s gorp.SqlExecutor) error {
	if m.WallsStr != "" && len(m.Walls) == 0 {
		// read walls slice from wallsStr column
		m.Walls = strings.Split(m.WallsStr, ",")
	}
	return nil
}
// PreInsert hook is executed before inserting maze into sqlite
func (m *Maze) PreInsert(s gorp.SqlExecutor) error {
	if m.WallsStr == "" && len(m.Walls) > 0 {
		// store walls slice in wallsStr column
		m.WallsStr = strings.Join(m.Walls, ",")
	}
	return nil
}

// MazeResponse represents a JSON reponse with created maze id
// swagger:model MazeResponse
type MazeResponse struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Maze ID
	// required: true
	// type: integer
	ID int64 `json:"id"`
}

// MazeSolutionResponse represents a JSON reponse with maze solution path
// swagger:model MazeSolutionResponse
type MazeSolutionResponse struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Cells path
	// required: true
	Path []string `json:"path"`
}

// MazeSearchResponse represents a JSON reponse with mazes list
// swagger:model MazeSearchResponse
type MazeSearchResponse struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Mazes list
	// required: true
	// type: array
	Items []*Maze `json:"items"`
}

// Validate checks maze data
func (m *Maze) Validate(v *revel.Validation) {
	v.Check(m.OwnerID, revel.Required{})

	v.Check(m.Entrance,
		revel.Required{},
		revel.ValidMatch(regexp.MustCompile("^[A-Z][1-9][0-9]?$")),
	).Key("entrance")

	v.Check(m.GridSize,
		revel.Required{},
		// max size is 99x27
		revel.ValidMatch(regexp.MustCompile("^[1-9][0-9]?x([1-9]|1[0-9]|2[0-7])$")),
	).Key("gridSize")
}
