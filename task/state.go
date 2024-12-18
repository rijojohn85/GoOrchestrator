package task

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if state == s {
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dest State) bool {
	return Contains(stateTransitionMap[src], dest)
}
