package erlangc

import (
	"testing"
)

func TestCalculateFte(t *testing.T) {

	volume := 0.5
	answer := int64(2)
	num := GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 1
	answer = int64(3)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 2
	answer = int64(3)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 10
	answer = int64(8)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 50
	answer = int64(27)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 100
	answer = int64(53)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 200
	answer = int64(105)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

	volume = 500
	answer = int64(262)
	num = GetNumberOfAgents(FteParams{
		Volume:             volume,
		IntervalLength:     900,
		MaxOccupancy:       0.8,
		Shrinkage:          0.2,
		Aht:                300,
		TargetServiceLevel: 0.8,
		TargetTime:         60,
	})

	if num != answer {
		t.Errorf("CalculateFte with %f volume = %d; want %d", volume, num, answer)
	}

}
