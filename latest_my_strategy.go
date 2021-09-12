package main

import (
	. "aicup2020/model"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

var (
	white = NewColor(255.0, 255.0, 255.0, 1.0)
	black = NewColor(0.0, 0.0, 0.0, 1.0)
	red   = NewColor(150.0, 0.0, 0.0, 1.0)
	green = NewColor(0.0, 150.0, 0.0, 1.0)
	blue  = NewColor(0.0, 0.0, 150.0, 1.0)
)

type MyStrategy struct {
	stopBuilder    bool
	debugInterface *DebugInterface
	settings       map[EntityType]EntityProperties
	e              [][]Entity
	miner          []Entity
	minerH         map[int32]bool
	melee          []Entity
	ranged         []Entity
	turel          []Entity
	house          []Entity
	builderBase    []Entity
	meleeBase      []Entity
	rangedBase     []Entity
	one            sync.Once
	buildChan      chan EntityType

	squards    [][]Entity
	squardSize int

	position       int
	prevPopulation int
	population     int
	buildAction    []EntityType
	defendSquard   int
	defendPoint    []Entity
	defended       map[int]bool

	stopCallback      func(Action)
	houseVec          *Vec2Int32
	turretVec         []Vec2Int32
	last              int
	builder           []Entity
	lastBuilder       int32
	lastMelee         int32
	lastRanged        int32
	enimies           map[int32][]Entity
	countBuilder      int
	myId              int32
	housePos          *Vec2Int32
	enemyBasePosition *Entity
	currentTick       int32

	center  Vec2Int32
	radius  float32
	build   [][]Vec2Int32
	builded []*Entity

	chanTick chan int32
	size     int
	tick     int32
	res      int
}

func NewMyStrategy() *MyStrategy {
	return &MyStrategy{
		countBuilder: 2,
		position:     0,
		build: [][]Vec2Int32{
			[]Vec2Int32{
				Vec2Int32{5, 11},
				Vec2Int32{4, 11},
				Vec2Int32{4, 12},
				Vec2Int32{4, 13},
			},
			[]Vec2Int32{
				Vec2Int32{11, 11},
				Vec2Int32{10, 11},
				Vec2Int32{10, 12},
				Vec2Int32{10, 13},
			},
			[]Vec2Int32{
				Vec2Int32{11, 5},
				Vec2Int32{10, 5},
				Vec2Int32{10, 6},
				Vec2Int32{10, 7},
			},
			[]Vec2Int32{
				Vec2Int32{11, 1},
				Vec2Int32{10, 1},
				Vec2Int32{10, 2},
				Vec2Int32{10, 3},
			},
			[]Vec2Int32{
				Vec2Int32{15, 1},
				Vec2Int32{14, 1},
				Vec2Int32{14, 2},
				Vec2Int32{14, 3},
			},
			[]Vec2Int32{
				Vec2Int32{7, 1},
				Vec2Int32{6, 1},
				Vec2Int32{6, 2},
				Vec2Int32{6, 3},
			},
			[]Vec2Int32{
				Vec2Int32{3, 1},
				Vec2Int32{2, 1},
				Vec2Int32{2, 2},
				Vec2Int32{2, 3},
			},
			[]Vec2Int32{
				Vec2Int32{1, 5},
				Vec2Int32{0, 5},
				Vec2Int32{0, 6},
				Vec2Int32{0, 7},
			},
			[]Vec2Int32{
				Vec2Int32{1, 9},
				Vec2Int32{0, 9},
				Vec2Int32{0, 10},
				Vec2Int32{0, 11},
			},
			[]Vec2Int32{
				Vec2Int32{1, 13},
				Vec2Int32{0, 14},
				Vec2Int32{0, 15},
				Vec2Int32{0, 16},
			},
			[]Vec2Int32{
				Vec2Int32{1, 17},
				Vec2Int32{0, 18},
				Vec2Int32{0, 19},
				Vec2Int32{0, 20},
			},
			[]Vec2Int32{
				Vec2Int32{1, 21},
				Vec2Int32{0, 22},
				Vec2Int32{0, 23},
				Vec2Int32{0, 24},
			},
			[]Vec2Int32{
				Vec2Int32{5, 21},
				Vec2Int32{4, 22},
				Vec2Int32{4, 23},
				Vec2Int32{4, 24},
			},
			[]Vec2Int32{
				Vec2Int32{9, 21},
				Vec2Int32{8, 22},
				Vec2Int32{8, 23},
				Vec2Int32{8, 24},
			},
			[]Vec2Int32{
				Vec2Int32{13, 21},
				Vec2Int32{12, 22},
				Vec2Int32{12, 23},
				Vec2Int32{12, 24},
			},
			[]Vec2Int32{
				Vec2Int32{17, 21},
				Vec2Int32{16, 22},
				Vec2Int32{16, 23},
				Vec2Int32{16, 24},
			},
			[]Vec2Int32{
				Vec2Int32{21, 21},
				Vec2Int32{20, 22},
				Vec2Int32{20, 23},
				Vec2Int32{20, 24},
			},
			[]Vec2Int32{
				Vec2Int32{21, 17},
				Vec2Int32{20, 18},
				Vec2Int32{20, 19},
				Vec2Int32{20, 20},
			},
			[]Vec2Int32{
				Vec2Int32{21, 13},
				Vec2Int32{20, 14},
				Vec2Int32{20, 15},
				Vec2Int32{20, 16},
			},
			[]Vec2Int32{
				Vec2Int32{21, 9},
				Vec2Int32{20, 10},
				Vec2Int32{20, 11},
				Vec2Int32{20, 12},
			},
			[]Vec2Int32{
				Vec2Int32{21, 5},
				Vec2Int32{20, 6},
				Vec2Int32{20, 7},
				Vec2Int32{20, 8},
			},
			[]Vec2Int32{
				Vec2Int32{21, 1},
				Vec2Int32{20, 2},
				Vec2Int32{20, 3},
				Vec2Int32{20, 4},
			},
		},
		buildAction: []EntityType{
			EntityTypeBuilderUnit, EntityTypeRangedUnit, EntityTypeMeleeUnit,
		},
		turretVec: []Vec2Int32{
			Vec2Int32{15, 15},
		},
		lastMelee:  int32(1),
		lastRanged: int32(1),
		squardSize: int(20),
		defended: map[int]bool{
			0: false,
			1: false,
			2: false,
		},
		defendPoint: []Entity{},
		buildChan:   make(chan EntityType),
		chanTick:    make(chan int32),
		radius:      float32(10),
		center:      Vec2Int32{10, 10},
	}
}

func (s *MyStrategy) NextBuild() EntityType {
	if len(s.miner) < 5 {
		return s.buildAction[0]
	}

	p := s.position
	if p >= len(s.buildAction) {
		p = 0
		s.position = 0
	}
	return s.buildAction[p]
}

func (s *MyStrategy) Population() int {
	result := 0
	for _, m := range s.miner {
		result += int(s.settings[m.EntityType].PopulationUse)
	}
	for _, m := range s.melee {
		result += int(s.settings[m.EntityType].PopulationUse)
	}
	for _, m := range s.ranged {
		result += int(s.settings[m.EntityType].PopulationUse)
	}
	for _, m := range s.builder {
		result += int(s.settings[m.EntityType].PopulationUse)
	}
	for _, ss := range s.squards {
		for _, m := range ss {
			result += int(s.settings[m.EntityType].PopulationUse)
		}
	}

	return result
}

func (s *MyStrategy) Capacity() int {
	result := 0
	for _, h := range s.house {
		if h.Active {
			result += int(s.settings[EntityTypeHouse].PopulationProvide)
		}
	}
	for _, h := range s.builderBase {
		if h.Active {
			result += int(s.settings[EntityTypeBuilderBase].PopulationProvide)
		}
	}
	for _, h := range s.meleeBase {
		if h.Active {
			result += int(s.settings[EntityTypeMeleeBase].PopulationProvide)
		}
	}
	for _, h := range s.rangedBase {
		if h.Active {
			result += int(s.settings[EntityTypeRangedBase].PopulationProvide)
		}
	}
	return result
}

func (s *MyStrategy) FindResource(e Vec2Int32) (Vec2Int32, EntityType, float64) {
	x := int32(e.X)
	y := int32(e.Y)
	distance := int32(70)

	dx := x + distance
	dy := y + distance

	if dx > int32(s.size) {
		dx = int32(s.size) - int32(1)
	}
	if dy > int32(s.size) {
		dy = int32(s.size) - int32(1)
	}

	if x > distance {
		x = x - distance
	}
	if y > distance {
		y = y - distance
	}

	min := float64(s.size)
	var fe Entity

	for i := x; i < dx; i++ {
		for j := y; j < dy; j++ {
			pe := s.e[i][j]
			if pe.EntityType != EntityTypeResource {
				continue
			}
			if d := distantion(e, pe.Position); d < min {
				min = d
				fe = pe
			}
		}
	}
	// s.PrintRectange(Vec2Int32{x, y}, Vec2Int32{dx, dy}, red)
	return fe.Position, fe.EntityType, min

}

func (s *MyStrategy) Mining(units []Entity, a Action) {
	center := ByID(units).Center()
	pos, _, _ := s.FindResource(center)
	for _, e := range units {
		var ma MoveAction
		pose, t, d := s.FindEnimy2(e.Position)
		if t.EntityType == EntityTypeMeleeUnit && d < 3 {
			s.DoMoveFrom(e, pose, a)
			continue
		}
		if t.EntityType == EntityTypeTurret && d < 7 {
			s.DoMoveFrom(e, pose, a)
			continue
		}
		if t.EntityType == EntityTypeRangedUnit && d < 6 {
			s.DoMoveFrom(e, pose, a)
			continue
		}
		// if (e.Position.X < 13 || e.Position.Y < 13) && s.tick > 400 {
		// ma = NewMoveAction(Vec2Int32{16, 19}, true, true)
		// if s.tick > 800 {
		//
		// ma = NewMoveAction(Vec2Int32{40, 40}, true, true)
		// }
		// }

		ma = NewMoveAction(Vec2Int32{pos.X, pos.Y}, true, true)
		aaa := NewAutoAttack(int32(s.size), []EntityType{EntityTypeResource})
		aa := NewAttackAction(nil, &aaa)
		ea, ok := a.EntityActions[e.Id]
		if !ok {
			a.EntityActions[e.Id] = EntityAction{
				MoveAction:   &ma,
				AttackAction: &aa,
			}

		} else {
			ea.MoveAction = &ma
			ea.AttackAction = &aa
			a.EntityActions[e.Id] = ea
		}
	}
}

func (s *MyStrategy) BuilderPositionHouse() Vec2Int32 {
	v := s.builderBase[0].Position
	size := s.settings[EntityTypeHouse].Size
	v.X = v.X - size - 1
	v.Y = v.Y + size*int32(len(s.house))
	return v
}

func (s *MyStrategy) NextHousePosition() Vec2Int32 {
	v := s.builderBase[0].Position
	size := s.settings[EntityTypeHouse].Size
	v.X = v.X - size - 1
	v.Y = v.Y + size*int32(len(s.house)) + 1
	return v
}

func (s *MyStrategy) FindSpaceForTurret(pos Vec2Int32) Vec2Int32 {
	x := pos.X - 3
	y := pos.Y - 3
	if y < 1 || x < 1 {
		return Vec2Int32{}
	}

	v := Vec2Int32{x, y}
	// if !s.IsFree(v, EntityTypeTurret) {
	// return Vec2Int32{}
	// }
	return v
}

func (s *MyStrategy) FindSpaceForMeleeBase() Vec2Int32 {
	if len(s.builderBase) < 0 {

		return Vec2Int32{11, 11}
	}
	size := s.settings[EntityTypeBuilderBase].Size
	p := s.builderBase[0].Position
	return Vec2Int32{p.X + size, p.Y + size}
}

func (s *MyStrategy) FindSpaceForRangedBase() Vec2Int32 {
	return Vec2Int32{26, 5}
	// if s.IsFree(Vec2Int32{18, 11}, EntityTypeRangedBase) {
	// return Vec2Int32{18, 11}
	// }
	// if s.IsFree(Vec2Int32{11, 18}, EntityTypeRangedBase) {
	// return Vec2Int32{11, 18}
	// } else {
	// return Vec2Int32{18, 18}
	// }
	// size := s.settings[EntityTypeBuilderBase].Size
	// p := s.builderBase[0].Position
	// return Vec2Int32{p.X + size + 8, p.Y + size + 8}
}

func (s *MyStrategy) FindSpaceForHouse(offset int) []Vec2Int32 {
	pos := len(s.house) + offset
	if pos >= len(s.build) {
		v1 := s.roundFunc(EntityTypeHouse)
		return []Vec2Int32{
			v1,
			Vec2Int32{v1.X - 1, v1.Y + 1},
			Vec2Int32{v1.X - 1, v1.Y + 2},
			Vec2Int32{v1.X - 1, v1.Y + 3},
		}
	}

	return s.build[pos]
	// if len(s.house) == 0 {
	// return Vec2Int32{5, 11}
	// } else if len(s.house) == 1 {
	// return Vec2Int32{8, 11}
	// } else if (len(s.house)) == 2 {
	// return Vec2Int32{11, 11}
	// } else if (len(s.house)) == 3 {
	// return Vec2Int32{11, 4}
	// }
}

func (s *MyStrategy) BuildTurret(a Action) {
	if len(s.miner) < 20 {
		return
	}

	for _, h := range s.turel {
		units := s.miner[:2]
		s.miner = s.miner[2:]
		if !h.Active || h.Health < s.settings[EntityTypeTurret].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}

	cost := s.settings[EntityTypeTurret].InitialCost
	if int32(s.res) < 2*cost+int32(len(s.turel)) {
		fmt.Println(">>>cost")
		return
	}

	var (
		p Vec2Int32
		e Entity
	)
	for i := 0; i < len(s.miner); i++ {
		e = s.miner[i]

		p = s.FindSpaceForTurret(e.Position)
		if myVector(p).IsZero() {
			continue
		}
		if p.X-1 < 0 {
			continue
		}
		s.miner = append(s.miner[:i], s.miner[i+1:]...)
		break
	}
	if myVector(p).IsZero() {
		return
	}

	tgt := Vec2Int32{p.X - 1, p.Y}
	ma := NewMoveAction(
		tgt, true, true,
	)
	ba := NewBuildAction(
		EntityTypeTurret,
		p,
	)
	a.EntityActions[e.Id] = EntityAction{
		MoveAction:  &ma,
		BuildAction: &ba,
	}
}

func (s *MyStrategy) BuildMeleeBase(c int32, units []Entity, a Action) {
	if len(s.meleeBase) > 2 || len(s.house) < 7 {
		return
	}

	for _, h := range s.meleeBase {
		if !h.Active || h.Health < s.settings[EntityTypeMeleeBase].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}

	for i, e := range units {

		var p Vec2Int32
		p = s.FindSpaceForMeleeBase()
		v := p

		if len(s.house)%2 != 0 {
			v.X = v.X - int32(1)
			v.Y = v.Y + int32(i)

		} else {
			v.X = v.X + int32(i)
			v.Y = v.Y - int32(1)
		}
		ma := NewMoveAction(
			v, true, true,
		)
		ba := NewBuildAction(
			EntityTypeMeleeBase,
			p,
		)
		ea, ok := a.EntityActions[e.Id]
		if !ok {
			a.EntityActions[e.Id] = EntityAction{
				MoveAction:  &ma,
				BuildAction: &ba,
			}
		} else {
			ea.MoveAction = &ma
			ea.BuildAction = &ba
			a.EntityActions[e.Id] = ea
		}
	}
}

func (s *MyStrategy) BuildRangedBase(c int32, units []Entity, a Action) {
	t := EntityTypeRangedUnit
	if len(s.rangedBase) > 2 || len(s.house) < 7 {
		return
	}

	for _, h := range s.rangedBase {
		if !h.Active || h.Health < s.settings[t].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}

	for i, e := range units {

		var p Vec2Int32
		p = s.FindSpaceForRangedBase()
		v := p

		if len(s.house)%2 != 0 {
			v.X = v.X - int32(1)
			v.Y = v.Y + int32(i)

		} else {
			v.X = v.X + int32(i)
			v.Y = v.Y - int32(1)
		}
		ma := NewMoveAction(
			v, true, true,
		)
		ba := NewBuildAction(
			EntityTypeRangedBase,
			p,
		)
		ea, ok := a.EntityActions[e.Id]
		if !ok {

			a.EntityActions[e.Id] = EntityAction{
				MoveAction:  &ma,
				BuildAction: &ba,
			}
		} else {
			ea.MoveAction = &ma
			ea.BuildAction = &ba
			a.EntityActions[e.Id] = ea
		}

	}
}

func (s *MyStrategy) BuildHouse(c int32, units []Entity, a Action) {
	for _, h := range s.house {
		if !h.Active || h.Health < s.settings[EntityTypeHouse].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}
	for _, h := range s.builderBase {
		if !h.Active || h.Health < s.settings[EntityTypeBuilderBase].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}
	for _, h := range s.meleeBase {
		if !h.Active || h.Health < s.settings[EntityTypeMeleeBase].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}
	for _, h := range s.rangedBase {
		if !h.Active || h.Health < s.settings[EntityTypeRangedBase].MaxHealth {
			s.Repair(units, a, h)
			return
		}
	}
	p := s.FindSpaceForHouse(0)
	if myVector(p[0]).IsZero() {
		return
	}

	// size := s.settings[EntityTypeHouse].Size
	// h := p[0]

	sort.Sort(ByID(units))
	// build := false
	for i, e := range units {
		ea := EntityAction{}

		ba := NewBuildAction(
			EntityTypeHouse,
			p[0],
		)
		ea.BuildAction = &ba

		buildPosition := p[i+1]
		if !myVector(e.Position).Equal(myVector(buildPosition)) {

			ma := NewMoveAction(
				buildPosition, true, true,
			)
			ea.MoveAction = &ma
		}

		a.EntityActions[e.Id] = ea
	}
}

func (s *MyStrategy) Repair(units []Entity, a Action, tgt Entity) {
	for _, e := range units {
		ma := NewMoveAction(
			Vec2Int32{X: tgt.Position.X, Y: tgt.Position.Y},
			true,
			true,
		)
		ra := NewRepairAction(tgt.Id)
		a.EntityActions[e.Id] = EntityAction{MoveAction: &ma, RepairAction: &ra}
	}
}

func (s *MyStrategy) incCurrentTick(c int32) {
	s.currentTick = c + 2
}

func (s *MyStrategy) MakeRanged(c int32, a Action) {
	t := EntityTypeRangedBase
	size := s.settings[t].Size
	// cost := s.settings[EntityTypeRangedUnit].InitialCost
	// pop := s.settings[EntityTypeRangedUnit].PopulationUse
	// if cost > int32(s.res) || s.Population()+int(pop) >= s.Capacity() {
	// return
	// }
	s.incCurrentTick(c)

	for _, e := range s.rangedBase {
		pos := Vec2Int32{X: e.Position.X + size, Y: e.Position.Y + size - s.lastRanged}
		ba := NewBuildAction(
			EntityTypeRangedUnit,
			// добавить смещение координаты при построении
			pos,
		)
		a.EntityActions[e.Id] = EntityAction{BuildAction: &ba}
		s.lastRanged += int32(1)
		if s.lastRanged > size {
			s.lastRanged = int32(1)
		}
	}
}

func (s *MyStrategy) StopMakeRanged(a Action) {
	for _, e := range s.rangedBase {
		a.EntityActions[e.Id] = EntityAction{}
	}
}

func (s *MyStrategy) MakeMelee(c int32, a Action) {
	t := EntityTypeMeleeBase
	size := s.settings[t].Size
	// cost := s.settings[EntityTypeMeleeUnit].InitialCost
	// pop := s.settings[EntityTypeMeleeUnit].PopulationUse
	if len(s.melee) > len(s.rangedBase)/2 {
		return
	}
	s.incCurrentTick(c)

	for _, e := range s.meleeBase {
		pos := Vec2Int32{X: e.Position.X + size, Y: e.Position.Y + size - s.lastMelee}
		ba := NewBuildAction(EntityTypeMeleeUnit, pos)
		a.EntityActions[e.Id] = EntityAction{BuildAction: &ba}
		s.lastMelee += int32(1)
		if s.lastMelee > size {
			s.lastMelee = int32(1)
		}
	}
}

func (s *MyStrategy) StopMakeMelee(a Action) {
	for _, e := range s.meleeBase {
		a.EntityActions[e.Id] = EntityAction{}
	}
}

func (s *MyStrategy) MakeBuilder(c int32, a Action) {
	t := EntityTypeBuilderUnit
	unit := s.settings[t].Size
	// cost := s.settings[t].InitialCost
	// pop := s.settings[t].PopulationUse

	size := s.settings[EntityTypeBuilderBase].Size
	if s.stopBuilder {
		for _, e := range s.builderBase {
			a.EntityActions[e.Id] = EntityAction{}
		}
		return

	}
	s.incCurrentTick(c)

	for _, e := range s.builderBase {
		ba := NewBuildAction(
			EntityTypeBuilderUnit,
			Vec2Int32{X: e.Position.X + size, Y: e.Position.Y + size - s.lastBuilder},
		)
		a.EntityActions[e.Id] = EntityAction{BuildAction: &ba}
		s.lastBuilder += unit
		if s.lastBuilder > 5*unit {
			s.lastBuilder = 0
		}
	}
}

func (s *MyStrategy) StopMakeBuilder(a Action) {
	for _, e := range s.builderBase {
		a.EntityActions[e.Id] = EntityAction{}
	}
}

func distantion(v1 Vec2Int32, v2 Vec2Int32) float64 {
	return math.Sqrt(math.Pow(float64(v2.X-v1.X), float64(2)) + math.Pow(float64(v2.Y-v1.Y), float64(2)))
}

func PointOnVector(v1 Vec2Int32, v2 Vec2Int32, length float32) Vec2Int32 {
	dx := float64(v2.X - v1.X)
	dy := float64(v2.Y - v1.Y)
	r := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	k := float64(length) / r
	x := v1.X + int32(dx*k)
	y := v1.Y + int32(dy*k)
	return Vec2Int32{x, y}
}

func (s *MyStrategy) FindEnimyBase(en Vec2Int32) (Vec2Int32, Entity, float64) {
	if s.enemyBasePosition != nil {
		return s.enemyBasePosition.Position, *s.enemyBasePosition, distantion(en, s.enemyBasePosition.Position)
	}
	min := float64(s.size)
	var fe Entity
	for id, e := range s.enimies {
		if id == s.myId {
			continue
		}
		sort.Sort(ByID(e))
		for _, ee := range e {
			switch ee.EntityType {
			case EntityTypeHouse, EntityTypeRangedBase, EntityTypeMeleeBase, EntityTypeTurret:
				if d := distantion(en, ee.Position); d < min {
					min = d
					fe = ee
				}
			}
		}
	}
	s.enemyBasePosition = &fe
	return fe.Position, fe, min
}

// func (s *MyStrategy) FindEnimyVector() ([][]Entity) {
// x := int32(e.X)
// y := int32(e.Y)
// distance := int32(20)
//
// dx := x + distance
// dy := y + distance
//
// if dx > int32(s.size) {
// dx = int32(s.size) - int32(1)
// }
// if dy > int32(s.size) {
// dy = int32(s.size) - int32(1)
// }
//
// if x > distance {
// x = x - distance
// }
// if y > distance {
// y = y - distance
// }
//
// min := float64(s.size)
// var fe Entity
//
// for i := x; i < dx; i++ {
// for j := y; j < dy; j++ {
// pe := s.e[i][j]
// pid := pe.PlayerId
// if pid != nil && *pid != s.myId {
// if d := distantion(e, pe.Position); d < min {
// min = d
// fe = pe
// }
// }
// }
// }
// s.PrintRectange(Vec2Int32{x, y}, Vec2Int32{dx, dy}, red)
// return fe.Position, fe, min
// }

func (s *MyStrategy) FindEnimy2(e Vec2Int32) (Vec2Int32, Entity, float64) {
	x := int32(e.X)
	y := int32(e.Y)
	distance := int32(25)
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > int32(s.size) {
		x = int32(s.size)
	}
	if y > int32(s.size) {
		y = int32(s.size)
	}

	dx := x + distance
	dy := y + distance

	if dx > int32(s.size) {
		dx = int32(s.size) - int32(1)
	}
	if dy > int32(s.size) {
		dy = int32(s.size) - int32(1)
	}

	if x > distance {
		x = x - distance
	}
	if y > distance {
		y = y - distance
	}

	min := float64(s.size)
	var fe Entity

	for i := x; i < dx; i++ {
		for j := y; j < dy; j++ {
			pe := s.e[i][j]
			pid := pe.PlayerId
			if pid != nil && *pid != s.myId {
				if d := distantion(e, pe.Position); d < min {
					min = d
					fe = pe
				}
			}
		}
	}
	// s.PrintRectange(Vec2Int32{x, y}, Vec2Int32{dx, dy}, red)
	return fe.Position, fe, min
}

func (s *MyStrategy) FindEnimy(t EntityType) (*Vec2Int32, EntityType) {
	return nil, EntityTypeRangedUnit
}

func (s *MyStrategy) AttackRage(e Entity, a Action, p *Vec2Int32, et EntityType) {
	return
}

func (s *MyStrategy) RangedMove(e Entity, a Action) {
	return
}

func (s *MyStrategy) RnagedAction(a Action) {
	t := EntityTypeRangedUnit
	for _, e := range s.ranged {
		pos, etype := s.FindEnimy(t)
		if pos != nil {
			s.AttackRage(e, a, pos, etype)
			continue
		}
		s.RangedMove(e, a)
	}
}

func (s *MyStrategy) DoDefendTurret(a Action) {
	for _, e := range s.turel {
		aaa := NewAutoAttack(20, []EntityType{EntityTypeMeleeUnit, EntityTypeRangedUnit, EntityTypeRangedBase, EntityTypeHouse, EntityTypeMeleeBase, EntityTypeTurret, EntityTypeRangedBase, EntityTypeBuilderBase, EntityTypeBuilderUnit})
		aa := NewAttackAction(nil, &aaa)
		ea := EntityAction{
			AttackAction: &aa,
		}
		ea, ok := a.EntityActions[e.Id]
		if !ok {

			ea := EntityAction{
				AttackAction: &aa,
			}
			a.EntityActions[e.Id] = ea
		} else {
			ea.AttackAction = &aa
			a.EntityActions[e.Id] = ea

		}

	}
}

func (s *MyStrategy) DoAttack(e Entity, et Entity, pos Vec2Int32, a Action) {
	// ma := NewMoveAction(
	// Vec2Int32{X: pos.X, Y: pos.Y},
	// true,
	// true,
	// )

	aaa := NewAutoAttack(2, []EntityType{EntityTypeBuilderUnit, EntityTypeRangedUnit, EntityTypeMeleeUnit, EntityTypeRangedBase, EntityTypeHouse, EntityTypeMeleeBase, EntityTypeTurret, EntityTypeRangedBase, EntityTypeBuilderBase})
	// aa := NewAttackAction(nil, &aaa)
	aa := NewAttackAction(&et.Id, &aaa)
	ea := EntityAction{
		// MoveAction:   &ma,
		AttackAction: &aa,
	}
	ea, ok := a.EntityActions[e.Id]
	if !ok {

		ea := EntityAction{
			// MoveAction:   &ma,
			AttackAction: &aa,
		}
		a.EntityActions[e.Id] = ea
	} else {
		// ea.MoveAction = &ma
		ea.AttackAction = &aa
		a.EntityActions[e.Id] = ea

	}
}

func (s *MyStrategy) DoMoveTo(e Entity, pos Vec2Int32, a Action) {
	s.DoMove(e, pos, a, 1)
}

func (s *MyStrategy) DoMoveFrom(e Entity, pos Vec2Int32, a Action) {
	s.DoMove(e, pos, a, -1)
}

func (s *MyStrategy) DoMove(e Entity, pos Vec2Int32, a Action, direction int32) {
	x := e.Position.X + direction*(pos.X-e.Position.X)
	y := e.Position.Y + direction*(pos.Y-e.Position.Y)

	if x < int32(0) {
		x = int32(0)
	} else if x > int32(s.size) {
		x = int32(s.size)
	}

	if y < int32(0) {
		y = int32(0)
	} else if y > int32(s.size) {
		y = int32(s.size)
	}
	ma := NewMoveAction(
		Vec2Int32{X: x, Y: y},
		true,
		true,
	)
	ea, ok := a.EntityActions[e.Id]
	if !ok {
		a.EntityActions[e.Id] = EntityAction{
			MoveAction: &ma,
		}
	} else {
		ea.MoveAction = &ma
		a.EntityActions[e.Id] = ea
	}

}

func (a ByID) Shape(c Vec2Int32, to Vec2Int32, s *MyStrategy, ac Action, angle float64, direction int32) {

	cos := math.Cos(angle * math.Pi / 180)
	sin := math.Sin(angle * math.Pi / 180)

	for j := int32(0); j < int32(len(a)); j++ {
		l := len(a)
		i := int32(l/2) - j
		p := PointOnVector(to, c, float32(i))

		x := int32(float64(c.X) + cos*float64(c.X+p.X-c.X) - sin*float64(c.Y+p.Y-c.Y))
		y := int32(float64(c.Y) + sin*float64(c.X+p.X-c.X) + cos*float64(c.Y+p.Y-c.Y))

		s.DoMoveTo(a[j], Vec2Int32{x, y}, ac)
	}
}

type myVector Vec2Int32

func (v1 myVector) Equal(v2 myVector) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func (v1 myVector) IsZero() bool {
	return v1.X == int32(0) && v1.Y == int32(0)
}

func (s *MyStrategy) DoDefend(a Action) {
	s.defendSquard = -1
	for i, e := range s.defendPoint {
		center := ByID([]Entity{e}).Center()
		pos, t, d := s.FindEnimy2(center)
		skip := false
		if myVector(pos).IsZero() {
			skip = true
		}
		s.PrintVector(center, pos, black)
		s.PrintText(fmt.Sprintf("%d", i), center, black)
		if skip {
			ByID(s.squards[i]).Shape(center, pos, s, a, 90.0, 1)
			s.defended[i] = false
			continue
		}

		if t.EntityType != EntityTypeMeleeUnit && t.EntityType != EntityTypeRangedUnit {
			return
		}
		if d < 15 {
			d := float64(s.size)
			for i, ss := range s.squards {
				if len(ss) == 0 {
					continue
				}
				centerSS := ByID(ss).Center()
				s.PrintText(fmt.Sprintf("%d", i), centerSS, blue)
				d2 := distantion(center, centerSS)
				// ByID(s.squards[i]).Shape(centerSS, pos, s, a, 90.0, 1)
				if d2 < d {
					d = d2
					s.defended[i] = true
				}
			}
			for _, e := range s.squards[i] {
				pos, t, d := s.FindEnimy2(e.Position)
				s.PrintVector(e.Position, pos, blue)
				if e.EntityType == EntityTypeRangedUnit && t.EntityType == EntityTypeMeleeUnit && int32(d) < s.settings[EntityTypeMeleeUnit].Attack.AttackRange+1 {
					s.DoMoveFrom(e, pos, a)
				} else {
					s.DoMoveTo(e, pos, a)
					s.DoAttack(e, t, e.Position, a)
				}
			}
		}
		s.defended[i] = false
	}
}

func (s *MyStrategy) DoAction(a Action) {

	if pid := s.enemyBasePosition; pid != nil && pid.PlayerId != nil {
		if en, ok := s.enimies[*pid.PlayerId]; ok {
			found := false
			for _, ee := range en {
				if ee.Id == s.enemyBasePosition.Id {
					found = true
					break
				}
			}
			if !found {
				s.enemyBasePosition = nil
			}

		} else {
			s.enemyBasePosition = nil
		}
	}

	do := func(d float64, e Entity, pos Vec2Int32, a Action, t Entity, min Entity) {
		switch e.EntityType {
		case EntityTypeBuilderUnit:
			if d < 5 {
				s.DoMoveFrom(e, pos, a)
			}
			if min.Id != int32(0) && min.PlayerId != nil && *min.PlayerId == s.myId {
				s.Repair([]Entity{e}, a, min)
			} else {
				s.DoMoveTo(e, pos, a)
			}
		case EntityTypeMeleeUnit:
			s.DoAttack(e, t, pos, a)
			s.DoMoveTo(e, pos, a)
		case EntityTypeRangedUnit:
			s.DoAttack(e, t, pos, a)
			if t.EntityType == EntityTypeMeleeUnit && d < 3 {
				s.DoMoveFrom(e, pos, a)
			} else {
				s.DoAttack(e, t, pos, a)
				s.DoMoveTo(e, pos, a)
			}
		}
	}

	for i, ss := range s.squards {
		if defend := s.defended[i]; defend {
			continue
		}
		min := Entity{}

		if len(ss) == 0 {
			continue
		}
		center := s.defendPoint[i].Position
		// center := ByID(ss).MedianCenter()
		pos, t, d := s.FindEnimy2(center)
		if myVector(pos).IsZero() {
			pos = s.defendPoint[i].Position
		}

		ld := ByID(ss).SumDamage(s)
		if d < 20 {
			for _, e := range ss {
				pos, t, d := s.FindEnimy2(e.Position)
				do(d, e, pos, a, t, min)
			}
			continue
		}

		if ld > 45 {
			pos, _, _ := s.FindEnimyBase(center)
			if myVector(pos).Equal(myVector(center)) {
				rx := int32(rand.Intn(5))
				ry := int32(rand.Intn(5))
				pos.X = pos.X + rx
				pos.Y += pos.Y + ry
			}
			for _, e := range ss {
				pos, t, d := s.FindEnimy2(e.Position)
				do(d, e, pos, a, t, min)
			}
			continue
		} else {
			s.enemyBasePosition = nil
		}

		// if myVector(pos).Equal(myVector(center)) {
		// rx := int32(rand.Intn(5))
		// ry := int32(rand.Intn(5))
		// pos.X = pos.X + rx
		// pos.Y = pos.Y + ry
		// }

		for _, e := range ss {
			if min.Id == int32(0) {
				min = e
			} else if min.Health > e.Health {
				min = e
			}

			do(d, e, pos, a, t, min)
		}
	}
}

func (s *MyStrategy) MoveFrom(pos Entity, e Vec2Int32) bool {
	size := s.settings[pos.EntityType].Size + 1
	p := pos.Position
	if (p.X < 1 || p.X > int32(len(s.e)-1)) || (p.Y < 1 || p.Y > int32(len(s.e)-1)) {
		return false
	}

	start := s.e[p.X][p.Y].Position
	s.PrintSpace(start, "F", green, EntityTypeHouse)
	return (p.X >= e.X && p.X+size <= e.X) || (p.Y >= e.Y && p.Y+size <= e.Y)
}

func (s *MyStrategy) IsFree(pos Vec2Int32, t EntityType) bool {
	size := s.settings[t].Size + 1
	free := true
	if (pos.X < 1 || pos.X > int32(len(s.e)-1)) || (pos.Y < 1 || pos.Y > int32(len(s.e)-1)) {
		return false
	}

	start := s.e[pos.X][pos.Y].Position
	s.PrintSpace(start, "F", green, EntityTypeHouse)
	for i := pos.X; i < pos.X+size; i++ {
		for j := pos.Y; j < pos.Y+size; j++ {
			entity := s.e[i][j]
			free = (free && entity.Id == int32(0))
			if !free {
				switch entity.EntityType {
				case EntityTypeBuilderUnit:
					continue
				case EntityTypeMeleeUnit:
					continue
				case EntityTypeRangedUnit:
					continue
				}

				break
			}
		}
	}

	s.e[pos.X][pos.Y] = Entity{Id: -1}

	return free
}

func (s *MyStrategy) roundFunc(t EntityType) (v Vec2Int32) {
	if s.houseVec != nil {
		return *s.houseVec
	}
	size := s.settings[t].Size
	count := int32(len(s.house))
	v = Vec2Int32{count, count}
	for i := int32(0); i < int32(s.size)-size-1; i++ {
		for k := int32(0); k <= i; k++ {
			v = Vec2Int32{i - k, k}
			if s.IsFree(v, t) {
				s.houseVec = &v
				return v
			}
		}
	}
	return v
}

func (s *MyStrategy) run() {
	s.buildChan <- s.NextBuild()
	go func() {
		for {
			select {
			case _, ok := <-s.chanTick:
				if !ok {
					return
				}
			}
			e := s.NextBuild()
			if s.res < int(s.settings[e].InitialCost) {
				continue
			}
			if s.prevPopulation != s.Population() {
				s.buildChan <- e
				s.position++
			}
			// s.buildChan
		}
	}()
}

func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
	go func() {
		s.chanTick <- s.currentTick
	}()
	s.debugInterface = debugInterface
	if s.settings == nil {
		s.settings = playerView.EntityProperties
	}
	go s.one.Do(s.run)
	s.squards = [][]Entity{}
	s.size = int(playerView.MapSize)

	s.tick = playerView.CurrentTick
	s.prevPopulation = s.population
	s.e = make([][]Entity, int(playerView.MapSize))
	for i := range s.e {
		s.e[i] = make([]Entity, int(playerView.MapSize))
	}
	s.miner = []Entity{}
	s.builder = []Entity{}
	s.melee = []Entity{}
	s.ranged = []Entity{}
	s.turel = []Entity{}
	s.house = []Entity{}
	s.builderBase = []Entity{}
	s.meleeBase = []Entity{}
	s.rangedBase = []Entity{}
	s.enimies = map[int32][]Entity{}
	s.defendPoint = []Entity{}

	a := Action{
		EntityActions: make(map[int32]EntityAction),
	}
	s.myId = playerView.MyId

	for _, p := range playerView.Players {
		if p.Id == s.myId {
			s.res = int(p.Resource)
		}
	}

	for _, e := range playerView.Entities {
		t := e.EntityType
		size := s.settings[t].Size
		shift := int32(1)
		if t == EntityTypeHouse {
			shift = int32(0)
		}
		size = size + shift
		x := e.Position.X
		y := e.Position.Y

		for i := x; i <= x+size && i < int32(s.size); i++ {
			for j := y; j <= y+size && j < int32(s.size); j++ {
				s.e[i][j] = e
			}
		}
		if e.PlayerId == nil {
			continue
		}
		if *e.PlayerId != playerView.MyId {
			// if e.EntityType == EntityTypeBuilderBase {
			c, ok := s.enimies[*e.PlayerId]
			if ok {
				c = append(c, e)
				s.enimies[*e.PlayerId] = c
				continue
			}
			s.enimies[*e.PlayerId] = []Entity{e}
			// }
			continue
		}

		switch e.EntityType {

		case EntityTypeWall:
		case EntityTypeHouse:
			s.house = append(s.house, e)
		case EntityTypeBuilderBase:
			s.builderBase = append(s.builderBase, e)
		case EntityTypeMeleeBase:
			s.meleeBase = append(s.meleeBase, e)
		case EntityTypeMeleeUnit:
			s.melee = append(s.melee, e)
		case EntityTypeRangedBase:
			s.rangedBase = append(s.rangedBase, e)
		case EntityTypeRangedUnit:
			s.ranged = append(s.ranged, e)
		case EntityTypeTurret:
			s.turel = append(s.turel, e)
		case EntityTypeBuilderUnit:
			s.miner = append(s.miner, e)
		case EntityTypeResource:
		}
	}

	for i, en := range s.enimies {
		sort.Sort(ByDistance(en))
		s.enimies[i] = en
		e := en
		if l := len(e); l > 10 {
			e = en[:int(l/3)]
		}
		mCenter := ByID(e).MedianCenter()
		if len(s.squards) > 0 {
			if distantion(s.center, mCenter) > 45.0 {
				continue
			}
		}
		s.defendPoint = append(s.defendPoint, Entity{Position: PointOnVector(s.center, mCenter, s.radius)})
		s.squards = append(s.squards, []Entity{})
	}

	sort.Sort(ByID(s.miner))
	sort.Sort(ByID(s.melee))
	sort.Sort(ByID(s.ranged))
	if len(s.miner) > 3 {
		s.builder = s.miner[:3]
		s.miner = s.miner[3:]
	}

	split := func(units []Entity) {
		for i := 0; i < len(units); i = i + len(s.squards) {
			for j := 0; j < len(s.squards); j++ {
				if len(units) <= i+j {
					return
				}
				e := units[i+j]
				s.squards[j] = append(s.squards[j], e)
			}
		}
	}

	arms := []Entity{}
	arms = append(s.melee, s.ranged...)
	split(arms)

	// if len(s.miner) > len(s.squards) {
	//
	// for j := 0; j < len(s.squards); j++ {
	// s.squards[j] = append(s.squards[j], s.miner[:1]...)
	// s.miner = s.miner[1:]
	// }
	// }

	s.population = s.Population()

	s.DoDefend(a)
	s.DoAction(a)
	s.DoDefendTurret(a)
	doRepair := false
	repair := func(e []Entity) {

		for _, h := range e {
			t := h.EntityType
			if h.Active && s.houseVec != nil && h.Position.Y == s.houseVec.Y && h.Position.X == s.houseVec.X {
				s.houseVec = nil
			}
			if !h.Active || h.Health < s.settings[t].MaxHealth {
				s.Repair(s.builder, a, h)
				doRepair = true
			}
		}
	}
	repair(s.rangedBase)
	repair(s.house)

	if s.currentTick > 500 && s.res < 100 {
		doRepair = true
	}

	if !doRepair {
		if s.res >= int(s.settings[EntityTypeRangedBase].InitialCost) && len(s.house) > 5 && len(s.rangedBase) < 2 && len(s.builder) >= 3 {
			s.BuildRangedBase(playerView.CurrentTick, s.builder[:2], a)
			if s.res >= int(s.settings[EntityTypeHouse].InitialCost) {
				s.BuildHouse(playerView.CurrentTick, s.builder[2:], a)
			}
		} else {
			if s.res >= int(s.settings[EntityTypeHouse].InitialCost) {
				s.BuildHouse(playerView.CurrentTick, s.builder, a)
			}
		}
	}
	// s.BuildTurret(a)
	s.Mining(s.miner, a)

	// if playerView.CurrentTick%2 == 0 {

	s.StopMakeMelee(a)
	s.StopMakeRanged(a)
	s.StopMakeBuilder(a)
	// }
	select {
	case e := <-s.buildChan:
		switch e {
		case EntityTypeBuilderUnit:
			if s.res > 300 && len(s.miner) > int(float64(s.Population()*0.6) {
				s.StopMakeBuilder(a)

				if playerView.CurrentTick%2 == 0 {
					s.MakeRanged(playerView.CurrentTick, a)
				} else {

					s.MakeMelee(playerView.CurrentTick, a)
				}
			} else {
				s.MakeBuilder(playerView.CurrentTick, a)
			}
		case EntityTypeRangedUnit:
			s.MakeRanged(playerView.CurrentTick, a)
		case EntityTypeMeleeUnit:
			s.MakeMelee(playerView.CurrentTick, a)
		}
	default:
	}
	if s.stopCallback == nil && int32(s.res) > s.settings[EntityTypeMeleeUnit].InitialCost {
		// if playerView.CurrentTick < 40 {

		// } else {
		// if len(s.miner) > 3 {
		// s.MakeRanged(playerView.CurrentTick, a)
		// s.MakeMelee(playerView.CurrentTick, a)
		// if ByID(squards).SumDamage(s) < 40 {
		// s.MakeBuilder(playerView.CurrentTick, a)
		// }
		// } else {
		// s.MakeBuilder(playerView.CurrentTick, a)
		// s.MakeRanged(playerView.CurrentTick, a)
		// s.MakeMelee(playerView.CurrentTick, a)
		//
		// }
		// }
	}

	if s.stopCallback != nil && s.currentTick == playerView.CurrentTick {
		s.stopCallback = nil
	}

	e := []Entity{}
	e = append(e, s.miner...)
	s.PrintId(e, "", black)
	s.PrintId(s.builder, "B", red)
	for i, sq := range s.squards {
		s.PrintId(sq, fmt.Sprintf("%dS", i), blue)
	}
	if s.houseVec != nil {
		s.PrintSpace(*s.houseVec, "", red, EntityTypeHouse)
	}

	s.PrintInfo()
	// for id, aa := range a.EntityActions {
	// if m := aa.MoveAction; m != nil {
	// fmt.Println("move", *m, id)
	// }
	// if m := aa.BuildAction; m != nil {
	// fmt.Println("build", *m, id)
	// }
	// if m := aa.AttackAction; m != nil {
	// fmt.Println("attack", *m, id)
	// }
	// if m := aa.RepairAction; m != nil {
	// fmt.Println("repair", *m, id)
	// }
	// }
	return a
}

func (s *MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
	// debugInterface.Send(DebugCommandClear{})
	debugInterface.GetState()
}

type ByDistance []Entity

func (a ByDistance) Len() int      { return len(a) }
func (a ByDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool {
	return distantion(Vec2Int32{}, a[i].Position) < distantion(Vec2Int32{}, a[j].Position)
}

type ByID []Entity

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].Id < a[j].Id }

func (a ByID) Center() Vec2Int32 {
	x := int32(0)
	y := int32(0)
	for _, e := range a {
		x += e.Position.X
		y += e.Position.Y
	}
	s := int32(len(a))
	if s == int32(0) {
		s = 1
	}
	return Vec2Int32{X: int32(x / s), Y: int32(y / s)}
}

func (a ByID) MedianCenter() Vec2Int32 {
	var x []int
	var y []int
	var v Vec2Int32
	for _, e := range a {
		x = append(x, int(e.Position.X))
		y = append(y, int(e.Position.Y))
	}
	sort.Ints(x)
	sort.Ints(y)
	l := len(a)
	if l == 1 {
		v.X = int32(x[0])
		v.Y = int32(y[0])
	}
	if l == 0 {
		v.X = 0
		v.Y = 0
	} else if l%2 == 0 {
		v.X = int32((x[l/2-1] + x[l/2]) / 2)
		v.Y = int32((y[l/2-1] + y[l/2]) / 2)
	} else {
		v.X = int32(x[l/2])
		v.Y = int32(y[l/2])
		//return Vec2Int32{X: medX, Y: medY}
	}
	return v
}

func (a ByID) SumDamage(s *MyStrategy) float32 {
	result := float32(0.0)
	for _, e := range a {
		result += float32(s.settings[e.EntityType].Attack.Damage)
	}
	return result
}

func (a ByID) AverageDistance(pos Vec2Int32) float32 {
	sum := float64(0.0)
	for _, e := range a {
		sum += distantion(pos, e.Position)
	}
	return float32(sum / float64(len(a)))
}

func (s *MyStrategy) PrintText(text string, pos Vec2Int32, color Color) {
	if s.debugInterface == nil {
		return
	}
	cv := NewColoredVertex(&Vec2Float32{float32(pos.X), float32(pos.Y)}, Vec2Float32{2.0, 2.0}, color)
	dd := NewDebugDataPlacedText(cv, text, 0.0, 12)
	s.debugInterface.Send(NewDebugCommandAdd(dd))
}

func (s *MyStrategy) PrintInfo() {
	if s.debugInterface == nil {
		return
	}
	pos := float32(0.0)
	size := float32(12.0)
	shift := float32(2.0)
	next := func() float32 {
		pos = pos + shift
		return pos
	}

	p := func(text string) {
		cv := NewColoredVertex(&Vec2Float32{next(), 0.0}, Vec2Float32{2.0, 2.0}, white)
		dd := NewDebugDataPlacedText(cv, text, 0.0, size)
		s.debugInterface.Send(NewDebugCommandAdd(dd))
	}

	text := fmt.Sprintf(`Builders: %d`, len(s.builder))
	p(text)

	text = fmt.Sprintf(`Miners: %d AVG: %f`, len(s.miner), ByID(s.miner).AverageDistance(s.center))
	p(text)

	for i, ss := range s.squards {
		squardLen := len(ss)
		ld := ByID(ss).SumDamage(s)
		prefix := ""
		if s.defendSquard == i {
			prefix = "Defender"
		}
		center := ByID(ss).Center()
		text = fmt.Sprintf(`%s Squard %d: %d, X: %d Y: %d SumDamage: %f`, prefix, i, squardLen, center.X, center.Y, ld)
		p(text)

	}
	for i, dp := range s.defendPoint {
		text = fmt.Sprintf(`Point %d, X: %d Y: %d`, i, dp.Position.X, dp.Position.Y)
		p(text)
	}

	for i, en := range s.enimies {
		center := Vec2Int32{}
		if len(en) > 10 {
			center = ByID(en[:10]).MedianCenter()
		} else {
			center = ByID(en).MedianCenter()
		}
		s.PrintText(fmt.Sprintf("%d", i), center, red)
		s.PrintVector(s.center, center, red)
		text = fmt.Sprintf(`Enimy: %d X: %d Y: %d`, i, center.X, center.Y)
		p(text)
	}

	text = fmt.Sprintf(`Melee: %d`, len(s.melee))
	p(text)

	text = fmt.Sprintf(`Ranged: %d`, len(s.ranged))
	p(text)

	text = fmt.Sprintf(`House: %d`, len(s.house))
	p(text)

	if en := s.enemyBasePosition; en != nil {
		text = fmt.Sprintf(`Enimy position: {x: %d, y: %d} ID: %d`, en.Position.X, en.Position.Y, en.PlayerId)
		p(text)
	}
}

func (s *MyStrategy) PrintId(e []Entity, prefix string, color Color) {
	if s.debugInterface == nil {
		return
	}
	for _, m := range e {
		cv := NewColoredVertex(&Vec2Float32{float32(m.Position.X), float32(m.Position.Y)}, Vec2Float32{2.0, 2.0}, color)

		dd := NewDebugDataPlacedText(cv, fmt.Sprintf("%s", prefix), 0.0, 12.0)
		s.debugInterface.Send(NewDebugCommandAdd(dd))
	}
}

func (s *MyStrategy) PrintSpace(e Vec2Int32, prefix string, color Color, t EntityType) {
	if s.debugInterface == nil {
		return
	}
	color.A = 0.5
	size := float32(s.settings[t].Size)
	colors := []ColoredVertex{}
	x := float32(e.X)
	y := float32(e.Y)

	cv := NewColoredVertex(&Vec2Float32{float32(x), float32(y)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{float32(x + size), float32(y)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{float32(x), float32(y + size)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)

	dd := NewDebugDataPrimitives(colors, PrimitiveTypeTriangles)
	s.debugInterface.Send(NewDebugCommandAdd(dd))

	cv = NewColoredVertex(&Vec2Float32{float32(x + size), float32(y + size)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{float32(x + size), float32(y)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{float32(x), float32(y + size)}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	dd = NewDebugDataPrimitives(colors, PrimitiveTypeTriangles)
	s.debugInterface.Send(NewDebugCommandAdd(dd))
}

func (s *MyStrategy) PrintVector(from Vec2Int32, to Vec2Int32, color Color) {
	if s.debugInterface == nil {
		return
	}
	colors := []ColoredVertex{}
	cv := NewColoredVertex(&Vec2Float32{float32(from.X), float32(from.Y)}, Vec2Float32{}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{float32(to.X), float32(to.Y)}, Vec2Float32{}, color)
	colors = append(colors, cv)

	dd := NewDebugDataPrimitives(colors, PrimitiveTypeLines)
	s.debugInterface.Send(NewDebugCommandAdd(dd))
}

func (s *MyStrategy) PrintRectange(from Vec2Int32, to Vec2Int32, color Color) {
	if s.debugInterface == nil {
		return
	}
	color.A = 0.5
	colors := []ColoredVertex{}
	fromX := float32(from.X)
	fromY := float32(from.Y)
	toX := float32(to.X)
	toY := float32(to.Y)

	cv := NewColoredVertex(&Vec2Float32{fromX, fromY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{toX, fromY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{fromX, toY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)

	dd := NewDebugDataPrimitives(colors, PrimitiveTypeTriangles)
	s.debugInterface.Send(NewDebugCommandAdd(dd))

	cv = NewColoredVertex(&Vec2Float32{toX, toY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{toX, fromY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	cv = NewColoredVertex(&Vec2Float32{fromY, toY}, Vec2Float32{2.0, 2.0}, color)
	colors = append(colors, cv)
	dd = NewDebugDataPrimitives(colors, PrimitiveTypeTriangles)
	s.debugInterface.Send(NewDebugCommandAdd(dd))
}
