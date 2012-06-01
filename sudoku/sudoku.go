
package sudoku

import "time"

const (
  NumAgents = 9
	NumVals = 3
)

func Solve(masterList Group, grps []Group) {
	ch := make(tunnel)

  for _, g := range grps {
    for _, a := range g {
      a.done = ch
    }
  }
	
	assignGroups(grps)
	dispatch(masterList)

	for count := 0; count < NumAgents; {
		select {
			case <-ch:
				count++
      case <-time.After(10 * time.Second):
        count = NumAgents
		}
	}
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

func New(id int) *agent {
	ops := make(options)
	for j := 1; j <= NumVals; j++ {
		ops[j] = true
	}
	return &agent{
			id: id,
			grps: []Group{},
			ops: ops,
			ch: make(tunnel, NumAgents),
		}
}

func (a *agent) AddGroup(g Group) {
	a.grps = append(a.grps, g)
}

func (a *agent) SetVal(val int) {
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

func assignGroups(grps []Group) {
	for _, g := range grps {
		for _, a := range g {
			a.AddGroup(g)
		}
	}
}

func dispatch(g Group) {
	for _, a := range g {
		go a.Run()
	}
}

