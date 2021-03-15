package divisors

import (
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Generate array of primes
func genArrayOfPrimes(n int64) []uint {

	primes := []uint{2}
	top := uint(math.Ceil(math.Sqrt(float64(n))))
	for i := uint(3); i < top; i += 2 {
		isPrime := true
		for j := 0; j < len(primes); j++ {
			if i%primes[j] == 0 {
				isPrime = false
			}
		}
		if isPrime {
			primes = append(primes, i)
		}

	}

	return primes
}

var (
	cachedPrimes   = []uint{}
	gemArrayMutex  sync.Mutex
	AtomicDiscards = int64(0)
)

func initCachePrimes(n int64) {
	gemArrayMutex.Lock()
	cachedPrimes = genArrayOfPrimes(n)
	gemArrayMutex.Unlock()
}

// PrimeFactorization Prime factorization of a given number
//
// A map is returned.
//
//   key of map: prime
//   value of map: prime exponents
//
func PrimeFactorization(n int64) (pfs map[uint]int) {
	pfs = make(map[uint]int)
	if len(cachedPrimes) == 0 {
		initCachePrimes(n)
	}
	// log.Print("cachedPrimes: ", cachedPrimes)
	// Get the number of 2s that divide n
	for n%2 == 0 {
		if _, ok := pfs[2]; ok {
			pfs[2]++
		} else {
			pfs[2] = 1
		}
		n = n / 2
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	for i := 1; i < len(cachedPrimes) && int64(cachedPrimes[i])*int64(cachedPrimes[i]) <= n; i++ {
		// while i divides n, append i and divide n
		for n%int64(cachedPrimes[i]) == 0 {
			if _, ok := pfs[cachedPrimes[i]]; ok {
				pfs[cachedPrimes[i]]++
			} else {
				pfs[cachedPrimes[i]] = 1
			}
			n = n / int64(cachedPrimes[i])
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n > 2 {
		pfs[uint(n)] = 1
	}

	return
}

// NumberOfDivisors calculate number of divisors of a given number
func NumberOfDivisors(n int64) (numd int, pfs map[uint]int) {
	pfs = PrimeFactorization(n)

	numd = 1
	for _, exponents := range pfs {
		numd *= (exponents + 1)
	}

	return
}

// FutureNumDiv result
type FutureNumDiv chan *NumDivElement

//Element class
type NumDivElement struct {
	Num    int64
	NumDiv int
}

/* AsyncNumberOfDivisors version NumberOfDivisors
func AsyncNumberOfDivisors(n int64, fut FutureNumDiv) {

	go func() { fut <- NumberOfDivisors(n) }()
} */

// HighlyComposite Calculate Highly Composite Numbers
func HighlyComposite(top int64, handler func(max int64, numd int)) int64 {
	cachedPrimes = genArrayOfPrimes(top + 1)
	topDiv := 0
	max := int64(0)
	for n := int64(2); n <= top; n += 2 {
		nd, _ := NumberOfDivisors(n)
		if nd > topDiv {
			topDiv = nd
			max = n
			handler(max, nd)
		}
	}
	return max
}

var numCores = 1
var waitGroup sync.WaitGroup

func init() {
	numCores = runtime.NumCPU() * 4
}

func checkPGapAndExponents(pfs map[uint]int) bool {
	foundPs := make([]int, 0)
	for _, p := range cachedPrimes {
		exp, ok := pfs[p]
		if !ok {
			break
		}
		foundPs = append(foundPs, exp)
	}
	ok := len(foundPs) == len(pfs)
	if ok {
		lastExp := -1
		for _, e := range foundPs {
			if lastExp != -1 && !(lastExp >= e) {
				ok = false
				break
			}
			lastExp = e
		}
	}
	return ok
}

func woker(start, step, end int64, results FutureNumDiv) func() {
	waitGroup.Add(1)
	return func() {
		go func() {
			defer waitGroup.Done()
			for n := int64(start); n <= end; n += step {
				nd, pfs := NumberOfDivisors(n)
				if checkPGapAndExponents(pfs) {
					res := NumDivElement{NumDiv: nd, Num: n}
					results <- &res
				} else {
					atomic.AddInt64(&AtomicDiscards, int64(1))
				}
			}
		}()
	}
}

type divMap map[int]NumDivElement

// ParallelHighlyComposite Calculate Highly Composite Numbers with parallelism
func ParallelHighlyComposite(top int64, handler func(idx int, max int64, numd int)) int64 {
	initCachePrimes(top)
	dmap := make(divMap)
	fchan := make(FutureNumDiv, numCores)

	var wchan sync.WaitGroup

	// var workers []func()
	for w := 0; w < numCores; w++ {
		// log.Printf("woker(int64(2+2*%d), int64(2*%d), %d, fchan)()", w, numCores, top)
		woker(int64(2+2*w), int64(2*numCores), top, fchan)()
	}

	wchan.Add(1)
	go func() {
		defer wchan.Done()
		for n := range fchan {
			currv, exists := dmap[n.NumDiv]
			if !exists || n.Num < currv.Num {
				dmap[n.NumDiv] = *n
				//	log.Println(n)
			}
		}

	}()

	time.Sleep(200 * time.Millisecond)
	// Wait for workers end
	waitGroup.Wait()
	close(fchan)
	// Wait for consumer end
	wchan.Wait()
	keys := make([]NumDivElement, len(dmap))
	idx := 0
	for _, ele := range dmap {
		keys[idx] = ele
		idx++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Num < keys[j].Num })
	// alen := len(keys)
	idx = 2
	maxDiv := 0
	for i := range keys {
		ele := keys[i]
		if ele.NumDiv > maxDiv {
			handler(idx, ele.Num, ele.NumDiv)
			idx++
			maxDiv = ele.NumDiv
		}

	}
	return dmap[maxDiv].Num
}
