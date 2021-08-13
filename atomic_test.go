package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestAtomicBool_Set(t *testing.T) {
	ab := NewAtomicBool()

	if ab.Value() != false {
		t.Error("expected value to start false")
	}

	ab.Set(true)
	if ab.Value() == false {
		t.Error("expected value to change to true, still false")
	}
}

func TestAtomicBool_WaitForTrue(t *testing.T) {
	ab := NewAtomicBool()
	var readyGroup sync.WaitGroup
	var doneGroup sync.WaitGroup

	worker := func(count int) {
		readyGroup.Add(count)
		doneGroup.Add(count)
		for i := 1; i <= count; i++ {
			go func(num int) {
				readyGroup.Done()
				fmt.Printf("go routine %d, launched\n", num)
				ab.WaitForTrue()
				fmt.Printf("go routine %d received true\n", num)
				doneGroup.Done()
			}(i)
		}
	}
	worker(2)
	readyGroup.Wait()

	fmt.Println("setting true")
	ab.Set(true)
	fmt.Println("set")
	doneGroup.Wait()
}
