package f

import (
	"errors"
	"fmt"
	"log"
)

// Result output returned after execution of a flow, state and tasks
// TODO - convert this to an interface
type Output struct {
	Error error
	Data  interface{}
}

type Input struct {
	Data interface{}
}

// Flow encapsulates a reusable sequence of steps that can execute
// in different contexts. A flow consists of series of steps. A flow
// also has a context which can be used to maintain `state`
// TODO - consider deriving an interface ? to better represent contract
type Flow struct {
	name   string
	states []State
}

// New create a new flow.
// TODO : create and manage flows with a flow executor
func New(name string) *Flow {
	return newFlow(name)
}

// TODO: whether we need to intialize a `FlowContext` at this point ?
func newFlow(name string) *Flow {
	f := &Flow{
		name: name,
	}
	return f
}

// Validate a flow before it executes.
func (f Flow) validateFlow() error {
	// do not proceed with execution if no valid states found
	if f.states == nil || len(f.states) == 0 {
		return errors.New("No states found for flow " + f.name)
	}

	return nil
}

// AddState add states to a flow. This is needs to be called before a
// flow starts to execute
func (f *Flow) AddStates(states ...State) {
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
	}
	// add to already existing states
	f.states = append(f.states, states...)
}

//TODO: Use more idiomatic way to handle errors - so than a calling
// program can decide and handle errors in a more nicer way

//TODO: Refator flow execution to handle asynchronous
// concurrent execution of steps

// Execute a flow.
func (f *Flow) Execute(in *Input) (*Output, error) {
	err := f.validateFlow()
	// currently we only deal with sequential synchronous execution
	if err != nil {
		log.Fatal(err)
	}
	var out *Output
	for _, state := range f.states {
		// fmt.Println("Executing state         - ", state.GetName())
		// fmt.Println("Supplied data for state - ", in.Data)
		out, err = state.Execute(in)
		if err != nil {
			return nil, err
		}
		checkOutput(out, state.GetName())
		in.Data = out.Data
	}
	return out, nil
}

type Validator interface {
	Validate() error
}

//State contract to be implemented by state
type State interface {
	GetName() string
	fmt.Stringer
	// Validator - state can be validated
	Validator
	// AddTasks adds tasks to a state. It is expected that AddState(..)
	// to be called before a state begins to execute
	AddTasks(...Task)
	GetTasks() []Task
	Execute(in *Input) (*Output, error)
}

type flowState struct {
	Name  string
	Tasks []Task
}

// NewState create, initiaize and return a new flow state
//TODO: make this function to return State
func NewState(name string) State /* State */ {
	fs := &flowState{Name: name}
	return fs
}

func (fs *flowState) GetName() string {
	return fs.Name
}

func (fs flowState) String() string {
	return fmt.Sprintf("State %v  ", fs.Name)
}

func (fs flowState) Validate() error {
	// fmt.Println("Validating state - ", fs.GetName())
	return nil
}

func checkOutput(out *Output, name string) {
	if out == nil {
		errStr := fmt.Errorf(
			"Error in %s - Output is found to be nil. \n", name)
		log.Fatal(errStr)

	}
	if out.Data == nil {
		errStr := fmt.Errorf(
			"Error in %s - Output.Data  found to be nil. \n",
			name)
		log.Fatal(errStr)
	}
}

func (fs flowState) Execute(in *Input) (*Output, error) {
	if fs.Tasks == nil {
		fmt.Println("No tasks found for state ", fs.GetName())
	}

	var out *Output
	var err error
	for _, task := range fs.Tasks {
		fmt.Println("Executing tasks         - ", task.Name)
		// fmt.Println("In data  - ", in.Data)
		out, err = task.Execute(in)
		if err != nil {
			return nil, err
		}
		// fmt.Println(fmt.Sprintf("Out data - %v", out.Data))
		checkOutput(out, task.Name)
		fmt.Println("Executed tasks         - ", task.Name)
		in.Data = out.Data
	}
	return out, nil
}

type Task struct {
	Name     string
	Validate func() error
	Execute  func(in *Input) (*Output, error)
}

func (fs *flowState) AddTasks(tasks ...Task) {
	if tasks == nil {
		log.Fatal("Valid tasks not found for state - ", fs.Name)
	}
	for _, task := range tasks {
		if task.Validate != nil {
			err := task.Validate()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fs.Tasks = append(fs.Tasks, tasks...)
}
func (fs flowState) GetTasks() []Task {
	return fs.Tasks
}
