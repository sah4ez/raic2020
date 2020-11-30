package main

//
// import (
// . "aicup2020/model"
// )
//
// type MyStrategy struct {
// e      [][]*Entity
// miner  []*Entity
// melee  []*Entity
// ranged []*Entity
// res    int
// }
//
// func NewMyStrategy() *MyStrategy {
// return &MyStrategy{}
// }
//
// func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
// s.e = make([][]*Entity, int(playerView.MapSize))
// for i := range s.e {
// s.e[i] = make([]*Entity, int(playerView.MapSize))
// }
// s.melee = []*Entity{}
// s.ranged = []*Entity{}
// s.miner = []*Entity{}
//
// a := Action{
// EntityActions: make(map[int32]EntityAction),
// }
//
// for _, p := range playerView.Players {
// if p.Id == playerView.MyId {
// s.res = int(p.Resource)
// }
// }
//
// for _, e := range playerView.Entities {
// s.e[e.Position.X][e.Position.Y] = &e
// if e.PlayerId != nil && *e.PlayerId != playerView.MyId {
// continue
// }
// switch e.EntityType {
// case EntityTypeWall:
// case EntityTypeHouse:
// case EntityTypeBuilderBase:
// case EntityTypeMeleeBase:
// s.melee = append(s.melee, &e)
// case EntityTypeMeleeUnit:
// case EntityTypeRangedBase:
// case EntityTypeRangedUnit:
// s.ranged = append(s.ranged, &e)
// case EntityTypeTurret:
// case EntityTypeBuilderUnit:
// s.miner = append(s.miner, &e)
// case EntityTypeResource:
// }
// }
// for _, e := range s.miner {
// a.EntityActions[e.Id] = EntityAction{
// AttackAction: NewAttackAction(nil, NewAutoAttack(100, []EntityType{EntityTypeResource})),
// }
// }
//
// return a
// }
//
// func (strategy MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
// debugInterface.Send(DebugCommandClear{})
// debugInterface.GetState()
// }
//
// import (
// . "aicup2020/model"
// "fmt"
// )
//
// type MyStrategy struct {
// e            [][]*Entity
// capacity     int
// mapSize      int32
// resouces     map[int32]*Resource
// builders     map[int32]*Builder
// miners       map[int32]*Builder
// add          int
// builderBases map[int32]*BuilderBase
// hash         map[int32]struct{}
//
// miner           []*Builder
// builderHouse    []*Builder
// buildHouseQueue []*House
//
// builderPos int
// res        int
//
// buildingq map[int32]EntityAction
// }
//
// func (s *MyStrategy) GetBuilderPos() int32 {
// s.builderPos += 1
// if s.builderPos > 6 {
// s.builderPos = 0
// }
// return int32(s.builderPos)
// }
//
// func NewMyStrategy() *MyStrategy {
// return &MyStrategy{
// capacity:     15,
// resouces:     make(map[int32]*Resource, 20),
// builders:     make(map[int32]*Builder, 20),
// miners:       make(map[int32]*Builder, 20),
// builderBases: make(map[int32]*BuilderBase, 20),
// hash:         make(map[int32]struct{}, 20),
// buildingq:    make(map[int32]EntityAction, 20),
// }
// }
//
// type Resource struct {
// *Entity
// }
//
// func NewResouce(e *Entity) *Resource {
// return &Resource{
// Entity: e,
// }
// }
//
// type BuilderBase struct {
// *Entity
// }
//
// func NewBuilderBase(e *Entity) *BuilderBase {
// return &BuilderBase{
// Entity: e,
// }
// }
//
// type House struct {
// *Entity
// }
//
// func NewHouse(e *Entity) *House {
// return &House{
// Entity: e,
// }
// }
//
// type Builder struct {
// *Entity
// mining   bool
// building bool
// }
//
// func NewBuilder(e *Entity) *Builder {
// return &Builder{
// Entity: e,
// }
// }
//
// func (b *Builder) GetPosition(mapSize int32) (int, int) {
// return int(b.Position.X), int(b.Position.Y)
// }
//
// func (s *MyStrategy) find(posx, posy int) (e *Entity) {
// for i := 0; i < int(s.mapSize); i++ {
// e = s.e[posx+i][posy+i]
// if e != nil {
// if e.EntityType == EntityTypeResource {
// return e
// }
// }
// e = s.e[posx-i][posy+i]
// if e != nil {
// if e.EntityType == EntityTypeResource {
// return e
// }
// }
// e = s.e[posx-i][posy-i]
// if e != nil {
// if e.EntityType == EntityTypeResource {
// return e
// }
// }
// e = s.e[posx+i][posy-i]
// if e != nil {
// if e.EntityType == EntityTypeResource {
// return e
// }
// }
// }
// return
// }
//
// func (s *MyStrategy) Mining(a Action) {
// if s.res > 1000 {
// return
// }
// for id := range s.miner {
// fmt.Println(">> make miner", id)
// if _, ok := a.EntityActions[id]; !ok {
// a.EntityActions[id] = EntityAction{
// AttackAction: NewAttackAction(
// nil,
// NewAutoAttack(100, []EntityType{EntityTypeResource}),
// ),
// }
// }
// }
// }
//
// func (s *MyStrategy) MakeBuilder(a Action) {
// for id, b := range s.builderBases {
// fmt.Println(">> make builder", id)
// if _, ok := a.EntityActions[id]; !ok {
// a.EntityActions[id] = EntityAction{
// BuildAction: NewBuildAction(
// EntityTypeBuilderUnit,
// Vec2Int32{b.Position.X + s.GetBuilderPos(), b.Position.Y - 1},
// ),
// }
// }
// }
// }
//
// func (s *MyStrategy) BuildHouse(a Action) {
// if s.res < 200 {
// return
// }
// for id, b := range s.miners {
// fmt.Println(">> make house", id)
// a.EntityActions[id] = EntityAction{
// MoveAction: NewMoveAction(
// Vec2Int32{X: 0, Y: 0}, true, true,
// ),
// BuildAction: NewBuildAction(
// EntityTypeHouse,
// Vec2Int32{b.Position.X + 1, b.Position.Y + 1},
// ),
// }
// }
// }
//
// func (s *MyStrategy) AddEntity(e *Entity) {
// x := e.Position.X
// y := e.Position.Y
// s.e[x][y] = e
// }
//
// func (s *MyStrategy) AddBuilderBase(e *Entity) {
// s.AddEntity(e)
// var (
// ok bool
// ee *BuilderBase
// )
// if ee, ok = s.builderBases[e.Id]; !ok {
// ee = NewBuilderBase(e)
// s.builderBases[e.Id] = ee
// return
// }
// ee.Position.X = e.Position.X
// ee.Position.Y = e.Position.Y
// s.builderBases[e.Id] = ee
// }
//
// func (s *MyStrategy) AddBuilder(e *Entity) {
// s.AddEntity(e)
// var (
// ok           bool
// foundBuilder *Builder
// )
// if foundBuilder, ok = s.builders[e.Id]; !ok {
// foundBuilder = NewBuilder(e)
// s.builders[e.Id] = foundBuilder
// return
// }
// foundBuilder.Position.X = e.Position.X
// foundBuilder.Position.Y = e.Position.Y
// s.builders[e.Id] = foundBuilder
// }
//
// func (s *MyStrategy) AddMiner(e *Entity) {
// if len(s.miner) < 10 {
// s.miner = append(s.miner, NewBuilder(e))
// } else if len(s.builderHouse) < 2 {
// s.builderHouse = append(s.builderHouse, NewBuilder(e))
// }
// s.AddEntity(e)
// if len(s.miners) == 12 {
// s.AddBuilder(e)
// return
// }
// var (
// ok           bool
// foundBuilder *Builder
// )
// if foundBuilder, ok = s.miners[e.Id]; !ok {
// foundBuilder = NewBuilder(e)
// s.miners[e.Id] = foundBuilder
// return
// }
// foundBuilder.Position.X = e.Position.X
// foundBuilder.Position.Y = e.Position.Y
// s.miners[e.Id] = foundBuilder
// }
//
// func (s *MyStrategy) AddResource(e *Entity) {
// var (
// ok           bool
// foundResouce *Resource
// )
//
// if e.Active {
// s.AddEntity(e)
// }
//
// if foundResouce, ok = s.resouces[e.Id]; !ok {
// if e.Active {
// s.resouces[e.Id] = NewResouce(e)
// }
// }
//
// if !e.Active && foundResouce != nil {
// foundResouce.Active = false
// s.resouces[e.Id] = foundResouce
// }
// }
//
// func (s *MyStrategy) getAction(playerView PlayerView, debugInterface *DebugInterface) Action {
// s.mapSize = playerView.MapSize
// s.e = make([][]*Entity, int(playerView.MapSize))
// for i := range s.e {
// s.e[i] = make([]*Entity, int(playerView.MapSize))
// }
// s.miner = make([]*Builder, len(s.miner))
//
// a := Action{
// EntityActions: make(map[int32]EntityAction),
// }
// for _, p := range playerView.Players {
// if p.Id == playerView.MyId {
// s.res = int(p.Resource)
// }
// }
//
// for _, e := range playerView.Entities {
// switch e.EntityType {
//
// case EntityTypeWall:
// case EntityTypeHouse:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// case EntityTypeBuilderBase:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// s.AddBuilderBase(&e)
// a.EntityActions[e.Id] = EntityAction{
// BuildAction: NewBuildAction(
// EntityTypeBuilderUnit,
// Vec2Int32{e.Position.X + s.GetBuilderPos(), e.Position.Y - 1},
// ),
// }
// case EntityTypeMeleeBase:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// case EntityTypeMeleeUnit:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// if _, ok := s.hash[e.Id]; !ok {
// s.hash[e.Id] = struct{}{}
// s.add = s.add + 1
// }
// case EntityTypeRangedBase:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// case EntityTypeRangedUnit:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// if _, ok := s.hash[e.Id]; !ok {
// s.hash[e.Id] = struct{}{}
// s.add = s.add + 1
// }
// case EntityTypeTurret:
// case EntityTypeBuilderUnit:
// if *e.PlayerId != playerView.MyId {
// continue
// }
// s.AddMiner(&e)
// s.AddBuilder(&e)
// case EntityTypeResource:
// s.AddResource(&e)
// }
// }
//
// s.MakeBuilder(a)
// s.Mining(a)
// s.BuildHouse(a)
//
// return a
// }
//
// func (strategy MyStrategy) debugUpdate(playerView PlayerView, debugInterface DebugInterface) {
// debugInterface.Send(DebugCommandClear{})
// debugInterface.GetState()
// }
