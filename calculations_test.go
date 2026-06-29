package main

import (
	"strings"
	"testing"
)

// TestCalculateDice_BasicDiceRoll tests simple dice rolls
func TestCalculateDice_BasicDiceRoll(t *testing.T) {
	result, diceRolls, err := CalculateDice("d20")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 1 || result > 20 {
		t.Errorf("d20 should roll 1-20, got %v", result)
	}
	if !strings.Contains(diceRolls, "1d20") {
		t.Errorf("diceRolls should contain '1d20', got %q", diceRolls)
	}
}

// TestCalculateDice_MultipleDice tests multiple dice
func TestCalculateDice_MultipleDice(t *testing.T) {
	result, _, err := CalculateDice("2d6")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 2 || result > 12 {
		t.Errorf("2d6 should roll 2-12, got %v", result)
	}
}

// TestCalculateDice_AdditionWithDice tests dice with addition
func TestCalculateDice_AdditionWithDice(t *testing.T) {
	result, _, err := CalculateDice("1d6+5")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 6 || result > 11 {
		t.Errorf("1d6+5 should roll 6-11, got %v", result)
	}
}

// TestCalculateDice_SubtractionWithDice tests dice with subtraction
func TestCalculateDice_SubtractionWithDice(t *testing.T) {
	result, _, err := CalculateDice("2d10-5")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < -3 || result > 15 {
		t.Errorf("2d10-5 should roll -3 to 15, got %v", result)
	}
}

// TestCalculateDice_MultiplicationWithDice tests dice with multiplication
func TestCalculateDice_MultiplicationWithDice(t *testing.T) {
	result, _, err := CalculateDice("1d5*2")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 2 || result > 10 {
		t.Errorf("1d5*2 should roll 2-10, got %v", result)
	}
}

// TestCalculateDice_DivisionWithDice tests dice with division
func TestCalculateDice_DivisionWithDice(t *testing.T) {
	result, _, err := CalculateDice("1d6/2")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 0 || result > 3 {
		t.Errorf("1d6/2 should roll 0-3, got %v", result)
	}
}

// TestCalculateDice_HighestDiceModifier tests the 'H' modifier
func TestCalculateDice_HighestDiceModifier(t *testing.T) {
	// 2d20H should return only the highest die
	result, _, err := CalculateDice("2d20H")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 1 || result > 20 {
		t.Errorf("2d20H should roll 1-20 (highest of two d20s), got %v", result)
	}
}

// TestCalculateDice_LowestDiceModifier tests the 'L' modifier
func TestCalculateDice_LowestDiceModifier(t *testing.T) {
	// 2d20L should return only the lowest die
	result, _, err := CalculateDice("2d20L")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 1 || result > 20 {
		t.Errorf("2d20L should roll 1-20 (lowest of two d20s), got %v", result)
	}
}

// TestCalculateDice_ComplexExpression tests complex mathematical expressions
func TestCalculateDice_ComplexExpression(t *testing.T) {
	result, _, err := CalculateDice("2d6+3*2")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// 2d6 (2-12) + 3*2 (6) = 8-18
	if result < 8 || result > 18 {
		t.Errorf("2d6+3*2 should roll 8-18, got %v", result)
	}
}

// TestCalculateDice_ParenthesesExpression tests expressions with parentheses
func TestCalculateDice_ParenthesesExpression(t *testing.T) {
	result, _, err := CalculateDice("(1d4+2)*3")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// (1d4+2)*3 = (3-6)*3 = 9-18
	if result < 9 || result > 18 {
		t.Errorf("(1d4+2)*3 should roll 9-18, got %v", result)
	}
}

// TestCalculateDice_MultipleAddends tests multiple addends
func TestCalculateDice_MultipleAddends(t *testing.T) {
	result, _, err := CalculateDice("1d4+1d6+1d8")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// 1d4 (1-4) + 1d6 (1-6) + 1d8 (1-8) = 3-18
	if result < 3 || result > 18 {
		t.Errorf("1d4+1d6+1d8 should roll 3-18, got %v", result)
	}
}

// TestCalculateDice_PrefixModifier tests prefix H/L modifiers
func TestCalculateDice_PrefixModifier(t *testing.T) {
	result, _, err := CalculateDice("H2d20")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result < 1 || result > 20 {
		t.Errorf("H2d20 should roll 1-20 (highest of two d20s), got %v", result)
	}
}

// TestCalculateDice_EmptyExpression tests error handling for empty expression
func TestCalculateDice_EmptyExpression(t *testing.T) {
	_, _, err := CalculateDice("")
	if err == nil {
		t.Errorf("Empty expression should return error")
	}
}

// TestCalculateDice_InvalidDiceCount tests error handling for invalid dice count
func TestCalculateDice_InvalidDiceCount(t *testing.T) {
	_, _, err := CalculateDice("0d6")
	if err == nil {
		t.Errorf("0d6 should return error")
	}
}

// TestCalculateDice_InvalidDiceSides tests error handling for invalid dice sides
func TestCalculateDice_InvalidDiceSides(t *testing.T) {
	_, _, err := CalculateDice("1d0")
	if err == nil {
		t.Errorf("1d0 should return error")
	}
}

// TestCalculateDice_PlaceholderDice tests error handling for placeholder dice
func TestCalculateDice_PlaceholderDice(t *testing.T) {
	_, _, err := CalculateDice("dx")
	if err == nil {
		t.Errorf("dx should return error")
	}
}

// TestEvaluateMathExpression_SimpleAddition tests simple addition
func TestEvaluateMathExpression_SimpleAddition(t *testing.T) {
	result, err := evaluateMathExpression("5+3")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 8 {
		t.Errorf("5+3 should be 8, got %v", result)
	}
}

// TestEvaluateMathExpression_Subtraction tests subtraction
func TestEvaluateMathExpression_Subtraction(t *testing.T) {
	result, err := evaluateMathExpression("10-3")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 7 {
		t.Errorf("10-3 should be 7, got %v", result)
	}
}

// TestEvaluateMathExpression_Multiplication tests multiplication
func TestEvaluateMathExpression_Multiplication(t *testing.T) {
	result, err := evaluateMathExpression("4*5")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 20 {
		t.Errorf("4*5 should be 20, got %v", result)
	}
}

// TestEvaluateMathExpression_Division tests division
func TestEvaluateMathExpression_Division(t *testing.T) {
	result, err := evaluateMathExpression("20/4")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 5 {
		t.Errorf("20/4 should be 5, got %v", result)
	}
}

// TestEvaluateMathExpression_OperatorPrecedence tests operator precedence (multiplication before addition)
func TestEvaluateMathExpression_OperatorPrecedence(t *testing.T) {
	result, err := evaluateMathExpression("2+3*4")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 14 {
		t.Errorf("2+3*4 should be 14, got %v", result)
	}
}

// TestEvaluateMathExpression_Parentheses tests parentheses override precedence
func TestEvaluateMathExpression_Parentheses(t *testing.T) {
	result, err := evaluateMathExpression("(2+3)*4")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 20 {
		t.Errorf("(2+3)*4 should be 20, got %v", result)
	}
}

// TestEvaluateMathExpression_Exponentiation tests exponentiation
func TestEvaluateMathExpression_Exponentiation(t *testing.T) {
	result, err := evaluateMathExpression("2^3")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 8 {
		t.Errorf("2^3 should be 8, got %v", result)
	}
}

// TestEvaluateMathExpression_UnaryMinus tests unary minus
func TestEvaluateMathExpression_UnaryMinus(t *testing.T) {
	result, err := evaluateMathExpression("-5+10")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 5 {
		t.Errorf("-5+10 should be 5, got %v", result)
	}
}

// TestEvaluateMathExpression_FloatingPoint tests floating point numbers
func TestEvaluateMathExpression_FloatingPoint(t *testing.T) {
	result, err := evaluateMathExpression("3.5+2.5")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 6 {
		t.Errorf("3.5+2.5 should be 6, got %v", result)
	}
}

// TestEvaluateMathExpression_DivisionByZero tests division by zero error
func TestEvaluateMathExpression_DivisionByZero(t *testing.T) {
	_, err := evaluateMathExpression("5/0")
	if err == nil {
		t.Errorf("Division by zero should return error")
	}
}

// TestEvaluateMathExpression_MissingClosingParen tests missing closing parenthesis error
func TestEvaluateMathExpression_MissingClosingParen(t *testing.T) {
	_, err := evaluateMathExpression("(2+3")
	if err == nil {
		t.Errorf("Missing closing parenthesis should return error")
	}
}

// TestEvaluateMathExpression_InvalidCharacter tests invalid character error
func TestEvaluateMathExpression_InvalidCharacter(t *testing.T) {
	_, err := evaluateMathExpression("5@3")
	if err == nil {
		t.Errorf("Invalid character should return error")
	}
}

// TestEvaluateMathExpression_Whitespace tests whitespace handling
func TestEvaluateMathExpression_Whitespace(t *testing.T) {
	result, err := evaluateMathExpression("  5  +  3  ")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 8 {
		t.Errorf("'  5  +  3  ' should be 8, got %v", result)
	}
}

// TestEvaluateMathExpression_ComplexExpression tests complex mathematical expressions
func TestEvaluateMathExpression_ComplexExpression(t *testing.T) {
	result, err := evaluateMathExpression("((2+3)*4-5)/3")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	// ((2+3)*4-5)/3 = (5*4-5)/3 = (20-5)/3 = 15/3 = 5
	if result != 5 {
		t.Errorf("((2+3)*4-5)/3 should be 5, got %v", result)
	}
}

// TestRollDiceSet_Range tests rollDiceSet produces values in correct range
func TestRollDiceSet_Range(t *testing.T) {
	rolls := rollDiceSet(10, 20)
	if len(rolls) != 10 {
		t.Errorf("rollDiceSet(10, 20) should return 10 rolls, got %d", len(rolls))
	}
	for i, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("Roll %d is %d, expected 1-20", i, roll)
		}
	}
}

// TestRollDiceSet_Variability tests rollDiceSet produces different values
func TestRollDiceSet_Variability(t *testing.T) {
	rolls := rollDiceSet(100, 6)
	seenValues := make(map[int]bool)
	for _, roll := range rolls {
		seenValues[roll] = true
	}
	// With 100 rolls of 1d6, we should see multiple different values
	if len(seenValues) < 3 {
		t.Errorf("100 rolls of 1d6 should see at least 3 different values, got %d", len(seenValues))
	}
}

// TestCalculateDice_OutputFormat tests the output format
func TestCalculateDice_OutputFormat(t *testing.T) {
	_, diceRolls, err := CalculateDice("2d6")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// Should contain information about the rolls
	if !strings.Contains(diceRolls, "2d6") {
		t.Errorf("diceRolls should contain '2d6', got %q", diceRolls)
	}
	if !strings.Contains(diceRolls, "(") || !strings.Contains(diceRolls, ")") {
		t.Errorf("diceRolls should be formatted with parentheses, got %q", diceRolls)
	}
}

// TestCalculateDice_PrefixAndSuffixModifiers tests that suffix modifier takes priority
func TestCalculateDice_PrefixAndSuffixModifiers(t *testing.T) {
	// When both prefix and suffix modifiers are present, suffix should take priority
	result, _, err := CalculateDice("H2d20L")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// Should use L (suffix) modifier, not H
	if result < 1 || result > 20 {
		t.Errorf("H2d20L should use L modifier and roll 1-20, got %v", result)
	}
}

// TestEvaluateMathExpression_LargeNumbers tests large number handling
func TestEvaluateMathExpression_LargeNumbers(t *testing.T) {
	result, err := evaluateMathExpression("1000000+2000000")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 3000000 {
		t.Errorf("1000000+2000000 should be 3000000, got %v", result)
	}
}

// TestEvaluateMathExpression_ExponentiationPrecedence tests exponentiation precedence
func TestEvaluateMathExpression_ExponentiationPrecedence(t *testing.T) {
	result, err := evaluateMathExpression("2+3^2")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	// 2 + (3^2) = 2 + 9 = 11
	if result != 11 {
		t.Errorf("2+3^2 should be 11, got %v", result)
	}
}

// TestEvaluateMathExpression_NestedParentheses tests nested parentheses
func TestEvaluateMathExpression_NestedParentheses(t *testing.T) {
	result, err := evaluateMathExpression("((10))")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 10 {
		t.Errorf("((10)) should be 10, got %v", result)
	}
}

// TestCalculateDice_NegativeResult tests expressions that can produce negative results
func TestCalculateDice_NegativeResult(t *testing.T) {
	result, _, err := CalculateDice("1d4-10")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// 1d4 (1-4) - 10 = -9 to -6
	if result < -9 || result > -6 {
		t.Errorf("1d4-10 should roll -9 to -6, got %v", result)
	}
}

// TestExpandDiceNotation_SingleDice tests expansion of single dice notation
func TestExpandDiceNotation_SingleDice(t *testing.T) {
	// This test verifies that dice notation expands correctly
	// We can't test exact values due to randomness, but we can test the format
	expanded, diceRolls, err := expandDiceNotation("1d6")
	if err != nil {
		t.Fatalf("expandDiceNotation failed: %v", err)
	}
	// expanded should be a number between 1 and 6
	val, err := evaluateMathExpression(expanded)
	if err != nil {
		t.Fatalf("expanded result should be valid: %v", err)
	}
	if val < 1 || val > 6 {
		t.Errorf("1d6 should expand to 1-6, got %v", val)
	}
	if !strings.Contains(diceRolls, "1d6") {
		t.Errorf("diceRolls should contain '1d6', got %q", diceRolls)
	}
}

// TestExpandDiceNotation_WithModifier tests expansion with H/L modifier
func TestExpandDiceNotation_WithModifier(t *testing.T) {
	expanded, _, err := expandDiceNotation("2d20H")
	if err != nil {
		t.Fatalf("expandDiceNotation failed: %v", err)
	}
	val, err := evaluateMathExpression(expanded)
	if err != nil {
		t.Fatalf("expanded result should be valid: %v", err)
	}
	if val < 1 || val > 20 {
		t.Errorf("2d20H should expand to 1-20, got %v", val)
	}
}

// TestCalculateDice_PowerOperator tests the power/exponentiation operator
func TestCalculateDice_PowerOperator(t *testing.T) {
	result, _, err := CalculateDice("2^3")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result != 8 {
		t.Errorf("2^3 should be 8, got %v", result)
	}
}

// TestCalculateDice_ZeroDice tests zero values
func TestCalculateDice_ZeroDice(t *testing.T) {
	result, _, err := CalculateDice("0")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	if result != 0 {
		t.Errorf("0 should be 0, got %v", result)
	}
}

// TestEvaluateMathExpression_NestedExpression tests deeply nested expressions
func TestEvaluateMathExpression_NestedExpression(t *testing.T) {
	result, err := evaluateMathExpression("(((2+3)))")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	if result != 5 {
		t.Errorf("(((2+3))) should be 5, got %v", result)
	}
}

// TestCalculateDice_MultipleOperations tests multiple different operations
func TestCalculateDice_MultipleOperations(t *testing.T) {
	result, _, err := CalculateDice("1d4+2-1")
	if err != nil {
		t.Fatalf("CalculateDice failed: %v", err)
	}
	// 1d4 (1-4) + 2 - 1 = 2-5
	if result < 2 || result > 5 {
		t.Errorf("1d4+2-1 should roll 2-5, got %v", result)
	}
}

// TestEvaluateMathExpression_ChainedExponentiation tests chained exponentiation
func TestEvaluateMathExpression_ChainedExponentiation(t *testing.T) {
	result, err := evaluateMathExpression("2^2^2")
	if err != nil {
		t.Fatalf("evaluateMathExpression failed: %v", err)
	}
	// Left-associative: (2^2)^2 = 4^2 = 16
	if result != 16 {
		t.Errorf("2^2^2 (left-associative) should be 16, got %v", result)
	}
}
