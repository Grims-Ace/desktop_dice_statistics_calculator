package main

import (
	"testing"
)

func TestCalculateDiceStatistics_Multiplication(t *testing.T) {
	expression := "5d10*3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Errorf("CalculateDiceStatistics(%q) failed: %v", expression, err)
		return
	}

	// For 5d10, min is 5, max is 50.
	// With *3, min should be 15, max should be 150.
	if stats.MinValue != 15 {
		t.Errorf("Expected MinValue 15, got %d", stats.MinValue)
	}
	if stats.MaxValue != 150 {
		t.Errorf("Expected MaxValue 150, got %d", stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_Basic(t *testing.T) {
	// 2d6 -> 2-12
	expression := "2d6"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 2 || stats.MaxValue != 12 {
		t.Errorf("Expected 2-12, got %d-%d", stats.MinValue, stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_Constant(t *testing.T) {
	// 5 + 3 -> 8
	expression := "5+3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 8 || stats.MaxValue != 8 {
		t.Errorf("Expected 8-8, got %d-%d", stats.MinValue, stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_Mixed_ConstantMult(t *testing.T) {
	// 1d4 + 2 * 3 -> 1d4 + 6 -> 7-10
	expression := "1d4+2*3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 7 || stats.MaxValue != 10 {
		t.Errorf("Expected 7-10, got %d-%d", stats.MinValue, stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_Mixed_DiceMult(t *testing.T) {
	// 1d4*2 + 3 -> (1..4)*2 + 3 -> {2,4,6,8} + 3 -> {5,7,9,11}
	expression := "1d4*2+3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 5 {
		t.Errorf("Expected MinValue 5, got %d", stats.MinValue)
	}
	if stats.MaxValue != 11 {
		t.Errorf("Expected MaxValue 11, got %d", stats.MaxValue)
	}
	// Check that 6 is NOT a possible outcome (since outcomes are 5, 7, 9, 11)
	if _, exists := stats.Results[6]; exists {
		t.Errorf("Did not expect outcome 6 to exist")
	}
}

func TestCalculateDiceStatistics_ParenthesisAndDiceMult(t *testing.T) {
	// (5d10+3)*d10
	// 5d10 ranges 5-50. +3 ranges 8-53.
	// d10 ranges 1-10.
	// Min: 8 * 1 = 8.
	// Max: 53 * 10 = 530.
	expression := "(5d10+3)*d10"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed to calculate %s: %v", expression, err)
	}
	if stats.MinValue != 8 {
		t.Errorf("Expected MinValue 8, got %d", stats.MinValue)
	}
	if stats.MaxValue != 530 {
		t.Errorf("Expected MaxValue 530, got %d", stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_ComplexParentheses(t *testing.T) {
	// 2 * (1d4 + 1) -> 2 * {2,3,4,5} -> {4,6,8,10}
	expression := "2 * (1d4 + 1)"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 4 {
		t.Errorf("Expected MinValue 4, got %d", stats.MinValue)
	}
	if stats.MaxValue != 10 {
		t.Errorf("Expected MaxValue 10, got %d", stats.MaxValue)
	}
	if _, exists := stats.Results[5]; exists {
		t.Errorf("Outcome 5 should not exist")
	}
}

func TestCalculateDiceStatistics_ImplicitMult_Dice(t *testing.T) {
	// d203d10 -> 1d203 * 1d10
	// Min: 1 * 1 = 1
	// Max: 203 * 10 = 2030
	expression := "d203d10"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed to calculate %s: %v", expression, err)
	}
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 2030 {
		t.Errorf("Expected MaxValue 2030, got %d", stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_ImplicitMult_ConstantDice(t *testing.T) {
	// 2d10 -> 2 * d10 -> 2,4,6...20?
	// NO! 2d10 should be parsed as "2 dice of 10 sides".
	// The parser MUST prioritize dice notation over implicit multiplication.
	expression := "2d10"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 2d10 is sum of 2 dice. Min 2, Max 20. All values 2..20 possible.
	if stats.MinValue != 2 {
		t.Errorf("Expected MinValue 2, got %d", stats.MinValue)
	}
	if stats.MaxValue != 20 {
		t.Errorf("Expected MaxValue 20, got %d", stats.MaxValue)
	}
	if _, exists := stats.Results[3]; !exists {
		t.Errorf("Expected outcome 3 to exist for 2d10")
	}
}

func TestCalculateDiceStatistics_ImplicitMult_NumberDice(t *testing.T) {
	// 3 d10 -> 3 * d10
	expression := "3 d10"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 3 {
		t.Errorf("Expected MinValue 3, got %d", stats.MinValue)
	}
	if stats.MaxValue != 30 {
		t.Errorf("Expected MaxValue 30, got %d", stats.MaxValue)
	}
	if _, exists := stats.Results[4]; exists {
		t.Errorf("Did not expect outcome 4 to exist for 3 * d10")
	}
}

func TestCalculateDiceStatistics_Division(t *testing.T) {
	// 1d6 / 2
	// Outcomes: 1/2=0, 2/2=1, 3/2=1, 4/2=2, 5/2=2, 6/2=3
	// Expected: 0, 1, 2, 3
	expression := "1d6/2"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 0 {
		t.Errorf("Expected MinValue 0, got %d", stats.MinValue)
	}
	if stats.MaxValue != 3 {
		t.Errorf("Expected MaxValue 3, got %d", stats.MaxValue)
	}
}

func TestCalculateDiceStatistics_Power(t *testing.T) {
	// 1d4 ^ 2
	// Outcomes: 1, 4, 9, 16
	expression := "1d4^2"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 16 {
		t.Errorf("Expected MaxValue 16, got %d", stats.MaxValue)
	}
	if _, exists := stats.Results[9]; !exists {
		t.Errorf("Expected outcome 9 to exist")
	}
}

func TestCalculateDiceStatistics_Decimal(t *testing.T) {
	// 2.5 + 2.5 = 2 + 2 = 4 (floor logic)
	expression := "2.5 + 2.5"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 4 {
		t.Errorf("Expected MinValue 4, got %d", stats.MinValue)
	}
}

// TestCalculateDiceStatistics_SingleD20 tests a single d20 roll
func TestCalculateDiceStatistics_SingleD20(t *testing.T) {
	stats, err := CalculateDiceStatistics("d20")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 20 {
		t.Errorf("Expected MaxValue 20, got %d", stats.MaxValue)
	}
	// All outcomes should be equally likely (1/20 probability)
	expectedProb := 1.0 / 20.0
	for value := 1; value <= 20; value++ {
		prob, exists := stats.Results[value]
		if !exists {
			t.Errorf("Expected outcome %d to exist", value)
		}
		if prob < expectedProb-0.01 || prob > expectedProb+0.01 {
			t.Errorf("d20 outcome %d has probability %f, expected %f", value, prob, expectedProb)
		}
	}
}

// TestCalculateDiceStatistics_AverageValue tests average calculation
func TestCalculateDiceStatistics_AverageValue(t *testing.T) {
	stats, err := CalculateDiceStatistics("1d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Average of 1d6 should be 3.5
	expectedAvg := 3.5
	if stats.Average < expectedAvg-0.1 || stats.Average > expectedAvg+0.1 {
		t.Errorf("Expected average ~3.5, got %f", stats.Average)
	}
}

// TestCalculateDiceStatistics_MostCommonValue tests most common outcome
func TestCalculateDiceStatistics_MostCommonValue(t *testing.T) {
	stats, err := CalculateDiceStatistics("2d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Most common roll for 2d6 should be 7
	if stats.MostCommon != 7 {
		t.Errorf("Expected MostCommon 7, got %d", stats.MostCommon)
	}
}

// TestCalculateDiceStatistics_SortedOutcomes tests sorted outcomes
func TestCalculateDiceStatistics_SortedOutcomes(t *testing.T) {
	stats, err := CalculateDiceStatistics("1d4")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	expected := []int{1, 2, 3, 4}
	if len(stats.SortedOutcomes) != len(expected) {
		t.Errorf("Expected %d outcomes, got %d", len(expected), len(stats.SortedOutcomes))
	}
	for i, outcome := range stats.SortedOutcomes {
		if outcome != expected[i] {
			t.Errorf("Expected outcome[%d] = %d, got %d", i, expected[i], outcome)
		}
	}
}

// TestCalculateDiceStatistics_Percentages tests percentage calculations
func TestCalculateDiceStatistics_Percentages(t *testing.T) {
	stats, err := CalculateDiceStatistics("d20")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Each outcome should be 5% (1/20)
	expectedPercentage := 5.0
	for value := 1; value <= 20; value++ {
		pct, exists := stats.Percentages[value]
		if !exists {
			t.Errorf("Expected percentage for outcome %d", value)
		}
		if pct < expectedPercentage-0.1 || pct > expectedPercentage+0.1 {
			t.Errorf("d20 outcome %d has percentage %f%%, expected %f%%", value, pct, expectedPercentage)
		}
	}
}

// TestCalculateDiceStatistics_TotalOutcomes tests total outcomes calculation
func TestCalculateDiceStatistics_TotalOutcomes(t *testing.T) {
	stats, err := CalculateDiceStatistics("2d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 2d6 should have 6*6 = 36 total outcomes
	if stats.TotalOutcomes.Int64() != 36 {
		t.Errorf("Expected 36 total outcomes for 2d6, got %d", stats.TotalOutcomes.Int64())
	}
}

// TestCalculateDiceStatistics_Subtraction tests subtraction in statistics
func TestCalculateDiceStatistics_Subtraction(t *testing.T) {
	expression := "1d6-1d4"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed to calculate %s: %v", expression, err)
	}
	// 1d6 (1-6) - 1d4 (1-4) gives range -3 to 5
	if stats.MinValue != -3 {
		t.Errorf("Expected MinValue -3, got %d", stats.MinValue)
	}
	if stats.MaxValue != 5 {
		t.Errorf("Expected MaxValue 5, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_MultiplyConstants tests multiplying constants
func TestCalculateDiceStatistics_MultiplyConstants(t *testing.T) {
	expression := "3 * 4"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 12 || stats.MaxValue != 12 {
		t.Errorf("Expected 12-12, got %d-%d", stats.MinValue, stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_DiceMultiplication tests dice multiplication
func TestCalculateDiceStatistics_DiceMultiplication(t *testing.T) {
	expression := "1d4 * 1d3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Min: 1*1 = 1, Max: 4*3 = 12
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 12 {
		t.Errorf("Expected MaxValue 12, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_HighestModifier tests highest die modifier
func TestCalculateDiceStatistics_HighestModifier(t *testing.T) {
	stats, err := CalculateDiceStatistics("4d6H")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 4d6H returns highest die (1-6)
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 6 {
		t.Errorf("Expected MaxValue 6, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_LowestModifier tests lowest die modifier
func TestCalculateDiceStatistics_LowestModifier(t *testing.T) {
	stats, err := CalculateDiceStatistics("4d6L")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 4d6L returns lowest die (1-6)
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 6 {
		t.Errorf("Expected MaxValue 6, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_EmptyExpression tests error handling
func TestCalculateDiceStatistics_EmptyExpression(t *testing.T) {
	_, err := CalculateDiceStatistics("")
	if err == nil {
		t.Errorf("Empty expression should return error")
	}
}

// TestCalculateDiceStatistics_InvalidDiceCount tests invalid dice count
func TestCalculateDiceStatistics_InvalidDiceCount(t *testing.T) {
	_, err := CalculateDiceStatistics("0d6")
	if err == nil {
		t.Errorf("0d6 should return error")
	}
}

// TestCalculateDiceStatistics_InvalidDiceSides tests invalid dice sides
func TestCalculateDiceStatistics_InvalidDiceSides(t *testing.T) {
	_, err := CalculateDiceStatistics("1d0")
	if err == nil {
		t.Errorf("1d0 should return error")
	}
}

// TestCalculateDiceStatistics_UnexpectedCharacter tests invalid character
func TestCalculateDiceStatistics_UnexpectedCharacter(t *testing.T) {
	_, err := CalculateDiceStatistics("1d6@")
	if err == nil {
		t.Errorf("Invalid character should return error")
	}
}

// TestCalculateDiceStatistics_LargeNumberOfDice tests statistics with many dice
func TestCalculateDiceStatistics_LargeNumberOfDice(t *testing.T) {
	stats, err := CalculateDiceStatistics("10d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 10d6 ranges from 10 to 60
	if stats.MinValue != 10 {
		t.Errorf("Expected MinValue 10, got %d", stats.MinValue)
	}
	if stats.MaxValue != 60 {
		t.Errorf("Expected MaxValue 60, got %d", stats.MaxValue)
	}
	// Average of 10d6 should be around 35 (10 * 3.5)
	expectedAvg := 35.0
	if stats.Average < expectedAvg-1 || stats.Average > expectedAvg+1 {
		t.Errorf("Expected average ~35, got %f", stats.Average)
	}
}

// TestCalculateDiceStatistics_DiceDivision tests dice division
func TestCalculateDiceStatistics_DiceDivision(t *testing.T) {
	expression := "1d10 / 2"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 1d10 (1-10) / 2: floor division gives 0-5
	if stats.MinValue != 0 {
		t.Errorf("Expected MinValue 0, got %d", stats.MinValue)
	}
	if stats.MaxValue != 5 {
		t.Errorf("Expected MaxValue 5, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_NegativeResults tests expressions with negative results
func TestCalculateDiceStatistics_NegativeResults(t *testing.T) {
	expression := "1d4 - 10"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 1d4 (1-4) - 10 = -9 to -6
	if stats.MinValue != -9 {
		t.Errorf("Expected MinValue -9, got %d", stats.MinValue)
	}
	if stats.MaxValue != -6 {
		t.Errorf("Expected MaxValue -6, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_Parentheses tests parentheses in expressions
func TestCalculateDiceStatistics_Parentheses(t *testing.T) {
	expression := "(1d4 + 2) * 3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// (1d4 + 2) * 3 = (3-6) * 3 = 9-18
	if stats.MinValue != 9 {
		t.Errorf("Expected MinValue 9, got %d", stats.MinValue)
	}
	if stats.MaxValue != 18 {
		t.Errorf("Expected MaxValue 18, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_Exponentiation tests exponentiation in expressions
func TestCalculateDiceStatistics_Exponentiation(t *testing.T) {
	expression := "1d4 ^ 2"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 1d4 ^ 2: outcomes 1, 4, 9, 16
	if stats.MinValue != 1 {
		t.Errorf("Expected MinValue 1, got %d", stats.MinValue)
	}
	if stats.MaxValue != 16 {
		t.Errorf("Expected MaxValue 16, got %d", stats.MaxValue)
	}
	// Check that 2 is NOT in results
	if _, exists := stats.Results[2]; exists {
		t.Errorf("Outcome 2 should not exist for 1d4^2")
	}
}

// TestCalculateDiceStatistics_MaxPercentage tests max percentage calculation
func TestCalculateDiceStatistics_MaxPercentage(t *testing.T) {
	stats, err := CalculateDiceStatistics("d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// For d6, all outcomes equally likely: 1/6 ≈ 16.67%
	expectedMaxPct := (1.0 / 6.0) * 100
	if stats.MaxPercentage < expectedMaxPct-1 || stats.MaxPercentage > expectedMaxPct+1 {
		t.Errorf("Expected MaxPercentage ~16.67%%, got %f%%", stats.MaxPercentage)
	}
}

// TestCalculateDiceStatistics_DistributionSum tests that distribution probabilities sum to 1
func TestCalculateDiceStatistics_DistributionSum(t *testing.T) {
	stats, err := CalculateDiceStatistics("2d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	sum := 0.0
	for _, prob := range stats.Results {
		sum += prob
	}
	// Probabilities should sum to ~1.0
	if sum < 0.999 || sum > 1.001 {
		t.Errorf("Probabilities should sum to 1.0, got %f", sum)
	}
}

// TestCalculateDiceStatistics_GetSortedOutcomes tests GetSortedOutcomes method
func TestCalculateDiceStatistics_GetSortedOutcomes(t *testing.T) {
	stats, err := CalculateDiceStatistics("1d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	outcomes := stats.GetSortedOutcomes()
	if len(outcomes) != 6 {
		t.Errorf("Expected 6 outcomes, got %d", len(outcomes))
	}
	// Check they're sorted
	for i := 1; i < len(outcomes); i++ {
		if outcomes[i] <= outcomes[i-1] {
			t.Errorf("Outcomes not properly sorted: %v", outcomes)
		}
	}
}

// TestCalculateDiceStatistics_GetMaxPercentage tests GetMaxPercentage method
func TestCalculateDiceStatistics_GetMaxPercentage(t *testing.T) {
	stats, err := CalculateDiceStatistics("d20")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	maxPct := stats.GetMaxPercentage()
	expectedMaxPct := 5.0 // 1/20 = 0.05 = 5%
	if maxPct < expectedMaxPct-0.1 || maxPct > expectedMaxPct+0.1 {
		t.Errorf("Expected MaxPercentage 5%%, got %f%%", maxPct)
	}
}

// TestCalculateDiceStatistics_ComplexExpression tests complex expression
func TestCalculateDiceStatistics_ComplexExpression(t *testing.T) {
	expression := "2d6 + 1d4 + 3"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// 2d6 (2-12) + 1d4 (1-4) + 3 = 6-19
	if stats.MinValue != 6 {
		t.Errorf("Expected MinValue 6, got %d", stats.MinValue)
	}
	if stats.MaxValue != 19 {
		t.Errorf("Expected MaxValue 19, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_ThreeDiceSum tests sum of three different dice
func TestCalculateDiceStatistics_ThreeDiceSum(t *testing.T) {
	expression := "1d4 + 1d6 + 1d8"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Min: 1+1+1 = 3, Max: 4+6+8 = 18
	if stats.MinValue != 3 {
		t.Errorf("Expected MinValue 3, got %d", stats.MinValue)
	}
	if stats.MaxValue != 18 {
		t.Errorf("Expected MaxValue 18, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_ImplicitMultiplication tests implicit multiplication
func TestCalculateDiceStatistics_ImplicitMultiplication(t *testing.T) {
	// 2d10 should be treated as 2 dice of 10 sides, not 2 * d10
	stats, err := CalculateDiceStatistics("2d10")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 2 {
		t.Errorf("Expected MinValue 2, got %d", stats.MinValue)
	}
	if stats.MaxValue != 20 {
		t.Errorf("Expected MaxValue 20, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_OutcomeAtBoundary tests outcomes at boundaries
func TestCalculateDiceStatistics_OutcomeAtBoundary(t *testing.T) {
	stats, err := CalculateDiceStatistics("1d20")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Check minimum boundary
	if _, exists := stats.Results[1]; !exists {
		t.Errorf("Expected outcome 1 for d20")
	}
	// Check maximum boundary
	if _, exists := stats.Results[20]; !exists {
		t.Errorf("Expected outcome 20 for d20")
	}
	// Check just outside boundaries
	if _, exists := stats.Results[0]; exists {
		t.Errorf("Outcome 0 should not exist for d20")
	}
	if _, exists := stats.Results[21]; exists {
		t.Errorf("Outcome 21 should not exist for d20")
	}
}

// TestCalculateDiceStatistics_SortedOutcomesAreUnique tests that sorted outcomes are unique
func TestCalculateDiceStatistics_SortedOutcomesAreUnique(t *testing.T) {
	stats, err := CalculateDiceStatistics("2d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	seen := make(map[int]bool)
	for _, outcome := range stats.SortedOutcomes {
		if seen[outcome] {
			t.Errorf("Outcome %d appears more than once in SortedOutcomes", outcome)
		}
		seen[outcome] = true
	}
}

// TestCalculateDiceStatistics_AverageInRange tests that average is within min/max
func TestCalculateDiceStatistics_AverageInRange(t *testing.T) {
	stats, err := CalculateDiceStatistics("1d20")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.Average < float64(stats.MinValue) || stats.Average > float64(stats.MaxValue) {
		t.Errorf("Average %f should be between min %d and max %d", stats.Average, stats.MinValue, stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_MostCommonInResults tests that MostCommon is in Results
func TestCalculateDiceStatistics_MostCommonInResults(t *testing.T) {
	stats, err := CalculateDiceStatistics("2d6")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if _, exists := stats.Results[stats.MostCommon]; !exists {
		t.Errorf("MostCommon value %d should exist in Results", stats.MostCommon)
	}
}

// TestCalculateDiceStatistics_WhitespaceHandling tests whitespace in expressions
func TestCalculateDiceStatistics_WhitespaceHandling(t *testing.T) {
	expression := "  1d6  +  2  "
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if stats.MinValue != 3 {
		t.Errorf("Expected MinValue 3, got %d", stats.MinValue)
	}
	if stats.MaxValue != 8 {
		t.Errorf("Expected MaxValue 8, got %d", stats.MaxValue)
	}
}

// TestCalculateDiceStatistics_ZeroOutcome tests expressions that can produce zero
func TestCalculateDiceStatistics_ZeroOutcome(t *testing.T) {
	expression := "1d4 - 1d4"
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	// Should include zero outcome
	if _, exists := stats.Results[0]; !exists {
		t.Errorf("Expected outcome 0 for 1d4 - 1d4")
	}
}
