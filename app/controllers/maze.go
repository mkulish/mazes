package controllers

import (
	"database/sql"
	"strings"

	"golang.org/x/sync/errgroup"
	"github.com/revel/revel"

	"github.com/mkulish/mazes/app/models"
	"github.com/mkulish/mazes/app/services"
)

// Maze controller
type Maze struct {
	App
}

// Search performs mazes search
// swagger:route GET /maze maze searchMazes
//
// Search mazes
//
//     Security:
//       oauth2: read
//
//     Responses:
//       200: MazeSearchResponse
//       400: ValidationError
//       401: UnauthorizedError
//       500: InternalError
func (c Maze) Search() revel.Result {
	user, _ := c.Session.Get("user")
	mazes, err := c.searchMazes(user.(*models.User).ID)
	if err != nil {
		return c.internalError()
	}

	return c.RenderJSON(models.MazeSearchResponse{OK: true, Items: mazes})
}

// Create performs maze validation, processing and insert
// swagger:route POST /maze maze createMaze
//
// Creates a maze, performs validation and path processing
//
//     Parameters:
//     + name: maze
//       in: body
//       description: Maze data
//       required: true
//       type: Maze
//
//     Security:
//       oauth2: write
//
//     Responses:
//       200: MazeResponse
//       400: ValidationError
//       401: UnauthorizedError
//       500: InternalError
func (c Maze) Create(maze models.Maze) revel.Result {	
	user, err := c.Session.Get("user")
	maze.OwnerID = user.(*models.User).ID

	maze.Validate(c.Validation)
	if ! c.Validation.HasErrors() {
		services.ValidateMazeGrid(&maze, c.Validation)
	}
	if ! c.Validation.HasErrors() {
		g := new(errgroup.Group)
		var minPath, maxPath []string
		g.Go(func() error {
			minPath, err = services.SolveMaze(&maze, true)
			return err
		})
		g.Go(func() error {
			maxPath, err = services.SolveMaze(&maze, false)
			return err
		})
		if err := g.Wait(); err != nil {
			c.Validation.Error(err.Error()).Key("walls")
		}

		maze.MinPathStr, maze.MaxPathStr = strings.Join(minPath, ","), strings.Join(maxPath, ",")
	}

	if c.Validation.HasErrors() {
		return c.validationError(c.Validation.Errors)
	}

	err = c.Txn.Insert(&maze)
	if err != nil {
		c.Log.Errorf("maze '%v' insert: %v", maze, err)
		return c.internalError()
	}

	return c.RenderJSON(models.MazeResponse{OK: true, ID: maze.ID})
}

// Solution returns maze solution
// swagger:route GET /maze/{mazeId}/solution maze getMazeSolution
//
// Get maze solution
//
//     Parameters:
//     + name: mazeId
//       in: path
//       description: Maze id
//       required: true
//       type: integer
//       example: 1
//     + name: steps
//       in: query
//       description: return _min_ or _max_ possible steps in solution path
//       required: true
//       type: string
//       example: min
//       pattern: ^min|max$
//
//     Security:
//       oauth2: read
//
//     Responses:
//       200: MazeSolutionResponse
//       400: ValidationError
//       401: UnauthorizedError
//       500: InternalError
func (c Maze) Solution(id int64, steps string) revel.Result {
	user, err := c.Session.Get("user")
	if user == nil || err != nil {
		// user should be injected in the auth interceptor
		return c.internalError()
	}

	if id == 0 {
		c.Validation.Error("Missing or incorrect maze id").Key("id")
		return c.validationError(c.Validation.Errors)
	}

	maze, err := c.getMaze(id)
	if err != nil {
		return c.internalError()
	}
	if maze == nil {
		c.Validation.Error("Not found").Key("id")
		return c.validationError(c.Validation.Errors)
	} else if maze.OwnerID != user.(*models.User).ID {
		return c.unauthorizedError()
	}

	var path []string
	switch steps {
		case "min": path = strings.Split(maze.MinPathStr, ",")
		case "max": path = strings.Split(maze.MaxPathStr, ",")
		default: c.Validation.Error("Should be one of: min, max").Key("steps")
	}

	if c.Validation.HasErrors() {
		return c.validationError(c.Validation.Errors)
	}

	return c.RenderJSON(models.MazeSolutionResponse{OK: true, Path: path})
}

// getUser performs maze lookup by id
func (c App) getMaze(id int64) (*models.Maze, error) {
	maze := &models.Maze{}
	err := c.Txn.SelectOne(maze, c.Db.SqlStatementBuilder.Select("*").From("Maze").Where("ID=?", id))

	if err == sql.ErrNoRows {
		// not found
		return nil, nil
	}
	if err != nil {
		c.Log.Error("Failed to find maze", "id", id, "error", err)
	}
	return maze, err
}

// searchMazes performs mazes search
func (c Maze) searchMazes(ownerID int64) ([]*models.Maze, error) {
	query := c.Db.SqlStatementBuilder.Select("*").From("Maze")
	if ownerID != 0 {
		query = query.Where("OwnerID=?", ownerID)
	}

	var mazes []*models.Maze
	_, err := c.Txn.Select(&mazes, query)

	if err == sql.ErrNoRows {
		// not found
		return nil, nil
	}
	if err != nil {
		c.Log.Error("Failed to search mazes", "error", err)
	}
	return mazes, err
}
