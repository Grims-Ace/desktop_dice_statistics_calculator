package main

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// barGraphCanvas is a custom widget that renders a bar graph
type barGraphCanvas struct {
	widget.BaseWidget
	stats *DiceStatistics
}

func newBarGraphCanvas(stats *DiceStatistics) *barGraphCanvas {
	graph := &barGraphCanvas{
		stats: stats,
	}
	graph.ExtendBaseWidget(graph)
	return graph
}

func (b *barGraphCanvas) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	return &barGraphCanvasRenderer{
		graph: b,
	}
}

type barGraphCanvasRenderer struct {
	graph   *barGraphCanvas
	objects []fyne.CanvasObject
}

func (r *barGraphCanvasRenderer) Layout(size fyne.Size) {
	r.Refresh()
}

func (r *barGraphCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(320, 240)
}

func (r *barGraphCanvasRenderer) Refresh() {
	r.objects = []fyne.CanvasObject{}

	if r.graph.stats == nil || len(r.graph.stats.Results) == 0 {
		return
	}

	stats := r.graph.stats
	outcomes := stats.GetSortedOutcomes()
	maxPercentage := stats.GetMaxPercentage()

	axisMaxPercent, tickStep := calculateYAxisScale(maxPercentage)

	size := r.graph.Size()
	if size.Width == 0 || size.Height == 0 {
		size = fyne.NewSize(900, 550)
	}

	// Padding and typography scale with the available area so the graph fits the window.
	topPadding := clamp(size.Height*0.14, 48, 75)
	bottomPadding := clamp(size.Height*0.15, 56, 80)
	leftPadding := clamp(size.Width*0.12, 60, 100)
	rightPadding := clamp(size.Width*0.03, 16, 24)

	graphWidth := size.Width - leftPadding - rightPadding
	graphHeight := size.Height - topPadding - bottomPadding
	if graphWidth <= 0 || graphHeight <= 0 {
		return
	}

	titleSize := clamp(size.Height*0.03, 12, 16)
	bodyTextSize := clamp(size.Height*0.02, 9, 11)
	labelTextSize := clamp(size.Height*0.022, 10, 12)
	tickTextSize := clamp(size.Height*0.018, 8, 10)

	// Background
	background := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 255})
	background.Move(fyne.NewPos(0, 0))
	background.Resize(size)
	r.objects = append(r.objects, background)

	// Y-axis
	yAxisLine := canvas.NewLine(color.White)
	yAxisLine.StrokeWidth = 2
	yAxisLine.Move(fyne.NewPos(leftPadding, topPadding))
	yAxisLine.Resize(fyne.NewSize(0, graphHeight))
	r.objects = append(r.objects, yAxisLine)

	// X-axis
	xAxisLine := canvas.NewLine(color.White)
	xAxisLine.StrokeWidth = 2
	xAxisLine.Move(fyne.NewPos(leftPadding, topPadding+graphHeight))
	xAxisLine.Resize(fyne.NewSize(graphWidth, 0))
	r.objects = append(r.objects, xAxisLine)

	// Title
	title := canvas.NewText("Probability Distribution", color.White)
	title.TextSize = titleSize
	title.Move(fyne.NewPos(leftPadding, 5))
	r.objects = append(r.objects, title)

	// Statistics info line 1
	statsLine1 := canvas.NewText(fmt.Sprintf("Range: %d to %d  |  Total Outcomes: %d", stats.MinValue, stats.MaxValue, stats.Total), color.White)
	statsLine1.TextSize = bodyTextSize
	statsLine1.Move(fyne.NewPos(leftPadding, 5+title.MinSize().Height))
	r.objects = append(r.objects, statsLine1)

	// Statistics info line 2
	statsLine2 := canvas.NewText(fmt.Sprintf("Average: %.2f  |  Most Common: %d", stats.Average, stats.MostCommon), color.White)
	statsLine2.TextSize = bodyTextSize
	statsLine2.Move(fyne.NewPos(leftPadding, 5+title.MinSize().Height+statsLine1.MinSize().Height))
	r.objects = append(r.objects, statsLine2)

	// Y-axis label
	yLabel := canvas.NewText("Probability (%)", color.White)
	yLabel.TextSize = labelTextSize
	yLabel.Move(fyne.NewPos(clamp(leftPadding*0.15, 8, 15), topPadding+graphHeight/2-yLabel.MinSize().Height/2))
	r.objects = append(r.objects, yLabel)

	// X-axis label
	xLabel := canvas.NewText("Result Value", color.White)
	xLabel.TextSize = labelTextSize
	xLabel.Move(fyne.NewPos(leftPadding+graphWidth/2-xLabel.MinSize().Width/2, topPadding+graphHeight+clamp(bottomPadding*0.35, 18, 28)))
	r.objects = append(r.objects, xLabel)

	// Y-axis tick marks and labels
	numYTicks := int(math.Round(axisMaxPercent / tickStep))
	for i := 0; i <= numYTicks; i++ {
		percent := float64(i) * tickStep

		yPos := topPadding + graphHeight - (float32(percent/axisMaxPercent) * graphHeight)

		// Tick mark
		tick := canvas.NewLine(color.White)
		tick.StrokeWidth = 1
		tick.Move(fyne.NewPos(leftPadding-5, yPos))
		tick.Resize(fyne.NewSize(5, 0))
		r.objects = append(r.objects, tick)

		// Label
		label := canvas.NewText(formatPercentLabel(percent), color.White)
		label.TextSize = tickTextSize
		label.Move(fyne.NewPos(leftPadding-50, yPos-7))
		r.objects = append(r.objects, label)
	}

	// Draw bars
	numBars := len(outcomes)
	barSpacing := float32(2)
	totalSpacing := float32(numBars+1) * barSpacing
	barWidth := (graphWidth - totalSpacing) / float32(numBars)
	if barWidth < 1 {
		barWidth = 1
	}

	// Calculate label step to prevent overlapping
	labelStep := calculateLabelStep(graphWidth, numBars)

	for i, value := range outcomes {
		percentage := stats.Percentages[value]

		// Bar height proportional to percentage
		barHeight := (float32(percentage) / float32(axisMaxPercent)) * graphHeight

		// X position
		xPos := leftPadding + barSpacing + float32(i)*(barWidth+barSpacing)

		// Draw bar
		bar := canvas.NewRectangle(color.NRGBA{R: 100, G: 180, B: 255, A: 255})
		bar.Move(fyne.NewPos(xPos, topPadding+graphHeight-barHeight))
		bar.Resize(fyne.NewSize(barWidth, barHeight))
		r.objects = append(r.objects, bar)

		// X-axis label
		// Always show first and last label
		isFirst := i == 0
		isLast := i == numBars-1

		// Determine if we should show this intermediate label
		// We show it if it matches the step, BUT we also need to make sure it doesn't clash with the last label
		// So if we are very close to the end, don't show it (unless it IS the end)
		showIntermediate := i%labelStep == 0 && i < numBars-labelStep

		if isFirst || isLast || showIntermediate {
			label := canvas.NewText(fmt.Sprintf("%d", value), color.White)
			label.TextSize = tickTextSize

			// Center label under bar
			label.Alignment = fyne.TextAlignCenter
			label.Move(fyne.NewPos(xPos+barWidth/2-label.MinSize().Width/2, topPadding+graphHeight+10))

			r.objects = append(r.objects, label)
		}
	}
}

func clamp(value, minValue, maxValue float32) float32 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func calculateLabelStep(graphWidth float32, numBars int) int {
	labelWidthEstimate := float32(35) // Estimate width of a label
	maxLabels := int(graphWidth / labelWidthEstimate)
	if maxLabels < 1 {
		maxLabels = 1
	}
	labelStep := int(math.Ceil(float64(numBars) / float64(maxLabels)))
	if labelStep < 1 {
		labelStep = 1
	}
	return labelStep
}

func calculateYAxisScale(maxPercentage float64) (axisMax float64, tickStep float64) {
	if maxPercentage <= 0 {
		return 1, 0.2
	}

	targetTicks := 5.0
	rawStep := maxPercentage / targetTicks
	magnitude := math.Pow(10, math.Floor(math.Log10(rawStep)))
	normalized := rawStep / magnitude

	switch {
	case normalized <= 1:
		tickStep = 1 * magnitude
	case normalized <= 2:
		tickStep = 2 * magnitude
	case normalized <= 2.5:
		tickStep = 2.5 * magnitude
	case normalized <= 5:
		tickStep = 5 * magnitude
	default:
		tickStep = 10 * magnitude
	}

	axisMax = math.Ceil(maxPercentage/tickStep) * tickStep
	return axisMax, tickStep
}

func formatPercentLabel(percent float64) string {
	if math.Abs(percent-math.Round(percent)) < 0.0001 {
		return fmt.Sprintf("%.0f%%", percent)
	}
	return fmt.Sprintf("%.1f%%", percent)
}

func (r *barGraphCanvasRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *barGraphCanvasRenderer) Destroy() {
}

// ShowStatisticsWindow creates and shows a statistics window for the given expression
func ShowStatisticsWindow(expression string) {
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		fmt.Printf("Error calculating statistics: %v\n", err)
		return
	}

	// Create the bar graph
	graph := newBarGraphCanvas(stats)

	// Create and show the window
	window := fyne.CurrentApp().NewWindow("Statistics: " + expression)
	window.SetContent(container.NewMax(graph))
	window.Resize(fyne.NewSize(900, 550))
	window.Show()
}
