# Demo exercise

An exercise project aiming to solve a code challenge and explore the [Revel framework](https://revel.github.io/)

TODO:
 - use dedicaded SQL server
 - env and dependencies injection
 - simple UI for maze editing / solutions visualization

# Exercise brief

Create an API (consuming and producing JSON) which allows users to register, persist mazes they create in the database and get solutions for those mazes. Please read the following instructions carefully, not following them will reflect negatively on your application. You should implement the necessary logic for the following flow:

1. User registers via POST /user endpoint with the following fields:
    1. username (i.e. happyUser)
    2. password (i.e. iTk19!n)
2. User logs in via POST /login endpoint
3. The API creates a session and responds with a token
4. **From this point on all of the mentioned endpoints should require a valid token to be supplied**
5. User creates a maze via a POST /maze endpoint with the following fields:
    1. gridSize (size of a maze grid i.e. 10x10)
    2. walls (an array of cells which contain a wall within a given grid)
    3. entrance (the cell where the path should begin i.e. A1)
    
    Note: the grid uses capital letters for columns and numbers for rows (i.e. A1)
    
6. User sends a request to GET /maze/{mazeId}/solution endpoint with steps query parameter which can be either min or max
7. The API returns an array of grid cells leading from the entrance of the maze to the exit of the maze with the following rules: 
    1. if steps parameter is min the API returns the path from the entrance to the exit with the **least number of steps possible**
    2. if steps parameter is max the API returns the path from the entrance to the exit with the **most number of steps possible**
    3. At each step, the API moves from one empty cell to an adjacent empty cell (horizontally or vertically, but not diagonally)
8. User can see their created mazes by sending a request to GET /maze (the user should be able to see just their own mazes)

**Example:**

POST /maze request body:

```
{
    "entrance":  "A1",
    "gridSize": "8x8",
    "walls": ["C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8"]
}
```

GET /maze/1/solution?steps=min response:

```
{
    "path": ["A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "A8"]
}
```

**Additional notes:**

- If the maze has no solution an error should be thrown
- API needs to detect the exit point automatically
- A maze can only have one exit point (at one of the bottom edge cells of the grid), otherwise an error should be thrown

### Start the web server locally on http://localhost:9000/:

   revel run

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).
