package concurrency

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
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

func TestAtomicBool_Wait(t *testing.T) {
	ab := NewAtomicBool()
	var readyGroup sync.WaitGroup
	var doneGroup sync.WaitGroup

	worker := func(count int, results []bool) {
		readyGroup.Add(count)
		doneGroup.Add(count)
		for i := 1; i <= count; i++ {
			go func(num int) {
				readyGroup.Done()
				fmt.Printf("go routine %d, launched\n", num)
				val := ab.Wait()
				results[num-1] = val
				fmt.Printf("go routine %d received true\n", num)
				doneGroup.Done()
			}(i)
		}
	}
	count := 2
	results := make([]bool, count)
	worker(2, results)
	readyGroup.Wait()

	fmt.Println("setting true")
	ab.Set(true)
	fmt.Println("set")
	doneGroup.Wait()

	for _, r := range results {
		if !r {
			t.Errorf("expected all results to be true, found false")
		}
	}
}

func TestAtomicBool_WaitWithContext(t *testing.T) {

	var readyGroup sync.WaitGroup
	var doneGroup sync.WaitGroup

	worker := func(count int, ab *AtomicBool, ctx context.Context, results []bool) {
		readyGroup.Add(count)
		doneGroup.Add(count)
		for i := 1; i <= count; i++ {
			go func(num int) {
				readyGroup.Done()
				fmt.Printf("go routine %d, launched\n", num)
				val := ab.WaitWithContext(ctx)
				results[num-1] = val
				fmt.Printf("go routine %d received %v\n", num, val)
				doneGroup.Done()
			}(i)
		}
	}

	t.Run("clears wait when set to true", func(t *testing.T) {
		ab := NewAtomicBool()

		ab.Set(false)
		count := 2
		results := make([]bool, count)
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		worker(2, ab, ctx, results)
		readyGroup.Wait()

		fmt.Println("setting true")
		ab.Set(true)
		fmt.Println("set")
		doneGroup.Wait()

		for _, r := range results {
			if !r {
				t.Errorf("expected all results to be true, found false")
			}
		}
	})

	t.Run("clears early when context is done", func(t *testing.T) {
		ab := NewAtomicBool()
		count := 2
		results := make([]bool, count)
		ctx, cancel := context.WithCancel(context.Background())

		ab.Set(false)
		worker(2, ab, ctx, results)
		readyGroup.Wait()

		time.Sleep(500 * time.Millisecond)
		fmt.Println("cancelling context...")
		cancel()
		fmt.Println("... cancelled")
		doneGroup.Wait()

		for _, r := range results {
			if r {
				t.Errorf("expected all results to be false due to cancellation, found true")
			}
		}
	})
}
