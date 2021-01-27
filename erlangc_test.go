package erlangc

import (
	"math/big"
	"testing"
)

func TestIntensity(t *testing.T) {
	volume := new(big.Float).SetFloat64(12)
	intervalLength := new(big.Float).SetInt64(900)
	aht := new(big.Float).SetInt64(600)

	intensity := getIntensity(volume, aht, intervalLength)
	expected := 8.000000
	if intensity.Cmp(big.NewFloat(expected)) != 0 {
		t.Errorf("intensity should be %f, got %f", expected, intensity)
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

	volume = 10
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
