package f

type FlowDefintion interface {
	New(name string, concurrent bool) Flow
	AddStates(states ...State)
}

type FlowExecutor interface {
	//Returns an error if name is already present in the register.
	ExecuteFlow(name string) (*Output, error)
}

type FlowRegistry interface {
	GetFlow(name string) (Flow, error)
	RegisterFlow(f Flow, name string) error
}

type registry struct {
	// internal structure used by the FlowExecutor for maintaining
	// state.
	flows map[string]Flow
}
