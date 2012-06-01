
package sudoku

import "fmt"
import "time"
import "errors"

func DefaultPuzzle() *Problem {
  return New(1, 1)
}

type Problem struct {
  masterList Group
  grps []Group
  numVals int
}

func New(numSquares , numVals int) *Problem {
  masterList := make(Group, numSquares)
  for i, _ := range masterList {
    masterList[i] = newSquare(i, numSquares, numVals)
  }
  return &Problem{
    masterList: masterList,
    numVals: numVals,
  }
}

func (p *Problem) Solve(initVals func(Group), makeGroups func(Group) []Group) error {
  initVals(p.masterList)
  p.grps = makeGroups(p.masterList)

	ch := make(tunnel)
  for _, g := range p.grps {
    for _, a := range g {
      a.done = ch
    }
  }
	
	p.assignGroups()
	p.dispatch()

	for count := 0; count < len(p.masterList); {
		select {
			case <-ch:
				count++
      case <-time.After(10 * time.Second):
        return errors.New("sudoku: failed to solve puzzle.")
		}
	}
  return nil
}
func (p *Problem) String() string {
  s := ""
  for _, a := range p.masterList {
    s += fmt.Sprint(a.Val, ",")
  }
  return s
}

type msg struct {
	id int
	val int
}

type tunnel chan msg

type Group []*agent
func (g Group) justOneWith(val int) bool {
  count := 0
  for _, a := range g {
    if a.ops[val] {
      count++
    }
  }
  return count == 1
}

type options map[int]bool

type agent struct {
	id int
	Val int
	ch tunnel
	done tunnel
	grps []Group
	ops options
}

func newSquare(id, numAgents, numVals int) *agent {
	ops := make(options)
	for j := 1; j <= numVals; j++ {
		ops[j] = true
	}
	return &agent{
			id: id,
			grps: []Group{},
			ops: ops,
			ch: make(tunnel, numAgents),
		}
}

func (a *agent) AddGroup(g Group) {
	a.grps = append(a.grps, g)
}

func (a *agent) SetVal(val int) {
  if val == 0 {
    return
  }

	for v, _ := range a.ops {
		a.ops[v] = false
	}
	a.ops[val] = true
	a.Val = val
}

func (a *agent) Run() {
  if a.Val != 0 {
    a.notifyAll()
  }

	var m msg
	for a.Val == 0 {
		select {
			case m = <-a.ch:
				a.markTaken(m)
		}
	}
}

func (a *agent) markTaken(m msg) {
	a.ops[m.val] = false

  // check if single option remains
	opCount := 0
	finalVal := 0
	for val, op := range a.ops {
    if op {
      opCount++
      finalVal = val

      // check if an option is not in any Groups
      for _, g := range a.grps {
        if g.justOneWith(val) {
          a.SetVal(val)
          a.notifyAll()
          return
        }
      }
    }
	}
	if opCount == 1 {
    a.SetVal(finalVal)
    a.notifyAll()
	}
}

func (a *agent) notifyAll() {
	m := msg{id:a.id, val:a.Val}
	a.done <- m
	for _, g := range a.grps {
		for _, friend := range g {
			friend.ch <- m
		}
	}
}

func (p *Problem) assignGroups() {
	for _, g := range p.grps {
		for _, a := range g {
			a.AddGroup(g)
		}
	}
}

func (p *Problem) dispatch() {
	for _, a := range p.masterList {
		go a.Run()
	}
}

