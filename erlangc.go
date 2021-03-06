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
	num = new(big.Int).Exp(num, y, nil)
	denom := x.Denom()
	denom = new(big.Int).Exp(denom, y, nil)
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
	factStr := fact.String()
	factInt, _ := new(big.Int).SetString(factStr, 10)
	factorailCache[n] = factInt
	return factInt
}

func getIntensity(volume float64, aht int64, intervalLength int64) float64 {
	return volume * math.Round((float64(aht)/float64(intervalLength))*10000) / 10000
}

func getAN(intensity *big.Rat, agents *big.Int) *big.Rat {
	res := ratioExp(intensity, agents)
	return res
}

func getX(AN *big.Rat, factorial *big.Int, intensity float64, agents int64) *big.Rat {
	agentsCoeff := math.Round((float64(agents)/(float64(agents)-intensity))*10000) / 10000
	res := new(big.Rat).Quo(AN, new(big.Rat).SetInt(factorial))
	return new(big.Rat).Mul(res, new(big.Rat).SetFloat64(agentsCoeff))
}

func getY(intensity *big.Rat, agents int64) *big.Rat {
	sum := new(big.Rat)
	for i := int64(0); i < agents; i++ {
		iFact := getFactorialSwing(i)
		aPowI := ratioExp(intensity, big.NewInt(i))
		div := new(big.Rat).Quo(aPowI, new(big.Rat).SetInt(iFact))
		sum = div.Add(sum, div)
	}
	return sum
}

func getPW(X *big.Rat, Y *big.Rat) float64 {
	YX := Y.Add(Y, X)
	res, _ := X.Quo(X, YX).Float64()
	return res
}

func getErlangC(AN *big.Rat, factorial *big.Int, intensity float64, agents int64) float64 {
	X := getX(AN, factorial, intensity, agents)
	Y := getY(new(big.Rat).SetFloat64(intensity), agents)
	PW := getPW(X, Y)
	return PW
}

func getServiceLevel(
	erlangC float64,
	intensity float64,
	agents int64,
	targetTime int64,
	aht int64,
) float64 {
	targetTimeToAht := float64(targetTime) / float64(aht)
	agentsSubInt := float64(agents) - intensity
	expInput := float64(agentsSubInt) * targetTimeToAht * -1
	exp := math.Exp(expInput)
	erlangCMul := erlangC * exp
	return 1 - erlangCMul
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
	for occupancy >= maxOccupancy {
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

	intensity := getIntensity(fteParams.Volume, fteParams.Aht, fteParams.IntervalLength)
	agents := int64(math.Floor(intensity + 1))

	for getFullServiceLevel(intensity, agents, fteParams.TargetTime, fteParams.Aht) < fteParams.TargetServiceLevel {
		agents++
	}

	if fteParams.MaxOccupancy > 0 {
		agents = checkMaxOccupancy(intensity, agents, fteParams.MaxOccupancy)
	}

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
