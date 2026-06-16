package npc

import "github/teohen/mgm-tto/cnts"

type NPCType string

const (
	VillagerNPCType NPCType = "villager"
	DeerNPCType     NPCType = "deer"
)

type NPC interface {
	Draw()
	Pos() cnts.Point
	ID() string
	Type() NPCType
	Tick()
}

type NPCList struct {
	list []NPC
}

var npcs NPCList

func NewNPCList() *NPCList {
	list := make([]NPC, 0)
	npcs = NPCList{
		list: list,
	}
	return &npcs
}

func (n *NPCList) AddNPC(npc NPC) {
	n.list = append(n.list, npc)
}

func (n *NPCList) NPCS() []NPC {
	return n.list
}

func (n *NPCList) GetNPCSOf(npcType NPCType) []NPC {
	npcs := make([]NPC, 0)
	for _, b := range n.list {
		if b.Type() == npcType {
			npcs = append(npcs, b)
		}
	}
	return npcs
}

func GetListInstance() *NPCList {
	return &npcs
}
