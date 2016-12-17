package main

import "fmt"

type Cell struct {
	Id int
	// Vals is the list of values the cell could possibly be.
	Vals map[int]bool
}

func NewCell(vals []int) *Cell {
	valsMap := map[int]bool{}
	for _, val := range vals {
		valsMap[val] = true
	}
	return &Cell{Vals: valsMap}
}

func (c *Cell) IsSolved() bool {
	for _, solved := range c.Vals {
		if !solved {
			return false
		}
	}
	return true
}

func (c *Cell) Remove(vals ...int) {
	for _, val := range vals {
		c.Vals[val] = false

	}
}

type ConstrGroup struct {
	Cells []*Cell
	// Val is the value that one of the cells in the group must have.
	Val int
}

type Grid struct {
	Groups [][]*Cell
	Cells  []*Cell
}

func colGroup(col, n int, cells []*Cell) (group []*Cell) {
	ncols := n * n
	for i := 0; i < ncols; i++ {
		group = append(group, cells[i*ncols+col])
	}
	return group
}

func boxGroup(i, n int, cells []*Cell) (group []*Cell) {
	ncols := n * n
	rowOffset := i / n * (n * ncols)
	colOffset := i % n * n
	for j := 0; j < ncols; j++ {
		group = append(group, cells[rowOffset+colOffset+ncols*(j/n)+j%n])
	}
	return group
}

func rowGroup(row, n int, cells []*Cell) (group []*Cell) {
	nrows := n * n
	for i := 0; i < nrows; i++ {
		group = append(group, cells[row*nrows+i])
	}
	return group
}

func NewGrid(n int) *Grid {
	nvals, ncells := n*n, n*n*n*n

	vals := make([]int, nvals)
	for i := range vals {
		vals[i] = i + 1
	}

	cells := make([]*Cell, ncells)
	for i := range cells {
		cells[i] = NewCell(vals)
		cells[i].Id = i + 1
	}

	groups := [][]*Cell{}
	for i := 0; i < n*n; i++ {
		groups = append(groups, boxGroup(i, n, cells))
		groups = append(groups, rowGroup(i, n, cells))
		groups = append(groups, colGroup(i, n, cells))
	}

	return &Grid{Groups: groups, Cells: cells}
}

func (g *Grid) Step() {
	for _, group := range g.Groups {

	}
}

func main() {
	n := 2
	grid := NewGrid(n)

	for i, cell := range grid.Cells {
		if i%(n*n) == 0 {
			fmt.Println()
		}
		fmt.Printf(" %2d", cell.Id)
	}
	fmt.Println()

	for i, group := range grid.Groups {
		fmt.Printf("group %2d: ", i)
		for _, cell := range group {
			fmt.Printf("%2d ", cell.Id)

		}
		fmt.Println()
	}
}
