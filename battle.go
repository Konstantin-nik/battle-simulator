package battle

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func CircleBattle(l []*Battle, cb chan *Battle) {
	for len(l) > 1 {
		// fmt.Println("----------------------------")
		// for _, lb := range l {
		// 	fmt.Println(*lb)
		// }

		var b1, b2 *Battle
		var a Arena
		ch := make(chan *Battle, len(l)/2)
		for len(l) > 1 {
			l, b1, b2 = l[:len(l)-2], l[len(l)-1], l[len(l)-2]
			a = &BattlePair{b1: *b1, b2: *b2}
			wg.Add(1)
			go a.Battle(ch)
		}
		wg.Wait()
		close(ch)
		for val := range ch {
			if *val != nil {
				l = append(l, val)
			}
		}
	}
	// fmt.Println("----------------------------")
	// for _, lb := range l {
	// 	fmt.Println(*lb)
	// }
	if len(l) == 0 {
		cb <- nil
	} else {
		cb <- l[0]
	}
}

// Arena interface used for organising battles with multiple players.
//
// List of acceptable structures:
// 		• type BattlePair struct {}
//
type Arena interface {
	StartBattle()
	UpdateStatus()
	GetResult() (a Battle)
	Battle(bc chan *Battle)
}

type BattlePair struct {
	b1, b2 Battle
	status bool
}

func (bp *BattlePair) StartBattle() {
	bp.UpdateStatus()
	for bp.status {
		bp.b1.DoDamage(bp.b2)
		bp.b2.DoDamage(bp.b1)
		bp.UpdateStatus()
	}
}

func (bp *BattlePair) UpdateStatus() {
	bp.status = bp.b1.IsAlive() && bp.b2.IsAlive()
}

func (bp *BattlePair) GetResult() (a Battle) {
	if bp.b1.IsAlive() {
		return bp.b1
	} else if bp.b2.IsAlive() {
		return bp.b2
	} else {
		return nil
	}
}

func (bp *BattlePair) Battle(bc chan *Battle) {
	// fmt.Println("11")
	bp.StartBattle()
	b := bp.GetResult()
	//time.Sleep(100 * time.Millisecond)
	// fmt.Println("22")
	bc <- &b
	defer wg.Done()
}

// Battle interface used for battle skills.
//
// List of acceptable structures:
// 		• type Warrior struct {}
//
type Battle interface {
	IsAlive() bool
	GetDamage(d float64)
	DoDamage(i interface{ GetDamage(d float64) })
	String() string
}

type Status struct {
	Name  string
	Value []int
}

type Person struct {
	Name   string
	Health float64
	Stat   Status
}

func (p *Person) IsAlive() bool {
	return p.Health > 0
}

func (p *Person) GetDamage(d float64) {
	if d > 0 {
		p.Health -= d
	}
}

func (p *Person) String() string {
	return fmt.Sprintf("Status:\n\tName: %s\n\tHealth: %.2f", p.Name, p.Health)
}

type Warrior struct {
	P               Person
	Damage          float64
	FlatArmor       float64
	Range           float64
	PercentageArmor float64
}

func (w *Warrior) IsAlive() bool {
	return w.P.IsAlive()
}

func (w *Warrior) GetDamage(d float64) {
	w.P.GetDamage(d - w.FlatArmor)
}

func (w *Warrior) DoDamage(i interface{ GetDamage(d float64) }) {
	i.GetDamage((w.Damage + w.Range) * w.PercentageArmor)
}

func (w *Warrior) String() string {
	return w.P.String()
}

// func describe(i interface{}) {
// 	fmt.Printf("(%v, %T)\n", i, i)
// }

// func makePlayer(name string, damage float64) *Battle {
// 	var p Battle = &Warrior{P: Person{Name: name, Health: 200}, Damage: damage,
// 		FlatArmor: 5, Range: 3, PercentageArmor: 0.69}
// 	return &p
// }

// var (
// 	player1 Battle = &Warrior{P: Person{Name: "Hero 1", Health: 200}, Damage: 5,
// 		FlatArmor: 5, Range: 3, PercentageArmor: 0.69}
// 	player2 Battle = &Warrior{P: Person{Name: "Hero 2", Health: 200}, Damage: 6,
// 		FlatArmor: 5, Range: 2.1, PercentageArmor: 0.68}
// 	arena Arena = &BattlePair{b1: player1, b2: player2}
// 	l     []*Battle
// )

// func main() {
// 	v := make(chan *Battle)
// 	l = append(l, &player1)
// 	l = append(l, &player2)
// 	l = append(l, makePlayer("Hero 3", 5))
// 	l = append(l, makePlayer("Hero 4", 5))
// 	l = append(l, makePlayer("Hero 5", 5))
// 	l = append(l, makePlayer("Hero 6", 5))
// 	l = append(l, makePlayer("Hero 7", 7))
// 	go CircleBattle(l, v)
// 	fmt.Println(*<-v)
// 	//fmt.Println("Welcome to the arena!")
// 	//arena.StartBattle()
// 	//fmt.Print(arena.GetResult())
// }
