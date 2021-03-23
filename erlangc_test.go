package erlangc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"testing"

	big "github.com/ncw/gmp"
)

func TestIntensity(t *testing.T) {
	volume := 12.0
	intervalLength := int64(900)
	aht := int64(600)

	intensity := getIntensity(volume, aht, intervalLength)
	expected := 8.000000
	if intensity != expected {
		t.Errorf("intensity should be %f, got %f", expected, intensity)
	}
}

func TestAN(t *testing.T) {
	intensity := big.NewRat(1, 2)
	res := getAN(intensity, big.NewInt(3))
	t.Error(res)
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
	t.Error(res)
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

	res := CalculateFte(params)
	fmt.Println(res)
	fmt.Println(len(res))
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
