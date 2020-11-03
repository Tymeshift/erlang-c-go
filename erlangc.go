package erlangc

import (
	"math"
	"math/big"
	"sync"

	"github.com/ALTree/bigfloat"
)

func bigMul(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).Mul(x, y)
}

func bigDiv(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).Quo(x, y)
}

func bigAdd(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).Add(x, y)
}

func bigSub(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).Sub(x, y)
}

func getFactorial(n int64) *big.Float {
	var fact big.Int
	fact.MulRange(1, n)
	return new(big.Float).SetInt(&fact)
}

func getIntensity(volume *big.Float, aht *big.Float, intervalLength *big.Float) *big.Float {
	return bigMul(volume, bigDiv(aht, intervalLength))
}

func getAN(intensity *big.Float, agents *big.Float) *big.Float {
	return bigfloat.Pow(intensity, agents)
}

func getX(AN *big.Float, factorial *big.Float, intensity *big.Float, agents *big.Float) *big.Float {
	return bigMul(bigDiv(AN, factorial), bigDiv(agents, bigSub(agents, intensity)))
}

func getY(intensity *big.Float, agents *big.Float) *big.Float {
	sum := big.NewFloat(0)
	n, _ := agents.Int64()
	for i := int64(0); i < n; i++ {
		iFact := getFactorial(i)
		aPowI := bigfloat.Pow(intensity, new(big.Float).SetInt64(i))
		sum = bigAdd(sum, bigDiv(aPowI, iFact))
	}
	return sum
}

func getPW(X *big.Float, Y *big.Float) *big.Float {
	YX := Y.Add(Y, X)
	return X.Quo(X, YX)
}

func getErlangC(AN *big.Float, factorial *big.Float, intensity *big.Float, agents *big.Float) *big.Float {
	X := getX(AN, factorial, intensity, agents)
	Y := getY(intensity, agents)
	PW := getPW(X, Y)
	return PW
}

func getServiceLevel(
	erlangC *big.Float,
	intensity *big.Float,
	agents *big.Float,
	targetTime *big.Float,
	aht *big.Float,
) *big.Float {

	targetTimeToAht := bigDiv(targetTime, aht)
	agentsSubInt := bigSub(agents, intensity)
	expInput := bigMul(bigMul(agentsSubInt, targetTimeToAht), big.NewFloat(-1))

	exp := bigfloat.Exp(expInput)
	erlangCMul := bigMul(erlangC, exp)
	return new(big.Float).Sub(big.NewFloat(1), erlangCMul)
	// return (
	//   // 1 - erlangC.times(Math.exp(-(agents - intensity) * (targetTime / aht)))
	// );
}

func getFullServiceLevel(intensity *big.Float, agents *big.Float, targetTime *big.Float, aht *big.Float) *big.Float {
	n, _ := agents.Int64()
	factorial := getFactorial(n)
	AN := getAN(intensity, agents)
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

func checkMaxOccupancy(intensity *big.Float, agents *big.Float, maxOccupancy *big.Float) *big.Float {
	occupancy := bigDiv(intensity, agents)
	for occupancy.Cmp(maxOccupancy) == 0 || occupancy.Cmp(maxOccupancy) == 1 {
		agents.Add(agents, big.NewFloat(1.0))
		occupancy = bigDiv(intensity, agents)
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
	volume := new(big.Float).SetFloat64(fteParams.Volume)
	intervalLength := new(big.Float).SetInt64(fteParams.IntervalLength)
	aht := new(big.Float).SetInt64(fteParams.Aht)
	targetServiceLevel := new(big.Float).SetFloat64(fteParams.TargetServiceLevel)
	targetTime := new(big.Float).SetInt64(fteParams.TargetTime)
	maxOccupancy := new(big.Float).SetFloat64(fteParams.MaxOccupancy)

	intensity := getIntensity(volume, aht, intervalLength)
	intensityRounded, _ := new(big.Float).Add(intensity, new(big.Float).SetFloat64(0.5)).Int(nil)
	agents := new(big.Float).SetInt(intensityRounded)
	agents = agents.Add(agents, new(big.Float).SetInt64(1))

	for getFullServiceLevel(intensity, agents, targetTime, aht).Cmp(targetServiceLevel) == -1 {
		agents.Add(agents, big.NewFloat(1.0))
	}

	if fteParams.MaxOccupancy > 0 {
		agents = checkMaxOccupancy(intensity, agents, maxOccupancy)
	}

	agentsInt, _ := new(big.Float).Add(agents, new(big.Float).SetFloat64(0.5)).Int64()

	if fteParams.Shrinkage == 1 {
		fteParams.Shrinkage = 0.99999
	}
	agentsInt = int64(math.Ceil(float64(agentsInt) / (1 - fteParams.Shrinkage)))

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
