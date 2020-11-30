package main

import (
	. "aicup2020/model"
	"fmt"
	"sync"
)

type MyStrategy struct {
	e           [][]Entity
	miner       []Entity
	melee       []Entity
	ranged      []Entity
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

	building bool

	size int
	res  int
}

func NewMyStrategy() *MyStrategy {
	return &MyStrategy{
		position:    0,
		buildAction: []string{"melee", "ranged", "builder"},
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
	return len(s.miner) + len(s.melee) + len(s.ranged)
}

func (s *MyStrategy) Attack(units []Entity, a Action) {
	for _, e := range units {
		a.EntityActions[e.Id] = EntityAction{
			MoveAction: NewMoveAction(
				Vec2Int32{X: e.Position.X + 1, Y: e.Position.Y + 1},
				true,
				true,
			),
			AttackAction: NewAttackAction(
				nil,
				NewAutoAttack(1000, []EntityType{EntityTypeMeleeUnit}),
			),
		}
	}
}

func (s *MyStrategy) Mining(units []Entity, a Action) {
	for _, e := range units {
		fmt.Println(">>mine", e.Id)
		a.EntityActions[e.Id] = EntityAction{
			AttackAction: NewAttackAction(
				nil,
				NewAutoAttack(1000, []EntityType{EntityTypeMeleeUnit}),
			),
		}
	}
}

func (s *MyStrategy) MakeRanged(a Action) {
	for _, e := range s.rangedBase {
		a.EntityActions[e.Id] = EntityAction{
			BuildAction: NewBuildAction(
				EntityTypeRangedUnit,
				Vec2Int32{X: e.Position.X, Y: e.Position.Y - 1},
			),
		}
	}
}

func (s *MyStrategy) MakeMelee(c int32, a Action) {
	for _, e := range s.meleeBase {
		s.stopBuild = c + 2
		a.EntityActions[e.Id] = EntityAction{
			BuildAction: NewBuildAction(
				EntityTypeMeleeUnit,
				Vec2Int32{X: e.Position.X, Y: e.Position.Y - 1},
			),
		}
	}
}

func (s *MyStrategy) StopMakeMelee(a Action) {
	for _, e := range s.meleeBase {
		a.EntityActions[e.Id] = EntityAction{}
	}
}

func (s *MyStrategy) MakeBuilder(a Action) {
	if len(s.miner) < 9 {
		return
	}
	for _, e := range s.builderBase {
		a.EntityActions[e.Id] = EntityAction{
			BuildAction: NewBuildAction(
				EntityTypeBuilderUnit,
				Vec2Int32{X: e.Position.X, Y: e.Position.Y - 1},
			),
		}
	}
}

func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
	s.prevPopulation = s.population
	s.e = make([][]Entity, int(playerView.MapSize))
	for i := range s.e {
		s.e[i] = make([]Entity, int(playerView.MapSize))
	}
	s.miner = []Entity{}
	s.melee = []Entity{}
	s.ranged = []Entity{}
	s.house = []Entity{}
	s.builderBase = []Entity{}
	s.meleeBase = []Entity{}
	s.rangedBase = []Entity{}

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
		case EntityTypeBuilderUnit:
			s.miner = append(s.miner, e)
		case EntityTypeResource:
		}
	}

	s.population = s.Population()

	s.Mining(s.miner, a)
	s.Attack(s.melee, a)
	s.Attack(s.ranged, a)

	if int32(s.res) == playerView.EntityProperties[EntityTypeMeleeUnit].Cost {
		s.MakeMelee(playerView.CurrentTick, a)

		// s.building = true
		// e := s.NextBuild()
		// fmt.Println(e, s.res, s.Population())
		// switch e {
		// case "melee":
		// case "ranged":
		// s.MakeRanged(a)
		// case "builder":
		// s.MakeBuilder(a)
		// }
	}

	if playerView.CurrentTick > 0 && playerView.CurrentTick == s.stopBuild {
		fmt.Println(">>>", "stop")
		s.StopMakeMelee(a)
	}

	return a
}

func (strategy MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
	debugInterface.Send(DebugCommandClear{})
	debugInterface.GetState()
}
