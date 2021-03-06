package battle

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func CircleBattle(l []*Player) *Player {
	ch := make(chan *Player)
	go circleBattle(l, ch)
	return <-ch
}

func circleBattle(l []*Player, cb chan *Player) {
	for len(l) > 1 {
		// fmt.Println("----------------------------")
		// for _, lb := range l {
		// 	fmt.Println(*lb)
		// }

		var b1, b2 *Player
		var a Arena
		ch := make(chan *Player, len(l))
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
	GetResult() (a Player)
	Battle(bc chan *Player)
}

type BattlePair struct {
	b1, b2 Player
	status bool
}

func (bp *BattlePair) StartBattle() {
	bp.UpdateStatus()
	for i := 0; (i < 5) && bp.status; i++ {
		bp.b1.DoDamage(bp.b2)
		bp.b2.DoDamage(bp.b1)
		bp.UpdateStatus()
	}
}

func (bp *BattlePair) UpdateStatus() {
	bp.b1.UpdateStatus()
	bp.b2.UpdateStatus()
	bp.status = bp.b1.IsAlive() && bp.b2.IsAlive()
}

func (bp *BattlePair) GetResult() (a Player) {
	if bp.b1.IsAlive() {
		return bp.b1
	} else if bp.b2.IsAlive() {
		return bp.b2
	} else {
		return nil
	}
}

func (bp *BattlePair) Battle(bc chan *Player) {
	// fmt.Println("11")
	bp.StartBattle()
	if bp.b1.IsAlive() && bp.b2.IsAlive() {
		bc <- &bp.b1
		bc <- &bp.b2
	} else {
		b := bp.GetResult()
		bc <- &b
	}
	defer wg.Done()
}

// Player interface used for battle classes.
//
// List of acceptable structures:
// 		• type Warrior struct {}
//
type Player interface {
	Name() string
	Health() string
	Status() string
	IsAlive() bool
	UpdateStatus()
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
	return fmt.Sprintf("Name: %s\n\tHealth: %s\n\tStatus: %.2f", p.Stat.Name, p.Name, p.Health)
}

type Warrior struct {
	P               Person
	Damage          float64
	FlatArmor       float64
	Range           float64
	PercentageArmor float64
}

func (w *Warrior) Name() string {
	return w.P.Name
}

func (w *Warrior) Health() string {
	return fmt.Sprintf("%.2f", w.P.Health)
}

func (w *Warrior) Status() string {
	return w.P.Stat.Name
}

func (w *Warrior) IsAlive() bool {
	return w.P.IsAlive()
}

func (w *Warrior) UpdateStatus() {
	if w.P.IsAlive() {
		w.P.Stat.Name = "Alive"
	} else {
		w.P.Stat.Name = "Dead"
	}

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

func MakePlayer(name string, health float64, damage float64, flatArmor float64, Range float64, percentageArmor float64) *Player {
	var p Player = &Warrior{P: Person{Name: name, Health: health, Stat: Status{Name: "Alive"}},
		Damage: damage, FlatArmor: flatArmor, Range: Range, PercentageArmor: percentageArmor}
	return &p
}

// func describe(i interface{}) {
// 	fmt.Printf("(%v, %T)\n", i, i)
// }

// var (
// 	player1 Player = &Warrior{P: Person{Name: "Hero 1", Health: 200}, Damage: 5,
// 		FlatArmor: 5, Range: 3, PercentageArmor: 0.69}
// 	player2 Player = &Warrior{P: Person{Name: "Hero 2", Health: 200}, Damage: 6,
// 		FlatArmor: 5, Range: 2.1, PercentageArmor: 0.68}
// 	arena Arena = &BattlePair{b1: player1, b2: player2}
// 	l     []*Player
// )

// func main() {
// 	v := make(chan *Player)
// 	l = append(l, &player1)
// 	l = append(l, &player2)
// 	l = append(l, MakePlayer("Hero 3", 200, 5, 5, 3, 0.69))
// 	l = append(l, MakePlayer("Hero 4", 200, 5, 5, 3, 0.69))
// 	l = append(l, MakePlayer("Hero 5", 200, 5, 5, 3, 0.69))
// 	l = append(l, MakePlayer("Hero 6", 200, 5, 5, 3, 0.69))
// 	l = append(l, MakePlayer("Hero 7", 200, 7, 5, 3, 0.69))
// 	go circleBattle(l, v)
// 	fmt.Println(*<-v)
// 	//fmt.Println("Welcome to the arena!")
// 	//arena.StartBattle()
// 	//fmt.Print(arena.GetResult())
// }
