package main

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gonum.org/v1/gonum/stat"
)

// DiceStatistics holds the theoretical statistics for a dice roll.
type DiceStatistics struct {
	MinValue          int
	MaxValue          int
	Results           Distribution // outcome -> probability
	TotalOutcomes     *big.Int     // exact number of equally likely underlying outcomes
	TotalOutcomesText string       // cached string form for UI
	Percentages       map[int]float64
	SortedOutcomes    []int
	MaxPercentage     float64
	Average           float64
	MostCommon        int
}

// Distribution represents the probability distribution of outcomes.
type Distribution map[int]float64

type distResult struct {
	dist          Distribution
	totalOutcomes *big.Int
}

type distributionEntry struct {
	value  int
	weight float64
}

const (
	combineParallelThreshold     = 4096
	convolutionParallelThreshold = 4096
	probabilityEpsilon           = 1e-12
)

// Regex patterns for parsing
var (
	diceTokenPattern   = regexp.MustCompile(`^([HL])?(\d*)d(\d+)([HL])?`)
	numberTokenPattern = regexp.MustCompile(`^(\d+(\.\d+)?)`)
)

// CalculateDiceStatistics calculates the theoretical distribution of possible outcomes for a dice expression.
func CalculateDiceStatistics(expression string) (*DiceStatistics, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, fmt.Errorf("empty expression")
	}

	parser := &statParser{expr: expression, pos: 0}
	result, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}

	parser.skipWhitespace()
	if parser.pos < len(parser.expr) {
		return nil, fmt.Errorf("unexpected character at position %d: '%c'", parser.pos, parser.expr[parser.pos])
	}

	if len(result.dist) == 0 {
		return nil, fmt.Errorf("no valid outcomes for expression")
	}

	stats := &DiceStatistics{
		Results:           normalizeDistribution(result.dist),
		TotalOutcomes:     cloneBigInt(result.totalOutcomes),
		TotalOutcomesText: cloneBigInt(result.totalOutcomes).String(),
		Percentages:       make(map[int]float64, len(result.dist)),
	}

	stats.populateDerivedFields()
	return stats, nil
}

// statParser implementation
type statParser struct {
	expr string
	pos  int
}

func (p *statParser) skipWhitespace() {
	for p.pos < len(p.expr) && (p.expr[p.pos] == ' ' || p.expr[p.pos] == '\t') {
		p.pos++
	}
}

// parseExpression handles addition and subtraction.
func (p *statParser) parseExpression() (distResult, error) {
	left, err := p.parseTerm()
	if err != nil {
		return distResult{}, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '+' {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependent(left, right, func(a, b int) int { return a + b })
		} else if p.expr[p.pos] == '-' {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependent(left, right, func(a, b int) int { return a - b })
		} else {
			break
		}
	}

	return left, nil
}

// parseTerm handles multiplication, division and implicit multiplication.
func (p *statParser) parseTerm() (distResult, error) {
	left, err := p.parsePower()
	if err != nil {
		return distResult{}, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		c := p.expr[p.pos]
		if c == '*' {
			p.pos++
			right, err := p.parsePower()
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependent(left, right, func(a, b int) int { return a * b })
		} else if c == '/' {
			p.pos++
			right, err := p.parsePower()
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependentFiltered(left, right, divideValues)
		} else if c == '(' || (c >= '0' && c <= '9') || c == 'd' || c == 'H' || c == 'L' {
			right, err := p.parsePower()
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependent(left, right, func(a, b int) int { return a * b })
		} else {
			break
		}
	}

	return left, nil
}

// parsePower handles exponentiation.
func (p *statParser) parsePower() (distResult, error) {
	left, err := p.parseFactor()
	if err != nil {
		return distResult{}, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}

		if p.expr[p.pos] == '^' {
			p.pos++
			right, err := p.parseFactor() // Left-associative to match calculator
			if err != nil {
				return distResult{}, err
			}
			left = combineIndependent(left, right, powerValues)
		} else {
			break
		}
	}

	return left, nil
}

// parseFactor handles parentheses, dice, and numbers.
func (p *statParser) parseFactor() (distResult, error) {
	p.skipWhitespace()
	if p.pos >= len(p.expr) {
		return distResult{}, fmt.Errorf("unexpected end of expression")
	}

	if p.expr[p.pos] == '(' {
		p.pos++
		dist, err := p.parseExpression()
		if err != nil {
			return distResult{}, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.expr) || p.expr[p.pos] != ')' {
			return distResult{}, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return dist, nil
	}

	remaining := p.expr[p.pos:]
	if loc := diceTokenPattern.FindStringIndex(remaining); loc != nil {
		token := remaining[loc[0]:loc[1]]
		p.pos += loc[1]
		return parseDiceToken(token)
	}

	if loc := numberTokenPattern.FindStringIndex(remaining); loc != nil {
		token := remaining[loc[0]:loc[1]]
		p.pos += loc[1]
		valFloat, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return distResult{}, fmt.Errorf("invalid number: %s", token)
		}
		return newDistResult(Distribution{int(valFloat): 1}), nil
	}

	return distResult{}, fmt.Errorf("unexpected character: %c", p.expr[p.pos])
}

func parseDiceToken(token string) (distResult, error) {
	matches := diceTokenPattern.FindStringSubmatch(token)
	if matches == nil {
		return distResult{}, fmt.Errorf("invalid dice term: %s", token)
	}

	prefixModifier := matches[1]
	countStr := matches[2]
	sidesStr := matches[3]
	suffixModifier := matches[4]

	count := 1
	if countStr != "" {
		c, err := strconv.Atoi(countStr)
		if err != nil {
			return distResult{}, err
		}
		count = c
	}

	sides, err := strconv.Atoi(sidesStr)
	if err != nil {
		return distResult{}, err
	}
	if count <= 0 || sides <= 0 {
		return distResult{}, fmt.Errorf("dice terms must use positive counts and sides")
	}

	modifier := ""
	if suffixModifier != "" {
		modifier = suffixModifier
	} else if prefixModifier != "" {
		modifier = prefixModifier
	}

	return getDiceOutcomes(count, sides, modifier), nil
}

func newDistResult(dist Distribution) distResult {
	return distResult{
		dist:          normalizeDistribution(dist),
		totalOutcomes: big.NewInt(1),
	}
}

func cloneBigInt(v *big.Int) *big.Int {
	if v == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(v)
}

func multiplyOutcomeCounts(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(cloneBigInt(a), cloneBigInt(b))
}

func normalizeDistribution(dist Distribution) Distribution {
	total := 0.0
	for _, weight := range dist {
		total += weight
	}
	if total == 0 {
		return dist
	}

	if math.Abs(total-1) <= probabilityEpsilon {
		pruned := make(Distribution, len(dist))
		for value, weight := range dist {
			if weight != 0 {
				pruned[value] = weight
			}
		}
		return pruned
	}

	normalized := make(Distribution, len(dist))
	for value, weight := range dist {
		probability := weight / total
		if probability != 0 {
			normalized[value] = probability
		}
	}
	return normalized
}

func combineIndependent(a, b distResult, op func(int, int) int) distResult {
	return distResult{
		dist:          combineDistributions(a.dist, b.dist, op),
		totalOutcomes: multiplyOutcomeCounts(a.totalOutcomes, b.totalOutcomes),
	}
}

func combineIndependentFiltered(a, b distResult, op func(int, int) (int, bool)) distResult {
	return distResult{
		dist:          combineDistributionsFiltered(a.dist, b.dist, op),
		totalOutcomes: multiplyOutcomeCounts(a.totalOutcomes, b.totalOutcomes),
	}
}

func combineDistributions(a, b Distribution, op func(int, int) int) Distribution {
	if len(a) == 0 || len(b) == 0 {
		return Distribution{}
	}

	entriesA := make([]distributionEntry, 0, len(a))
	for value, weight := range a {
		entriesA = append(entriesA, distributionEntry{value: value, weight: weight})
	}

	entriesB := make([]distributionEntry, 0, len(b))
	for value, weight := range b {
		entriesB = append(entriesB, distributionEntry{value: value, weight: weight})
	}

	workSize := len(entriesA) * len(entriesB)
	if workSize < combineParallelThreshold || len(entriesA) < 2 {
		return combineDistributionsSerial(entriesA, entriesB, op)
	}

	workerCount := runtime.GOMAXPROCS(0)
	if workerCount > len(entriesA) {
		workerCount = len(entriesA)
	}
	if workerCount < 2 {
		return combineDistributionsSerial(entriesA, entriesB, op)
	}

	chunkSize := (len(entriesA) + workerCount - 1) / workerCount
	partials := make([]Distribution, workerCount)
	var wg sync.WaitGroup

	for worker := 0; worker < workerCount; worker++ {
		start := worker * chunkSize
		if start >= len(entriesA) {
			break
		}
		end := start + chunkSize
		if end > len(entriesA) {
			end = len(entriesA)
		}

		wg.Add(1)
		go func(workerIndex, from, to int) {
			defer wg.Done()
			local := make(Distribution)
			for _, left := range entriesA[from:to] {
				for _, right := range entriesB {
					local[op(left.value, right.value)] += left.weight * right.weight
				}
			}
			partials[workerIndex] = local
		}(worker, start, end)
	}

	wg.Wait()

	result := make(Distribution)
	for _, partial := range partials {
		for value, weight := range partial {
			result[value] += weight
		}
	}
	return normalizeDistribution(result)
}

func combineDistributionsFiltered(a, b Distribution, op func(int, int) (int, bool)) Distribution {
	if len(a) == 0 || len(b) == 0 {
		return Distribution{}
	}

	result := make(Distribution)
	for leftValue, leftWeight := range a {
		for rightValue, rightWeight := range b {
			value, ok := op(leftValue, rightValue)
			if !ok {
				continue
			}
			result[value] += leftWeight * rightWeight
		}
	}
	return normalizeDistribution(result)
}

func combineDistributionsSerial(entriesA, entriesB []distributionEntry, op func(int, int) int) Distribution {
	result := make(Distribution)
	for _, left := range entriesA {
		for _, right := range entriesB {
			result[op(left.value, right.value)] += left.weight * right.weight
		}
	}
	return normalizeDistribution(result)
}

func divideValues(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func powerValues(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

// getDiceOutcomes returns a probability distribution for a dice roll and the exact total outcome count.
func getDiceOutcomes(count int, sides int, modifier string) distResult {
	totalOutcomes := new(big.Int).Exp(big.NewInt(int64(sides)), big.NewInt(int64(count)), nil)

	var outcomes Distribution
	switch modifier {
	case "H":
		outcomes = highestDieDistribution(count, sides)
	case "L":
		outcomes = lowestDieDistribution(count, sides)
	default:
		outcomes = sumDiceDistribution(count, sides)
	}

	return distResult{
		dist:          outcomes,
		totalOutcomes: totalOutcomes,
	}
}

func sumDiceDistribution(count int, sides int) Distribution {
	current := []float64{1}
	singleDieProbability := 1 / float64(sides)

	for die := 0; die < count; die++ {
		next := make([]float64, len(current)+sides)
		if len(current)*sides >= convolutionParallelThreshold {
			convolveStepParallel(current, next, sides, singleDieProbability)
		} else {
			convolveStepSerial(current, next, sides, singleDieProbability)
		}
		current = next
	}

	outcomes := make(Distribution, count*(sides-1)+1)
	for sum, probability := range current {
		if probability != 0 {
			outcomes[sum] = probability
		}
	}
	return normalizeDistribution(outcomes)
}

func convolveStepSerial(current, next []float64, sides int, singleDieProbability float64) {
	for sum, probability := range current {
		if probability == 0 {
			continue
		}
		for face := 1; face <= sides; face++ {
			next[sum+face] += probability * singleDieProbability
		}
	}
}

func convolveStepParallel(current, next []float64, sides int, singleDieProbability float64) {
	workerCount := runtime.GOMAXPROCS(0)
	if workerCount > len(current) {
		workerCount = len(current)
	}
	if workerCount < 2 {
		convolveStepSerial(current, next, sides, singleDieProbability)
		return
	}

	chunkSize := (len(current) + workerCount - 1) / workerCount
	partials := make([][]float64, workerCount)
	var wg sync.WaitGroup

	for worker := 0; worker < workerCount; worker++ {
		start := worker * chunkSize
		if start >= len(current) {
			break
		}
		end := start + chunkSize
		if end > len(current) {
			end = len(current)
		}

		wg.Add(1)
		go func(workerIndex, from, to int) {
			defer wg.Done()
			local := make([]float64, len(next))
			for sum := from; sum < to; sum++ {
				probability := current[sum]
				if probability == 0 {
					continue
				}
				for face := 1; face <= sides; face++ {
					local[sum+face] += probability * singleDieProbability
				}
			}
			partials[workerIndex] = local
		}(worker, start, end)
	}

	wg.Wait()

	for _, partial := range partials {
		if partial == nil {
			continue
		}
		for i, probability := range partial {
			next[i] += probability
		}
	}
}

func highestDieDistribution(count int, sides int) Distribution {
	outcomes := make(Distribution, sides)
	denominator := float64(sides)
	for value := 1; value <= sides; value++ {
		current := math.Pow(float64(value)/denominator, float64(count))
		previous := math.Pow(float64(value-1)/denominator, float64(count))
		probability := current - previous
		if probability != 0 {
			outcomes[value] = probability
		}
	}
	return normalizeDistribution(outcomes)
}

func lowestDieDistribution(count int, sides int) Distribution {
	outcomes := make(Distribution, sides)
	denominator := float64(sides)
	for value := 1; value <= sides; value++ {
		current := math.Pow(float64(sides-value+1)/denominator, float64(count))
		next := math.Pow(float64(sides-value)/denominator, float64(count))
		probability := current - next
		if probability != 0 {
			outcomes[value] = probability
		}
	}
	return normalizeDistribution(outcomes)
}

func (s *DiceStatistics) populateDerivedFields() {
	if len(s.Results) == 0 {
		s.TotalOutcomesText = cloneBigInt(s.TotalOutcomes).String()
		return
	}

	s.SortedOutcomes = make([]int, 0, len(s.Results))
	s.Percentages = make(map[int]float64, len(s.Results))

	// Prepare data for weighted mean calculation using gonum
	outcomes := make([]float64, 0, len(s.Results))
	weights := make([]float64, 0, len(s.Results))

	first := true
	maxProbability := 0.0
	for value, probability := range s.Results {
		s.SortedOutcomes = append(s.SortedOutcomes, value)
		s.Percentages[value] = probability * 100

		if first {
			s.MinValue = value
			s.MaxValue = value
			s.MostCommon = value
			first = false
		} else {
			if value < s.MinValue {
				s.MinValue = value
			}
			if value > s.MaxValue {
				s.MaxValue = value
			}
		}

		// Collect data for weighted mean calculation
		outcomes = append(outcomes, float64(value))
		weights = append(weights, probability)

		// Mode: find outcome with highest probability
		if probability > maxProbability || (math.Abs(probability-maxProbability) <= probabilityEpsilon && value < s.MostCommon) {
			maxProbability = probability
			s.MostCommon = value
		}
	}

	// Calculate weighted mean using gonum/stat
	s.Average = stat.Mean(outcomes, weights)

	sort.Ints(s.SortedOutcomes)
	s.MaxPercentage = maxProbability * 100
}

// GetSortedOutcomes returns cached sorted unique outcomes.
func (s *DiceStatistics) GetSortedOutcomes() []int {
	return s.SortedOutcomes
}

// GetMaxPercentage returns the cached maximum percentage value.
func (s *DiceStatistics) GetMaxPercentage() float64 {
	return s.MaxPercentage
}
