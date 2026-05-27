package entity

type TraitType int

const (
	TraitMove TraitType = iota
	TraitChop
	TraitDeposit
)

type PlanStep struct {
	Trait   TraitType
	TargetX int
	TargetY int
	Tree    *Tree
}
