package main

import "fmt"
import "./sudoku"

func main() {
	master := make(sudoku.Group, sudoku.NumAgents)
	for i := 0; i < sudoku.NumAgents; i++ {
		master[i] = sudoku.New(i)
	}
	
	grps := makeGroups(master)
	setPuzzle(master)
  sudoku.Solve(master, grps)

  for _, a := range master {
    fmt.Print(a.Val, ",")
  }
  fmt.Print("\n")
}

// user defined
func makeGroups(master sudoku.Group) (grps []sudoku.Group) {
  grps = make([]sudoku.Group, 2 * sudoku.NumVals)
  for i, _ := range grps {
    grps[i] = make(sudoku.Group, sudoku.NumVals)
  }
  for i := 0; i < sudoku.NumVals; i++ {
    for j := 0; j < sudoku.NumVals; j++ {
      n := i * sudoku.NumVals + j
      col := n % sudoku.NumVals
      row := n / sudoku.NumVals
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

