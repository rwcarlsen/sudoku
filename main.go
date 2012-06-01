package main

import "fmt"
import "./sudoku"

const (
  numAgents = 9
  numVals = 3
)

func main() {
  p := sudoku.New(numAgents, numVals)
  p.Solve(setPuzzle, makeGroups)
  fmt.Println(p)
}

// user defined
func makeGroups(master sudoku.Group) (grps []sudoku.Group) {
  grps = make([]sudoku.Group, 2 * numVals)
  for i, _ := range grps {
    grps[i] = make(sudoku.Group, numVals)
  }
  for i := 0; i < numVals; i++ {
    for j := 0; j < numVals; j++ {
      n := i * numVals + j
      col := n % numVals
      row := n / numVals
      grps[col][row] = master[n]
      grps[row + 3][col] = master[n]
    }
  }
	return grps
}

// user defined
func setPuzzle(g sudoku.Group) {
	g[0].SetVal(1)
	g[3].SetVal(3)
	g[7].SetVal(1)
}

