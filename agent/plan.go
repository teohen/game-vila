package agent

type Plan struct {
	goal       IGoal
	actions    []IAction
	currAction int
	cost       float32
}

func NewPlan() *Plan {
	p := Plan{
		currAction: 0,
	}
	return &p
}

func (p *Plan) Clear() {
	p.actions = nil
	p.goal = nil
	p.currAction = 0
}
func (p *Plan) SetGoal(g IGoal) {
	p.goal = g
}

func (p *Plan) SetActions(actions []IAction) {
	p.actions = actions
}

func (p *Plan) IsSet() bool {
	return len(p.actions) > 0
}

func (p *Plan) GetCurrentAction() IAction {
	return p.actions[p.currAction]
}

func (p *Plan) nextAction() {
	p.currAction += 1
}

func (p *Plan) hasAction() bool {
	return p.currAction < len(p.actions)
}

func (p *Plan) GetGoal() IGoal {
	return p.goal
}

func (p *Plan) IncrementCost(value float32) {
	p.cost += value
}
