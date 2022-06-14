package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel/testing"

	"github.com/mkulish/mazes/app/models"
)

var (
	validAuth = ""
	validMazeWithSolution1 = models.Maze{
		Entrance: "A1",
		GridSize: "4x3",
		Walls: []string{"B2", "B4", "C4"},
	}
)

// MazeTest contains basic integration tests for /maze endpoints
type MazeTest struct {
	testing.TestSuite
}

// Before called on every test
func (t *MazeTest) Before() {
	// Revel framework has poor support for test DB etc, temporary using existing endpoints
	if validAuth == "" {
		t.Post("/user", "application/json", strings.NewReader("{\"username\": \"test\", \"password\": \"12345\"}"))

		var resp models.LoginResponse
		json.Unmarshal(t.ResponseBody, &resp)
		validAuth = resp.Token
	}
}

// TestCreateShouldReturnUnauthorized ...
func (t *MazeTest) TestCreateShouldReturnUnauthorized() {
	t.Post("/maze", "application/json", nil)
	t.AssertStatus(401)
}

// TestCreateShouldCreateMaze ...
func (t *MazeTest) TestCreateShouldCreateMaze() {
	// should create new maze
	t.postObject(t.BaseUrl()+"/maze", validMazeWithSolution1)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")

	// should return new maze
	t.authGet(t.BaseUrl()+"/maze")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")

	var resp models.MazeSearchResponse
	json.Unmarshal(t.ResponseBody, &resp)
	t.AssertEqual(len(resp.Items), 1)
}

// TestCreateShouldValidateGrid ...
func (t *MazeTest) TestCreateShouldValidateGrid() {
	// should check entrance
	invalidMaze := validMazeWithSolution1

	invalidMaze.Entrance = ""
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	invalidMaze.Entrance = "A:1"
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	invalidMaze.Entrance = "A1000"
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	// should check gridSize
	invalidMaze = validMazeWithSolution1

	invalidMaze.GridSize = ""
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	invalidMaze.GridSize = "25/100"
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)
}

// TestCreateShouldValidateWalls ...
func (t *MazeTest) TestCreateShouldValidateWalls() {
	invalidMaze := validMazeWithSolution1

	// wall at the entrance
	invalidMaze.Walls = []string{invalidMaze.Entrance}
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	// duplicate wall cells
	invalidMaze.Walls = []string{"A2", "A2"}
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	// incorrect wall cell
	invalidMaze.Walls = []string{"A:2", "A200"}
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)
}

// TestCreateShouldValidateSolution ...
func (t *MazeTest) TestCreateShouldValidateSolution() {
	invalidMaze := validMazeWithSolution1

	// no solution
	invalidMaze.Walls = []string{"A2", "B2", "C2"}
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)

	// multiple exits
	invalidMaze.Walls = []string{"B2"}
	t.postObject(t.BaseUrl()+"/maze", invalidMaze)
	t.AssertStatus(400)
}

// TestSolutionShouldReturnUnauthorized ...
func (t *MazeTest) TestSolutionShouldReturnUnauthorized() {
	t.Get("/maze/1/solution?steps=min")
	t.AssertStatus(401)
}

// TestSolutionShouldReturnSolution ...
func (t *MazeTest) TestSolutionShouldReturnSolution() {
	t.postObject(t.BaseUrl()+"/maze", validMazeWithSolution1)
	t.AssertOk()

	var resp models.MazeSolutionResponse
	var mazeResp models.MazeResponse
	json.Unmarshal(t.ResponseBody, &mazeResp)
	id := mazeResp.ID

	// should return min solution
	t.authGet(fmt.Sprintf("%s/maze/%d/solution?steps=min", t.BaseUrl(), id))
	json.Unmarshal(t.ResponseBody, &resp)
	t.AssertEqual(resp.Path, []string{"A1", "A2", "A3", "A4"})

	// should return max solution
	t.authGet(fmt.Sprintf("%s/maze/%d/solution?steps=max", t.BaseUrl(), id))
	json.Unmarshal(t.ResponseBody, &resp)
	t.AssertEqual(resp.Path, []string{"A1", "B1", "C1", "C2", "C3", "B3", "A3", "A4"})
}

// TestSolutionShouldValidateParams ...
func (t *MazeTest) TestSolutionShouldValidateParams() {
	// should check maze id
	t.authGet(fmt.Sprintf("%s/maze/-/solution?steps=min", t.BaseUrl()))
	t.AssertStatus(400)

	t.authGet(fmt.Sprintf("%s/maze/0/solution?steps=min", t.BaseUrl()))
	t.AssertStatus(400)

	t.authGet(fmt.Sprintf("%s/maze/999/solution?steps=min", t.BaseUrl()))
	t.AssertStatus(400)

	// should check steps query param
	t.authGet(fmt.Sprintf("%s/maze/1/solution", t.BaseUrl()))
	t.AssertStatus(400)

	t.authGet(fmt.Sprintf("%s/maze/1/solution?steps=any", t.BaseUrl()))
	t.AssertStatus(400)
}

func (t *MazeTest) postObject(url string, obj any) {
	data, _ := json.Marshal(obj)
	req := t.PostCustom(url, "application/json", bytes.NewReader(data))
	req.Header.Add("Authorization", "Bearer " + validAuth)
	req.Send()
}


func (t *MazeTest) authGet(url string) {
	req := t.GetCustom(url)
	req.Header.Add("Authorization", "Bearer " + validAuth)
	req.Send()
}
