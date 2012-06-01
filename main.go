package main

import "fmt"
import "strconv"

const (
  numAgents = 9
	numVals = 3
)

func main() {
	ch := make(tunnel)
	
	master := make(group, numAgents)
	for i := 0; i < numAgents; i++ {
		master[i] = New(i, ch)
	}
	
	grps := makeGroups(master)
	assignGroups(grps)
	setPuzzle(master)
	dispatch(master)

	for count := 0; count < numAgents; {
		select {
			case <-ch:
				count++
		}
	}

  for _, a := range master {
    fmt.Print(a.val, ",")
  }
  fmt.Print("\n")
}

type tunnel chan msg

type group []*agent
func (g group) JustOneWith(val int) bool {
  count := 0
  for _, a := range g {
    if a.ops[val] {
      count++
    }
  }
  return count == 1
}
func (g *group) String() string {
	s := "group["
	for _, a := range *g {
		s += "id=" + strconv.Itoa(a.id) + ", "
	}
	s += "]"
	return s
}

type options map[int]bool
func (ops *options) String() string {
	s := "options["
	for val, op := range *ops {
		if op {
			s += strconv.Itoa(val) + ", "
		}
	}
	s += "]"
	return s
}

type msg struct {
	id int
	val int
}

type agent struct {
	id int
	val int
	ch tunnel
	done tunnel
	grps []group
	ops options
}

func New(id int, done tunnel) *agent {
	ops := make(options)
	for j := 1; j <= numVals; j++ {
		ops[j] = true
	}
	return &agent{
			id: id,
			done: done,
			grps: []group{},
			ops: ops,
			ch: make(tunnel, numAgents),
		}
}

func (a *agent) String() string {
	s := "agent:: id=" + strconv.Itoa(a.id) + ", val=" + strconv.Itoa(a.val)
	s += "\n"
	for _, g := range a.grps {
		s += g.String() + "\n"
	}
	s += a.ops.String()
	return s
}

func (a *agent) SetVal(val int) {
	for v, _ := range a.ops {
		a.ops[v] = false
	}
	a.ops[val] = true
	a.val = val
}

func (a *agent) AddGroup(g group) {
	a.grps = append(a.grps, g)
}

func (a *agent) markTaken(m msg) {
	a.ops[m.val] = false
	a.checkSolved()
}

func (a *agent) checkSolved() {
	opCount := 0
	finalVal := 0
	for val, op := range a.ops {
    if op {
      opCount++
      finalVal = val
    }
	}
	if opCount == 1 {
    a.SetVal(finalVal)
    a.notifyAll()
    return
	}

	for val, op := range a.ops {
    if !op {
      continue
    }
    for _, g := range a.grps {
      if g.JustOneWith(val) {
        a.SetVal(val)
        a.notifyAll()
        return
      }
    }
  }
}

func (a *agent) notifyAll() {
	m := msg{id:a.id, val:a.val}
	a.done <- m
	for _, g := range a.grps {
		for _, friend := range g {
			friend.ch <- m
		}
	}
}

func (a *agent) run() {
  if a.val != 0 {
    a.notifyAll()
  }

	var m msg
	for a.val == 0 {
		select {
			case m = <-a.ch:
				a.markTaken(m)
		}
	}
}

// user defined
func makeGroups(master group) (grps []group) {
  grps = make([]group, 2 * numVals)
  for i, _ := range grps {
    grps[i] = make(group, numVals)
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
func setPuzzle(g group) {
	g[0].SetVal(1)
	g[4].SetVal(2)
	g[8].SetVal(3)
}

func assignGroups(grps []group) {
	for _, g := range grps {
		for _, a := range g {
			a.AddGroup(g)
		}
	}
}

func dispatch(g group) {
	for _, a := range g {
		go a.run()
	}
}
