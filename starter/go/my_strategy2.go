package main

import (
	. "aicup2020/model"
	"fmt"
	"math/rand"
	"sync"
)

type MyStrategy struct {
	settings    map[EntityType]EntityProperties
	e           [][]Entity
	miner       []Entity
	minerH      map[int32]bool
	melee       []Entity
	ranged      []Entity
	turel       []Entity
	house       []Entity
	builderBase []Entity
	meleeBase   []Entity
	rangedBase  []Entity

	position       int
	stopBuild      int32
	prevPopulation int
	population     int
	buildAction    []string
	one            sync.Once
	two            sync.Once

	stopCallback func(Action)
	houseVec     []Vec2Int32
	turretVec    []Vec2Int32
	last         int
	builder      []Entity
	turelBuilder Entity
	lastBuilder  int32
	lastMelee    int32
	lastRanged   int32
	enimies      map[int32][]Entity
	countBuilder int

	size int
	tick int32
	res  int
}

func NewMyStrategy() *MyStrategy {
	return &MyStrategy{
		countBuilder: 2,
		position:     0,
		buildAction: []string{
			// "builder", "builder", "builder",
			"builder", "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "ranged", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
		},
		houseVec: []Vec2Int32{
			Vec2Int32{1, 1},
		},
		turretVec: []Vec2Int32{
			Vec2Int32{15, 15},
		},
		lastMelee:  int32(1),
		lastRanged: int32(1),
	}
}

func (s *MyStrategy) NextBuild() string {
	p := s.position
	if p >= len(s.buildAction) {
		p = 0
		s.position = 0
	}
	defer func() {
		s.position++
	}()
	return s.buildAction[p]
}

func (s *MyStrategy) Population() int {
	return len(s.miner)*int(s.settings[EntityTypeBuilderUnit].PopulationUse) +
		len(s.melee)*int(s.settings[EntityTypeMeleeUnit].PopulationUse) +
		len(s.ranged)*int(s.settings[EntityTypeRangedUnit].PopulationUse) +
		len(s.builder)*int(s.settings[EntityTypeBuilderUnit].PopulationUse)
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

func (s *MyStrategy) Attack(units []Entity, a Action) {
	if len(units) < 1 {
		return
	}
	if s.settings[units[0].EntityType].CanMove {
		if len(units) < 10 {
			for _, e := range units {

				var aaa AutoAttack
				if e.EntityType == EntityTypeRangedUnit {

					size := s.settings[EntityTypeRangedBase].Size
					i := int32(rand.Intn(10))
					lu := int32(len(units)) + size + i

					if len(s.rangedBase) > 0 {
						ma := NewMoveAction(
							Vec2Int32{X: s.rangedBase[0].Position.X + lu, Y: s.rangedBase[0].Position.Y + lu},
							true,
							true,
						)
						aaa = NewAutoAttack(10, []EntityType{EntityTypeMeleeUnit, EntityTypeRangedUnit, EntityTypeRangedBase, EntityTypeMeleeBase, EntityTypeTurret, EntityTypeRangedBase, EntityTypeBuilderBase, EntityTypeBuilderUnit})
						aa := NewAttackAction(nil, &aaa)
						a.EntityActions[e.Id] = EntityAction{MoveAction: &ma, AttackAction: &aa}
					}

				} else if e.EntityType == EntityTypeMeleeUnit {

					size := s.settings[EntityTypeMeleeBase].Size
					i := int32(rand.Intn(10))
					lu := int32(len(units)) + size + i
					if len(s.meleeBase) > 0 {
						ma := NewMoveAction(
							Vec2Int32{X: s.meleeBase[0].Position.X + lu, Y: s.meleeBase[0].Position.Y + lu},
							true,
							true,
						)
						aaa = NewAutoAttack(10, []EntityType{EntityTypeMeleeUnit, EntityTypeRangedUnit, EntityTypeRangedBase, EntityTypeMeleeBase, EntityTypeTurret, EntityTypeRangedBase, EntityTypeBuilderBase, EntityTypeBuilderUnit})
						aa := NewAttackAction(nil, &aaa)
						a.EntityActions[e.Id] = EntityAction{MoveAction: &ma, AttackAction: &aa}

					}

				}
			}
			return
		}
	}

	for _, e := range units {
		i := int32(rand.Intn(3))
		ma := NewMoveAction(
			Vec2Int32{X: e.Position.X + i, Y: e.Position.Y + i},
			true,
			true,
		)
		var aaa AutoAttack
		if e.EntityType == EntityTypeTurret {
			aaa = NewAutoAttack(1000, []EntityType{EntityTypeMeleeUnit, EntityTypeRangedUnit, EntityTypeMeleeBase})
		} else {
			aaa = NewAutoAttack(1000, []EntityType{EntityTypeMeleeUnit, EntityTypeRangedUnit, EntityTypeRangedBase, EntityTypeMeleeBase, EntityTypeTurret, EntityTypeRangedBase, EntityTypeBuilderBase, EntityTypeBuilderUnit})
		}
		aa := NewAttackAction(nil, &aaa)
		a.EntityActions[e.Id] = EntityAction{
			MoveAction:   &ma,
			AttackAction: &aa,
		}
	}
}

func (s *MyStrategy) Mining(units []Entity, a Action) {
	for _, e := range units {
		var ma MoveAction
		if (e.Position.X < 13 || e.Position.Y < 13) && s.tick > 400 {
			ma = NewMoveAction(Vec2Int32{16, 19}, true, true)
			if s.tick > 800 {

				ma = NewMoveAction(Vec2Int32{40, 40}, true, true)
			}
		}

		aaa := NewAutoAttack(int32(s.size), []EntityType{EntityTypeResource})
		aa := NewAttackAction(nil, &aaa)
		a.EntityActions[e.Id] = EntityAction{
			MoveAction:   &ma,
			AttackAction: &aa,
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

func (s *MyStrategy) FindSpaceForTurret() Vec2Int32 {
	l := int32(len(s.turel))
	v1 := s.turretVec[0]

	shift := int32(int(l/2) + 1)

	size := s.settings[EntityTypeTurret].Size
	if l%2 != 0 {
		return Vec2Int32{v1.X - shift*size - 1, v1.Y + shift*size + 2}
	}
	return Vec2Int32{v1.X + shift*size + 2, v1.Y - shift*size - 1}
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
	if len(s.builderBase) < 0 {

		return Vec2Int32{21, 4}
	}
	// if len(s.rangedBase) == 1 {
	//
	// return Vec2Int32{10, 5}
	// } else if len(s.rangedBase) == 0 {
	// return Vec2Int32{10, 5}
	// } else {
	// if len(s.meleeBase) > 1 {
	//
	// return Vec2Int32{11, 21}
	// } else {
	// return Vec2Int32{5, 10}
	// }
	// }
	size := s.settings[EntityTypeBuilderBase].Size
	p := s.builderBase[0].Position
	return Vec2Int32{p.X + size + 1, p.Y + size + 1}
}

func (s *MyStrategy) roundFunc(t EntityType) (v Vec2Int32) {
	v = Vec2Int32{1, 1}
	// fmt.Println("<<<<", v)
	for v.X <= int32(s.size-2) {
		// fmt.Println(">>>>", v)
		switch {
		case v.X == 1 || v.Y == 1:
			if s.IsFree(v, t) {
				return v
			} else {
				v.Y++
			}
		case v.X == 1 || v.Y != 1:
			for v.Y > 1 {
				if s.IsFree(v, t) {
					return v
				} else {
					v.Y--
					v.X++
				}
			}
		case v.X != 1 || v.Y == 1:
			v.X++
			if s.IsFree(v, t) {
				return v
			} else {
				for v.X > 1 {
					if s.IsFree(v, t) {
						return v
					} else {
						v.X--
						v.Y++
					}
				}
			}
			v.Y++
			if s.IsFree(v, t) {
				return v
			}
		}
	}
	fmt.Println("!>>>", v)
	return v
}

func (s *MyStrategy) FindSpaceForHouse() Vec2Int32 {
	l := int32(len(s.house))
	v1 := s.roundFunc(EntityTypeHouse)
	// v1 := s.houseVec[0]
	// поиск свободного места, не работает )))
	// if l >= 0 && len(s.builder) > 0 {
	// e := s.builder[0]
	// x := e.Position.X
	// y := e.Position.Y
	//
	// v := Vec2Int32{X: x, Y: y}
	//
	// next := nextFunc(int32(s.size), int32(s.size), x, y)
	// for s.IsFree(v, EntityTypeHouse) {
	// x, y := next()
	// v = Vec2Int32{X: x, Y: y}
	// }
	// return v
	// }

	size := s.settings[EntityTypeHouse].Size
	shift := int32(int(l/2) + 1)
	if l%2 != 0 {
		return Vec2Int32{v1.X, v1.Y + shift*size + 1}
	}
	return Vec2Int32{v1.X + shift*size + 1, v1.Y}
}

func (s *MyStrategy) BuildTurret(c int32, units []Entity, a Action) bool {
	for _, h := range s.turel {
		if !h.Active || h.Health < s.settings[EntityTypeTurret].MaxHealth {
			s.Repair(c, units, a, h)
			return true
		}
	}

	result := false
	for i, e := range units {

		coeff := 2
		if len(s.house) > 0 && len(s.house)%coeff == 0 && len(s.turel) < len(s.house)/coeff {
			p := s.FindSpaceForTurret()
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
				EntityTypeTurret,
				p,
			)
			a.EntityActions[e.Id] = EntityAction{
				MoveAction:  &ma,
				BuildAction: &ba,
			}
			result = true
		}
	}
	return result
}

func (s *MyStrategy) BuildMeleeBase(c int32, units []Entity, a Action) {
	if len(s.meleeBase) > 2 || len(s.house) < 7 {
		return
	}

	for _, h := range s.meleeBase {
		if !h.Active || h.Health < s.settings[EntityTypeMeleeBase].MaxHealth {
			s.Repair(c, units, a, h)
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
		a.EntityActions[e.Id] = EntityAction{
			MoveAction:  &ma,
			BuildAction: &ba,
		}

	}
}

func (s *MyStrategy) BuildRangedBase(c int32, units []Entity, a Action) {
	if len(s.rangedBase) > 2 || len(s.house) < 7 {
		return
	}

	for _, h := range s.rangedBase {
		if !h.Active || h.Health < s.settings[EntityTypeRangedBase].MaxHealth {
			s.Repair(c, units, a, h)
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
		a.EntityActions[e.Id] = EntityAction{
			MoveAction:  &ma,
			BuildAction: &ba,
		}

	}
}

func (s *MyStrategy) BuildHouse(c int32, units []Entity, a Action) {
	for _, h := range s.house {
		if !h.Active || h.Health < s.settings[EntityTypeHouse].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}
	for _, h := range s.builderBase {
		if !h.Active || h.Health < s.settings[EntityTypeBuilderBase].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}
	for _, h := range s.meleeBase {
		if !h.Active || h.Health < s.settings[EntityTypeMeleeBase].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}
	for _, h := range s.rangedBase {
		if !h.Active || h.Health < s.settings[EntityTypeRangedBase].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}

	// for i, e := range units {
	for _, e := range units {

		var p Vec2Int32
		p = s.FindSpaceForHouse()
		v := p

		v.X = v.X - int32(1)
		v.Y = v.Y - int32(1)
		// if len(s.house)%2 != 0 {
		// v.X = v.X - int32(1)
		// v.Y = v.Y + int32(i)
		//
		// } else {
		// v.X = v.X + int32(i)
		// v.Y = v.Y - int32(1)
		// }
		ma := NewMoveAction(
			v, true, true,
		)
		ba := NewBuildAction(
			EntityTypeHouse,
			p,
		)
		a.EntityActions[e.Id] = EntityAction{
			MoveAction:  &ma,
			BuildAction: &ba,
		}

	}
}

func (s *MyStrategy) Repair(c int32, units []Entity, a Action, tgt Entity) {
	for _, e := range units {
		s.stopBuild = c + 1
		ma := NewMoveAction(
			Vec2Int32{X: tgt.Position.X, Y: tgt.Position.Y},
			true,
			true,
		)
		ra := NewRepairAction(tgt.Id)
		a.EntityActions[e.Id] = EntityAction{MoveAction: &ma, RepairAction: &ra}
	}
}

func (s *MyStrategy) MakeRanged(c int32, a Action) {
	size := s.settings[EntityTypeRangedBase].Size

	for _, e := range s.rangedBase {
		s.stopBuild = c + 2
		pos := Vec2Int32{X: e.Position.X + size, Y: e.Position.Y + size - s.lastRanged}
		// fmt.Println(">>>r", pos)
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
	size := s.settings[EntityTypeMeleeBase].Size

	for _, e := range s.meleeBase {
		s.stopBuild = c + 2
		pos := Vec2Int32{X: e.Position.X + size, Y: e.Position.Y + size - s.lastMelee}
		// fmt.Println(">>>", pos)
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
	size := s.settings[EntityTypeBuilderBase].Size
	unit := s.settings[EntityTypeBuilderUnit].Size
	for _, e := range s.builderBase {
		s.stopBuild = c + 2
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

func (s *MyStrategy) IsFree(pos Vec2Int32, t EntityType) bool {
	size := s.settings[t].Size
	free := true
	if (pos.X < 0 || pos.X > int32(len(s.e)-1)) || (pos.Y < 0 || pos.Y > int32(len(s.e)-1)) {
		return false
	}
	for i := pos.X; i < pos.X+size; i++ {
		for j := pos.Y; j < pos.Y+size; j++ {
			free = (free && s.e[i][j].Id == int32(0))
			if !free {
				break
			}
		}
	}

	s.e[pos.X][pos.Y] = Entity{Id: int32(-1)}
	return free
}

func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
	if s.settings == nil {
		s.settings = playerView.EntityProperties
	}
	s.size = int(playerView.MapSize)

	s.tick = playerView.CurrentTick
	s.prevPopulation = s.population
	s.e = make([][]Entity, int(playerView.MapSize))
	for i := range s.e {
		s.e[i] = make([]Entity, int(playerView.MapSize))
	}
	s.turelBuilder = Entity{}
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

	a := Action{
		EntityActions: make(map[int32]EntityAction),
	}

	for _, p := range playerView.Players {
		if p.Id == playerView.MyId {
			s.res = int(p.Resource)
		}
	}

	for _, e := range playerView.Entities {
		s.e[e.Position.X][e.Position.Y] = e
		for i := e.Position.X; i <= e.Position.X+s.settings[e.EntityType].Size-1; i++ {
			for j := e.Position.Y; j <= e.Position.Y+s.settings[e.EntityType].Size-1; j++ {
				s.e[i][j] = e
			}
		}
		if e.PlayerId == nil {
			continue
		}
		if *e.PlayerId != playerView.MyId {
			if e.EntityType == EntityTypeBuilderBase {
				c, ok := s.enimies[*e.PlayerId]
				if ok {
					c = append(c, e)
					s.enimies[*e.PlayerId] = c
					continue
				}
				s.enimies[*e.PlayerId] = []Entity{e}
			}
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

			if len(s.miner) == 0 {
				s.miner = append(s.miner, e)
			} else if len(s.miner)%2 == 0 && len(s.builder) < s.countBuilder {
				s.builder = append(s.builder, e)
			} else if s.turelBuilder.Id == int32(0) {
				s.turelBuilder = e
			} else {
				s.miner = append(s.miner, e)
			}
		case EntityTypeResource:
		}
	}

	s.population = s.Population()

	s.Attack(s.melee, a)
	s.Attack(s.ranged, a)
	s.Attack(s.turel, a)

	// shift := 0
	for _, h := range s.house {
		if !h.Active || h.Health < s.settings[EntityTypeHouse].MaxHealth {
			s.Repair(playerView.CurrentTick, s.builder, a, h)
		}
	}
	if s.res >= int(s.settings[EntityTypeTurret].Cost) || len(s.rangedBase) > 1 {
		units := []Entity{s.turelBuilder}
		ok := s.BuildTurret(playerView.CurrentTick, units, a)
		if !ok {
			s.Mining(units, a)
		}
	}
	if s.res >= int(s.settings[EntityTypeRangedBase].Cost) && len(s.meleeBase) < 2 && s.IsFree(Vec2Int32{18, 18}, EntityTypeRangedBase) {
		s.BuildRangedBase(playerView.CurrentTick, s.builder, a)
	} else if s.Population() >= s.Capacity()-10 {
		if s.res >= int(s.settings[EntityTypeHouse].Cost) {
			s.BuildHouse(playerView.CurrentTick, s.builder, a)
		}
	} else {
		// s.Mining(s.builder, a)
	}
	s.Mining(s.miner, a)

	//if s.stopCallback == nil && int32(s.res) >= playerView.EntityProperties[EntityTypeMeleeUnit].Cost {
	//	// fmt.Println(">>>", "start")

	//	if len(s.miner) < 3 {
	//		s.MakeBuilder(playerView.CurrentTick, a)
	//		s.stopCallback = s.StopMakeBuilder
	//	} else {
	//		e := s.NextBuild()
	//		if len(s.miner)+len(s.builder) >= int(50) && e == "builder" {
	//			e = s.NextBuild()
	//		}
	//		// fmt.Println(e, s.res, s.Population())
	//		switch e {
	//		case "builder":
	//			s.MakeBuilder(playerView.CurrentTick, a)
	//			s.stopCallback = s.StopMakeBuilder
	//		case "melee":
	//			s.MakeMelee(playerView.CurrentTick, a)
	//			s.stopCallback = s.StopMakeMelee
	//		case "ranged":
	//			s.MakeRanged(playerView.CurrentTick, a)
	//			s.stopCallback = s.StopMakeRanged
	//		}
	//	}
	//}

	//if playerView.CurrentTick > 0 && playerView.CurrentTick == s.stopBuild {
	//	// fmt.Println(">>>", "stop")
	//	if s.stopCallback != nil {
	//		s.stopCallback(a)
	//		s.stopCallback = nil
	//	}
	//}
	// if  int32(s.res) >= playerView.EntityProperties[EntityTypeMeleeUnit].Cost {
	// fmt.Println(">>>", "start")

	if len(s.miner) < 3 {
		s.MakeBuilder(playerView.CurrentTick, a)
	} else {
		s.MakeBuilder(playerView.CurrentTick, a)
		s.MakeMelee(playerView.CurrentTick, a)
		s.MakeRanged(playerView.CurrentTick, a)
	}

	return a
}

func (strategy MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
	debugInterface.Send(DebugCommandClear{})
	debugInterface.GetState()
}

func nextFunc(m, n int32, x, y int32) func() (int32, int32) {
	top, down := int32(0), int32(m-1)
	left, right := int32(0), int32(n-1)
	// x, y := int32(0), int32(-1)
	dx, dy := int32(0), int32(1)
	return func() (int32, int32) {
		x += dx
		y += dy
		switch {
		case y+dy > right:
			top++
			dx, dy = 1, 0
		case x+dx > down:
			right--
			dx, dy = 0, -1
		case y+dy < left:
			down--
			dx, dy = -1, 0
		case x+dx < top:
			left++
			dx, dy = 0, 1
		}
		return x, y
	}
}
