package divisors

import (
	"fmt"
	"log"
)

//Element class
type NumDivElement struct {
	Num    int64
	NumDiv int
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	elements     []*NumDivElement
	elementCount int
}

// String method on Element class
func (element *NumDivElement) String() string {
	return fmt.Sprintf("%d, %d", element.Num, element.NumDiv)
}

// NewStack returns a new stack.
func (stack *Stack) New() {
	stack.elements = make([]*NumDivElement, 0)
}

func NewStack() *Stack {
	stack := &Stack{}
	stack.New()
	return stack
}

// Push adds a node to the stack.
func (stack *Stack) Push(element *NumDivElement) (int, *NumDivElement) {
	stack.elements = append(stack.elements, element)
	stack.elementCount = stack.elementCount + 1
	return stack.elementCount, element
}

// CondPush adds a node to the stack conditionally
func (stack *Stack) CondPush(element *NumDivElement) (int, *NumDivElement) {
	last := stack.elementCount - 1
	if stack.elementCount < 1 {
		return stack.Push(element)
	} else if stack.elements[last].Num > element.Num && stack.elements[last].NumDiv <= element.NumDiv {
		log.Printf("Remove %v", stack.Pop())
		return stack.Push(element)
	} else if stack.elements[last].Num < element.Num && stack.elements[last].NumDiv >= element.NumDiv {
		// Discard this element
		log.Printf("Discard %v", element)
		return -1, nil
	} else if stack.elements[last].Num < element.Num && stack.elements[last].NumDiv < element.NumDiv {
		return stack.Push(element)
	} else {
		head := stack.Pop()
		ok, _ := stack.CondPush(element)
		if ok != -1 {
			return stack.Push(head)
		} else {
			log.Printf("HRemove %v", head)
			return -1, nil
		}

	}
}

// Pop removes and returns a node from the stack in last to first order.
func (stack *Stack) Pop() *NumDivElement {
	if stack.elementCount == 0 {
		return nil
	}
	var length int = len(stack.elements)
	var element *NumDivElement = stack.elements[length-1]
	//stack.elementCount = stack.elementCount - 1
	stack.elements = stack.elements[:length-1]
	stack.elementCount = len(stack.elements)
	return element
}
