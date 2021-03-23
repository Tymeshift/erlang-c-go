package erlangc

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"testing"

	"math/big"
)

func TestIntensity(t *testing.T) {
	volume := 12.0
	intervalLength := int64(900)
	aht := int64(600)

	intensity := getIntensity(volume, aht, intervalLength)
	expected := 8.004
	if intensity != expected {
		t.Errorf("intensity should be %f, got %f", expected, intensity)
	}
}

func TestGetAN(t *testing.T) {
	res := getAN(new(big.Rat).SetFloat64(8.0), big.NewInt(10))
	expected := new(big.Rat).SetFloat64(math.Pow(8, 10))
	if res.Cmp(expected) != 0 {
		t.Errorf("AN should be %s, got %s", expected, res)
	}
}

func TestGetX(t *testing.T) {
	an := math.Pow(8, 10)
	fact := int64(3628800)
	intensity := 8.0
	agents := int64(10)
	res := getX(new(big.Rat).SetFloat64(an), big.NewInt(fact), intensity, agents)
	expected := 1479.4723
	resF, _ := res.Float64()

	if math.Round(resF*10000)/10000 != expected {
		t.Errorf("X should be %f, got %f", expected, resF)
	}
}

func TestGetY(t *testing.T) {
	intensity := 8.0
	agents := int64(10)
	res := getY(new(big.Rat).SetFloat64(intensity), agents)
	expected := 2136.2268
	resF, _ := res.Float64()
	if math.Round(resF*10000)/10000 != expected {
		t.Errorf("Y should be %f, got %f", expected, resF)
	}
}

func TestGetPW(t *testing.T) {
	X := 1479.4723
	Y := 2136.2268
	res := getPW(new(big.Rat).SetFloat64(X), new(big.Rat).SetFloat64(Y))
	expected := 0.40918
	if math.Round(res*100000)/100000 != expected {
		t.Errorf("PW should be %f, got %f", expected, res)
	}
}

func TestGetErlangC(t *testing.T) {
	an := math.Pow(8, 10)
	fact := int64(3628800)
	intensity := 8.0
	agents := 10.0
	res := getErlangC(new(big.Rat).SetFloat64(an), big.NewInt(fact), intensity, int64(agents))
	expected := 0.40918
	if math.Round(res*100000)/100000 != expected {
		t.Errorf("erlang should be %f, got %f", expected, res)
	}
}

func TestGetServiceLevel(t *testing.T) {
	erlang := 0.4091801508
	intensity := 8.0
	agents := 10
	targetTime := 1000
	aht := 1500
	res := getServiceLevel(erlang, intensity, int64(agents), int64(targetTime), int64(aht))
	expected := 0.89214
	if math.Round(res*100000)/100000 != expected {
		t.Errorf("service level should be %f, got %f", expected, res)
	}
}

func TestGetFactorial(t *testing.T) {
	res := getFactorial(5)
	expected := big.NewInt(120)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}

	res = getFactorial(20)
	expected = big.NewInt(2432902008176640000)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}

	res = getFactorial(80)
	expectedStr := "71569457046263802294811533723186532165584657342365752577109445058227039255480148842668944867280814080000000000000000000"
	expected = new(big.Int)
	expected, _ = expected.SetString(expectedStr, 10)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}
}

func TestGetFactorialSwing(t *testing.T) {
	res := getFactorialSwing(5)
	expected := big.NewInt(120)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}

	res = getFactorialSwing(20)
	expected = big.NewInt(2432902008176640000)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}

	res = getFactorialSwing(80)
	expectedStr := "71569457046263802294811533723186532165584657342365752577109445058227039255480148842668944867280814080000000000000000000"
	expected = new(big.Int)
	expected, _ = expected.SetString(expectedStr, 10)
	if res.Cmp(expected) != 0 {
		t.Errorf("factorial result should be %s, got %s", expected.String(), res.String())
	}
}

func BenchmarkGetFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFactorial(10000)
	}
}

func BenchmarkGetFactorialSwing10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFactorialSwing(10)
	}
}

func BenchmarkGetFactorialSwing100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFactorialSwing(100)
	}
}

func BenchmarkGetFactorialSwing1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFactorialSwing(1000)
	}
}

func BenchmarkGetFactorialSwing10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFactorialSwing(10000)
	}
}

func BenchmarkCaclulateFTE(b *testing.B) {
	file, err := os.Open("fteParams.json") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	var params []FteParams
	bytes, _ := io.ReadAll(file)
	json.Unmarshal(bytes, &params)
	// fmt.Println(params)
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < b.N; i++ {
		CalculateFte(params)
	}
}

func TestCalculateFte(t *testing.T) {
	volume := 0.5
	answer := int64(2)
	num := GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 1
	answer = int64(3)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       1,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 1
	answer = int64(3)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 2
	answer = int64(3)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 10.0
	answer = int64(8)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 50
	answer = int64(27)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 100
	answer = int64(53)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 200
	answer = int64(105)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 500
	answer = int64(262)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 5000
	answer = int64(2605)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}

	volume = 0.046481566
	answer = int64(6)
	num = GetNumberOfAgents(FteParams{
		ID:                 "1",
		Index:              0,
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.5,
		Aht:                29400,
		TargetServiceLevel: 0.8,
		TargetTime:         14400,
	})

	if num.Volume != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num.Volume, answer)
	}
}
