package main

import (
	. "aicup2020/model"
	"math/rand"
	"sync"
)

type MyStrategy struct {
	settings    map[EntityType]EntityProperties
	e           [][]Entity
	miner       []int32
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
	builder      []int32
	lastBuilder  int32
	lastMelee    int32
	lastRanged   int32
	enimies      map[int32][]Entity
	countBuilder int

	size int
	res  int
}

func NewMyStrategy() *MyStrategy {
	return &MyStrategy{
		countBuilder: 2,
		position:     0,
		buildAction: []string{
			// "builder", "builder", "builder",
			"builder", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "ranged", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
			// "melee", "melee", "ranged",
		},
		houseVec: []Vec2Int32{
			Vec2Int32{1, 2},
			Vec2Int32{2, 2},
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
		if len(units) < 3 {
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

func (s *MyStrategy) Mining(units []int32, a Action) {
	for _, e := range units {

		aaa := NewAutoAttack(32, []EntityType{EntityTypeResource})
		aa := NewAttackAction(nil, &aaa)
		a.EntityActions[e] = EntityAction{
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

func (s *MyStrategy) FindSpaceForTurret() []Vec2Int32 {
	l := int32(len(s.turel))
	v0 := s.houseVec[0]
	v1 := s.houseVec[1]
	shift := int32(17)

	size := s.settings[EntityTypeTurret].Size
	if l%2 != 0 {
		return []Vec2Int32{
			Vec2Int32{v0.X + shift, v0.Y + shift + l*size},
			Vec2Int32{v1.X + shift, v1.Y + shift + l*size},
		}
	}
	return []Vec2Int32{
		Vec2Int32{v0.X + shift + l*size, v0.Y + shift},
		Vec2Int32{v1.X + shift + l*size, v1.Y + shift},
	}
}

func (s *MyStrategy) FindSpaceForHouse() []Vec2Int32 {
	l := int32(len(s.house))
	v0 := s.houseVec[0]
	v1 := s.houseVec[1]

	shift := int32(0)
	if l > 2 {
		shift = -2
	}
	size := s.settings[EntityTypeHouse].Size
	if l%2 != 0 {
		return []Vec2Int32{
			Vec2Int32{v0.X, v0.Y + l*size - shift},
			Vec2Int32{v1.X, v1.Y + l*size - shift},
		}
	}
	return []Vec2Int32{
		Vec2Int32{v0.X + l*size - shift, v0.Y},
		Vec2Int32{v1.X + l*size - shift, v1.Y},
	}

}

func (s *MyStrategy) BuildHouse(c int32, units []int32, a Action) {
	for _, h := range s.turel {
		if !h.Active || h.Health < s.settings[EntityTypeTurret].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}
	for _, h := range s.house {
		if !h.Active || h.Health < s.settings[EntityTypeHouse].MaxHealth {
			s.Repair(c, units, a, h)
			return
		}
	}

	for i, e := range units {

		var p []Vec2Int32
		if len(s.turel) <= 1 && len(s.house) == 4 {
			p = s.FindSpaceForTurret()
			v := p[1]
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
				p[1],
			)
			a.EntityActions[e] = EntityAction{
				MoveAction:  &ma,
				BuildAction: &ba,
			}
		} else {
			p = s.FindSpaceForHouse()
			v := p[1]
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
				EntityTypeHouse,
				p[1],
			)
			a.EntityActions[e] = EntityAction{
				MoveAction:  &ma,
				BuildAction: &ba,
			}
		}

	}
}

func (s *MyStrategy) BuildTurret(c int32, units []int32, a Action) {
	t := EntityTypeTurret
	var he *Entity
	for _, h := range s.turel {
		if !h.Active || h.Health < s.settings[t].MaxHealth {
			he = &h
			break
		}
	}
	for i, e := range units {

		if he != nil {
			s.Repair(c, units, a, *he)
			return
		}

		p := s.FindSpaceForHouse()

		v := p[0]
		v.Y = v.Y + int32(i)
		ma := NewMoveAction(
			v, true, true,
		)
		ba := NewBuildAction(
			EntityTypeHouse,
			p[1],
		)
		a.EntityActions[e] = EntityAction{
			MoveAction:  &ma,
			BuildAction: &ba,
		}
	}
}

func (s *MyStrategy) Repair(c int32, units []int32, a Action, tgt Entity) {
	for _, e := range units {
		s.stopBuild = c + 1
		ma := NewMoveAction(
			Vec2Int32{X: tgt.Position.X, Y: tgt.Position.Y},
			true,
			true,
		)
		ra := NewRepairAction(tgt.Id)
		a.EntityActions[e] = EntityAction{MoveAction: &ma, RepairAction: &ra}
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

func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
	if s.settings == nil {
		s.settings = playerView.EntityProperties
	}

	s.prevPopulation = s.population
	s.e = make([][]Entity, int(playerView.MapSize))
	for i := range s.e {
		s.e[i] = make([]Entity, int(playerView.MapSize))
	}
	s.miner = []int32{}
	s.builder = []int32{}
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
				s.miner = append(s.miner, e.Id)
			} else if len(s.miner)%2 == 0 && len(s.builder) < s.countBuilder {
				s.builder = append(s.builder, e.Id)
			} else {
				s.miner = append(s.miner, e.Id)
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

	if s.Population() >= s.Capacity()-5 {
		if s.res >= int(s.settings[EntityTypeTurret].Cost) && len(s.turel) < 2 {

		}
		if s.res >= int(s.settings[EntityTypeHouse].Cost) {
			// fmt.Println(">>", len(s.builder))
			s.BuildHouse(playerView.CurrentTick, s.builder, a)
		}
	} else {
		// s.Mining(s.builder, a)
	}
	s.Mining(s.miner, a)

	if s.stopCallback == nil && int32(s.res) >= playerView.EntityProperties[EntityTypeMeleeUnit].Cost {
		// fmt.Println(">>>", "start")

		if len(s.miner) < 3 {
			s.MakeBuilder(playerView.CurrentTick, a)
			s.stopCallback = s.StopMakeBuilder
		} else {
			e := s.NextBuild()
			if len(s.miner)+len(s.builder) >= 30 && e == "builder" {
				e = s.NextBuild()
			}
			// fmt.Println(e, s.res, s.Population())
			switch e {
			case "builder":
				s.MakeBuilder(playerView.CurrentTick, a)
				s.stopCallback = s.StopMakeBuilder
			case "melee":
				s.MakeMelee(playerView.CurrentTick, a)
				s.stopCallback = s.StopMakeMelee
			case "ranged":
				s.MakeRanged(playerView.CurrentTick, a)
				s.stopCallback = s.StopMakeRanged
			}
		}
	}

	if playerView.CurrentTick > 0 && playerView.CurrentTick == s.stopBuild {
		// fmt.Println(">>>", "stop")
		if s.stopCallback != nil {
			s.stopCallback(a)
			s.stopCallback = nil
		}
	}

	return a
}

func (strategy MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
	debugInterface.Send(DebugCommandClear{})
	debugInterface.GetState()
}
