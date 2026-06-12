package building

import "github/teohen/mgm-tto/cnts"

type Building interface {
	Draw()
	Pos() cnts.Point
	ID() string
	Type() BuildingType
}

type BuildingsList struct {
	list []Building
}

var buildings BuildingsList

func NewBuildingsList() *BuildingsList {
	list := make([]Building, 0)
	buildings = BuildingsList{
		list: list,
	}
	return &buildings
}

func (bl *BuildingsList) AddBuilding(b Building) {

	bl.list = append(bl.list, b)
}

func (bl *BuildingsList) GetBuildingAt(at cnts.Point) Building {
	for _, b := range bl.list {
		if b.Pos() == at {
			return b
		}
	}
	return nil
}

func (bl *BuildingsList) Buldings() []Building {
	return buildings.list
}

func (bl *BuildingsList) GetBuild(id string) *Building {
	for _, building := range bl.list {
		if building.ID() == id {
			return &building
		}
	}
	return nil
}

func (bl *BuildingsList) GetBuildingsOf(typ BuildingType) []Building {
	buildings := make([]Building, 0)
	for _, b := range bl.list {
		if b.Type() == typ {
			buildings = append(buildings, b)
		}
	}
	return buildings
}

func Get() *BuildingsList {
	return &buildings
}
