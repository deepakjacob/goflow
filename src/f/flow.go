package f

import (
	"errors"
	"fmt"
	"log"
)

// FlowContext context associated with the flow
type FlowContext struct {
	Ctx map[string]interface{}
}

// Result output returned after execution of a flow, state and tasks
// TODO - convert this to an interface
type Result struct {
}

// Flow encapsulates a reusable sequence of steps that can execute
// in different contexts. A flow consists of series of steps. A flow
// also has a context which can be used to maintain `state`
// TODO - consider deriving an interface ? to better represent contract
type flow struct {
	name           string
	ctx            *FlowContext
	states         []State
	proceedOnError bool
}

// New create a new flow.
// TODO : create and manage flows with a flow executor
func New(name string, proceedOnError bool) *flow {
	return newFlow(name, proceedOnError)
}

// TODO: whether we need to intialize a `FlowContext` at this point ?
func newFlow(name string, proceedOnError bool) *flow {
	ctx := &FlowContext{Ctx: make(map[string]interface{})}
	f := &flow{
		name:           name,
		ctx:            ctx,
		proceedOnError: proceedOnError,
	}
	return f
}

// an internal  method called to check whether flow related state
// variables are initialized properly before starting a flow .Do not
// confuse with state's own Validate method the above method is only
// for validating user created state variables
func (f flow) validateFlowState(state State) error {
	return nil
}

// Validate a flow before it executes.
func (f flow) validateFlow() error {
	// do not proceed with execution if no valid states found
	if f.states == nil || len(f.states) == 0 {
		return errors.New("No states found for flow " + f.name)
	}

	return nil
}

// AddState add states to a flow. This is needs to be called before a
// flow starts to execute
func (f *flow) AddStates(states ...State) {
	if states == nil {
		log.Fatal("Valid states not found for flow - ", f.name)
	}
	for _, state := range states {
		// invokes state's own validate in this validate we
		// dont validate name and ctx variables
		err := state.Validate()
		if err != nil {
			log.Fatal(err)
		}
		// check whether flow related state variables are properly
		// initialized
		err = f.validateFlowState(state)
		if err != nil {
			log.Fatal(err)
		}
	}
	// add to already existing states
	f.states = append(f.states, states...)
}

//TODO: Use more idiomatic way to handle errors - so than a calling
// program can decide and handle errors in a more nicer way

//TODO: Refator flow execution to handle asynchronous
// concurrent execution of steps

// Execute a flow. Having an associated context and non nil states are
// mandatory. If `ProceedOnError` is set to false then entire flow is
// stopped even when a single step return a non nil error.
func (f *flow) Execute() {
	err := f.validateFlow()
	// currently we only deal with sequential synchronous execution
	if err != nil {
		log.Fatal(err)
	}
	for _, state := range f.states {
		err, result := state.Execute(f.ctx)
		if err != nil {
			if f.proceedOnError {
				fmt.Println(
					"Error: but continue executing ",
					state,
				)
				continue
			} else {
				log.Fatal(err)
				return
			}
		} else {
			fmt.Println("Result %v", result)
		}
	}
}

// StateContext a state with in a flow. It has a name and a `context`
// where it maintains it's own internal state
type stateContext struct {
	ctx map[string]interface{}
}

//State contract tobe implemented by state
type State interface {
	fmt.Stringer
	Validate() error
	Execute(ctx *FlowContext) (*Result, error)
}

type FlowState struct {
	Name  string
	Ctx   *stateContext
	Tasks []Task
}

// NewState create, initiaize and return a new flow state
//TODO: make this function to return State
func NewState(name string) *FlowState /* State */ {
	fs := &FlowState{Name: name}
	fs.Ctx = &stateContext{ctx: make(map[string]interface{})}
	return fs
}

func (fs *FlowState) GetName() string {
	return fs.Name
}

func (fs FlowState) String() string {
	return fmt.Sprintf("State %v  ", fs.Name)
}

func (fs FlowState) Validate() error {
	fmt.Println("Validating state - ", fs.GetName())
	return nil
}

func (fs FlowState) Execute(ctx *FlowContext) (*Result, error) {
	r := &Result{}
	// TODO: Handle this
	fmt.Println("Executing state - ", fs.GetName())
	if fs.Tasks == nil {
		fmt.Println("No tasks found for state ", fs.GetName())
	}
	for _, task := range fs.Tasks {
		err, result := task.Run()
		if err != nil {
			fmt.Println("Result %v", result)
		}
	}
	return r, nil
}

func (fs FlowState) validateTask(task Task) error {
	return nil
}

type TaskContext struct {
	Ctx map[string]interface{}
}

type Task struct {
	Name string
	Run  func() (*Result, error)
	Ctx  *TaskContext
}

func (fs *FlowState) AddTasks(tasks ...Task) {
	if tasks == nil {
		log.Fatal("Valid tasks not found for state - ", fs.Name)
	}
	for _, task := range tasks {
		// check whether flow related state variables are
		// properly initialized
		err := fs.validateTask(task)
		if err != nil {
			log.Fatal(err)
		}
	}
	// add to already existing states
	fs.Tasks = append(fs.Tasks, tasks...)
}
