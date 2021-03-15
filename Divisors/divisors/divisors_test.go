package divisors_test

import (
	divisors "mathalgol/Divisors/divisors"
	"testing"
)

func TestDivisors(t *testing.T) {
	num := int64(75600)
	nd, _ := divisors.NumberOfDivisors(num)
	if nd != 120 {
		t.Fatal("Expected 120 ")
	}
	nd, _ = divisors.NumberOfDivisors(int64(6983776800))
	if nd != 2304 {
		t.Fatal("Expected 2304 ", nd)
	}
}

func handler(max int64, numd int) {

}

func parHandler(idx int, max int64, numd int) {
	//	fmt.Printf("%3d -> %d, %d div.\n", idx, max, numd)
}

func TestHighlyComposite(t *testing.T) {
	res := divisors.HighlyComposite(100, handler)
	if res != int64(60) {
		t.Fatal("Expected 60 ", res)
	}
	res = divisors.HighlyComposite(5544, handler)
	if res != int64(5040) {
		t.Fatal("Expected 5040", res)
	}
}

/* func TestStack(t *testing.T) {
	stack := divisors.NewStack()
	stack.Push(&divisors.NumDivElement{1, 1})
	stack.Pop()
	if stack.Pop() != nil {
		t.Fatal("Nil expected")
	}
	stack.CondPush(&divisors.NumDivElement{24, divisors.NumberOfDivisors(24)})
	stack.CondPush(&divisors.NumDivElement{1, 1})
	stack.CondPush(&divisors.NumDivElement{25, divisors.NumberOfDivisors(25)})
	head := stack.Pop()
	if head.Num != int64(24) {
		t.Fatal("Expected 24", head.Num)
	}
}*/

func TestParHighlyComposite(t *testing.T) {
	res := divisors.ParallelHighlyComposite(100, parHandler)
	if res != int64(60) {
		t.Fatal("Expected 60: ", res)
	}
	res = divisors.ParallelHighlyComposite(5544, parHandler)
	if res != int64(5040) {
		t.Fatal("Expected 5040")
	}
}
