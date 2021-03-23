package erlangc

import (
	"math"
	"sync"

	big "github.com/ncw/gmp"

	"github.com/Tymeshift/erlang-c-go/factorial"
)

var factorailCache map[int64]*big.Int = make(map[int64]*big.Int)

func ratioExp(x *big.Rat, y *big.Int) *big.Rat {
	num := x.Num()
	num = num.Exp(num, y, nil)
	denom := x.Denom()
	denom = denom.Exp(denom, y, nil)
	return new(big.Rat).SetFrac(num, denom)
}

func getFactorial(n int64) big.Int {
	var fact big.Int
	fact.MulRange(1, n)
	return fact
}

func getFactorialSwing(n int64) *big.Int {
	cache, ok := factorailCache[n]
	if ok {
		return cache
	}
	fact := factorial.Factorial(uint64(n))
	factorailCache[n] = fact
	return fact
}

func getIntensity(volume float64, aht int64, intervalLength int64) float64 {
	return volume * math.Round((float64(aht)/float64(intervalLength))*100) / 100
}

func getAN(intensity *big.Rat, agents *big.Int) *big.Rat {
	res := ratioExp(intensity, agents)
	return res
}

func getX(AN *big.Rat, factorial *big.Int, intensity float64, agents int64) *big.Rat {
	agentsCoeff := math.Round((float64(agents)/(float64(agents)-intensity))*100) / 100
	res := new(big.Rat).Quo(AN, new(big.Rat).SetInt(factorial))
	return new(big.Rat).Mul(res, new(big.Rat).SetFloat64(agentsCoeff))
}

func getY(intensity *big.Rat, agents int64) *big.Rat {
	sum := new(big.Rat)
	for i := int64(0); i < agents; i++ {
		iFact := getFactorialSwing(i)
		aPowI := ratioExp(intensity, big.NewInt(i))
		div := new(big.Rat).Quo(aPowI, new(big.Rat).SetInt(iFact))
		sum = new(big.Rat).Add(sum, div)
	}
	return sum
}

func getPW(X *big.Rat, Y *big.Rat) *big.Rat {
	YX := Y.Add(Y, X)
	return X.Quo(X, YX)
}

func getErlangC(AN *big.Rat, factorial *big.Int, intensity float64, agents int64) *big.Rat {
	X := getX(AN, factorial, intensity, agents)
	Y := getY(new(big.Rat).SetFloat64(intensity), agents)
	PW := getPW(X, Y)
	return PW
}

func getServiceLevel(
	erlangC *big.Rat,
	intensity float64,
	agents int64,
	targetTime int64,
	aht int64,
) float64 {
	targetTimeToAht := float64(targetTime) / float64(aht)
	agentsSubInt := float64(agents) - intensity
	expInput := float64(agentsSubInt) * targetTimeToAht * -1
	exp := math.Round(math.Exp(expInput)*100) / 100
	erlangCMul := erlangC.Mul(erlangC, new(big.Rat).SetFloat64(exp))
	res, _ := erlangCMul.Float64()
	return 1 - res
	// return (
	//   // 1 - erlangC.times(Math.exp(-(agents - intensity) * (targetTime / aht)))
	// );
}

func getFullServiceLevel(intensity float64, agents int64, targetTime int64, aht int64) float64 {
	factorial := getFactorialSwing(agents)
	bigInensity := new(big.Rat).SetFloat64(intensity)
	AN := getAN(bigInensity, big.NewInt(agents))
	erlangC := getErlangC(AN, factorial, intensity, agents)
	serviceLevel := getServiceLevel(
		erlangC,
		intensity,
		agents,
		targetTime,
		aht,
	)
	return serviceLevel
}

func checkMaxOccupancy(intensity float64, agents int64, maxOccupancy float64) int64 {
	occupancy := intensity / float64(agents)
	for maxOccupancy >= occupancy {
		agents++
		occupancy = intensity / float64(agents)
	}
	return agents
}

// FteParams - parameters to calculate FTE
type FteParams struct {
	ID                 string
	Index              int64
	Volume             float64
	IntervalLength     int64
	Aht                int64
	TargetServiceLevel float64
	TargetTime         int64
	MaxOccupancy       float64
	Shrinkage          float64
}

type FteResult struct {
	ID     string
	Index  int64
	Volume int64
}

func GetNumberOfAgents(fteParams FteParams) FteResult {
	if fteParams.Volume < 0 || fteParams.Aht < 0 {
		return FteResult{
			ID:     fteParams.ID,
			Index:  fteParams.Index,
			Volume: 2,
		}
	}

	intensity := math.Round(getIntensity(fteParams.Volume, fteParams.Aht, fteParams.IntervalLength)*100) / 100
	agents := int64(math.Floor(intensity + 1))

	// s := getFullServiceLevel(intensity, agents, fteParams.TargetTime, fteParams.Aht)
	// fmt.Println(s)

	s := 0.0

	for s < fteParams.TargetServiceLevel {
		s = getFullServiceLevel(intensity, agents, fteParams.TargetTime, fteParams.Aht)
		agents++
	}

	// if fteParams.MaxOccupancy > 0 {
	// 	agents = checkMaxOccupancy(intensity, agents, fteParams.MaxOccupancy)
	// }

	if fteParams.Shrinkage == 1 {
		fteParams.Shrinkage = 0.99999
	}
	agentsInt := int64(math.Ceil(float64(agents) / (1 - fteParams.Shrinkage)))

	return FteResult{
		ID:     fteParams.ID,
		Index:  fteParams.Index,
		Volume: agentsInt,
	}
}

func getNumberOfAgentsParallel(fteParams FteParams, fteChan chan FteResult, wg *sync.WaitGroup) {
	agents := GetNumberOfAgents(fteParams)
	fteChan <- agents
	wg.Done()
}

// CalculateFte calculats number of agents needed for a specific service level to handle incoming volume of arrivals per time interval
//
// volume - incoming number of arrivals per time interval
// intervalLength - time interval in seconds
// aht - average handle time in seconds
// targetServiceLevel - service level goal, the percentage of calls answered within the acceptable waiting time (0 <= targetServiceLevel < 1)
// targetTime - target answer time, acceptable wait time in seconds
// maxOccupancy - maximum occupancy rate (0 <= maxOccupancy <= 1)
// shrinkage - shrinkage rate (0 <= shrinkage < 1)
func CalculateFte(params []FteParams) []FteResult {
	fte := make([]FteResult, len(params))
	for i, param := range params {
		fte[i] = GetNumberOfAgents(param)
	}

	return fte
}

// CalculateFteParallel calculats number of agents needed for a specific service level to handle incoming volume of arrivals per time interval
//
// volume - incoming number of arrivals per time interval
// intervalLength - time interval in seconds
// aht - average handle time in seconds
// targetServiceLevel - service level goal, the percentage of calls answered within the acceptable waiting time (0 <= targetServiceLevel < 1)
// targetTime - target answer time, acceptable wait time in seconds
// maxOccupancy - maximum occupancy rate (0 <= maxOccupancy <= 1)
// shrinkage - shrinkage rate (0 <= shrinkage < 1)
func CalculateFteParallel(params []FteParams) []FteResult {
	var fte []FteResult
	fteChan := make(chan FteResult, len(params))
	wg := sync.WaitGroup{}
	for _, param := range params {
		wg.Add(1)
		go getNumberOfAgentsParallel(param, fteChan, &wg)
	}
	wg.Wait()

	close(fteChan)

	for agents := range fteChan {
		fte = append(fte, agents)
	}

	return fte
}
