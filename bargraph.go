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

// barGraphCanvas is a custom widget that renders a bar graph.
type barGraphCanvas struct {
	widget.BaseWidget
	stats *DiceStatistics
}

func newBarGraphCanvas(stats *DiceStatistics) *barGraphCanvas {
	graph := &barGraphCanvas{stats: stats}
	graph.ExtendBaseWidget(graph)
	return graph
}

func (b *barGraphCanvas) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	background := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 255})
	yAxisLine := canvas.NewLine(color.White)
	yAxisLine.StrokeWidth = 2
	xAxisLine := canvas.NewLine(color.White)
	xAxisLine.StrokeWidth = 2
	title := canvas.NewText("Probability Distribution", color.White)
	statsLine1 := canvas.NewText("", color.White)
	statsLine2 := canvas.NewText("", color.White)
	yLabel := canvas.NewText("Probability (%)", color.White)
	xLabel := canvas.NewText("Result Value", color.White)

	renderer := &barGraphCanvasRenderer{
		graph:      b,
		background: background,
		yAxisLine:  yAxisLine,
		xAxisLine:  xAxisLine,
		title:      title,
		statsLine1: statsLine1,
		statsLine2: statsLine2,
		yLabel:     yLabel,
		xLabel:     xLabel,
		objects: []fyne.CanvasObject{
			background,
			yAxisLine,
			xAxisLine,
			title,
			statsLine1,
			statsLine2,
			yLabel,
			xLabel,
		},
	}

	renderer.Refresh()
	return renderer
}

type barGraphCanvasRenderer struct {
	graph       *barGraphCanvas
	background  *canvas.Rectangle
	yAxisLine   *canvas.Line
	xAxisLine   *canvas.Line
	title       *canvas.Text
	statsLine1  *canvas.Text
	statsLine2  *canvas.Text
	yLabel      *canvas.Text
	xLabel      *canvas.Text
	yTicks      []*canvas.Line
	yTickLabels []*canvas.Text
	bars        []*canvas.Rectangle
	barLabels   []*canvas.Text
	objects     []fyne.CanvasObject
}

func (r *barGraphCanvasRenderer) Layout(size fyne.Size) {
	r.layout(size)
}

func (r *barGraphCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(320, 240)
}

func (r *barGraphCanvasRenderer) Refresh() {
	stats := r.graph.stats
	if stats == nil || len(stats.Results) == 0 {
		r.hideAllButBackground()
		r.background.Show()
		r.background.Refresh()
		return
	}

	r.background.Show()
	r.syncData(stats)
	size := r.graph.Size()
	if size.Width == 0 || size.Height == 0 {
		size = fyne.NewSize(900, 550)
	}
	r.layout(size)
}

func (r *barGraphCanvasRenderer) syncData(stats *DiceStatistics) {
	r.title.Show()
	r.statsLine1.Show()
	r.statsLine2.Show()
	r.yLabel.Show()
	r.xLabel.Show()
	r.yAxisLine.Show()
	r.xAxisLine.Show()

	r.statsLine1.Text = fmt.Sprintf("Range: %d to %d  |  Total Outcomes: %s", stats.MinValue, stats.MaxValue, stats.TotalOutcomesText)
	r.statsLine2.Text = fmt.Sprintf("Average: %.2f  |  Most Common: %d", stats.Average, stats.MostCommon)

	axisMaxPercent, tickStep := calculateYAxisScale(stats.GetMaxPercentage())
	numYTicks := int(math.Round(axisMaxPercent / tickStep))
	r.ensureYTicks(numYTicks + 1)
	for i := 0; i <= numYTicks; i++ {
		percent := float64(i) * tickStep
		r.yTicks[i].Show()
		r.yTickLabels[i].Text = formatPercentLabel(percent)
		r.yTickLabels[i].Show()
	}
	for i := numYTicks + 1; i < len(r.yTicks); i++ {
		r.yTicks[i].Hide()
		r.yTickLabels[i].Hide()
	}

	outcomes := stats.GetSortedOutcomes()
	r.ensureBars(len(outcomes))
	for i := range outcomes {
		r.bars[i].Show()
		r.barLabels[i].Hide()
	}
	for i := len(outcomes); i < len(r.bars); i++ {
		r.bars[i].Hide()
		r.barLabels[i].Hide()
	}

	r.refreshObjects()
}

func (r *barGraphCanvasRenderer) layout(size fyne.Size) {
	r.background.Move(fyne.NewPos(0, 0))
	r.background.Resize(size)

	stats := r.graph.stats
	if stats == nil || len(stats.Results) == 0 {
		return
	}

	axisMaxPercent, tickStep := calculateYAxisScale(stats.GetMaxPercentage())
	outcomes := stats.GetSortedOutcomes()

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

	r.yAxisLine.Move(fyne.NewPos(leftPadding, topPadding))
	r.yAxisLine.Resize(fyne.NewSize(0, graphHeight))
	r.xAxisLine.Move(fyne.NewPos(leftPadding, topPadding+graphHeight))
	r.xAxisLine.Resize(fyne.NewSize(graphWidth, 0))

	r.title.TextSize = titleSize
	r.title.Move(fyne.NewPos(leftPadding, 5))

	r.statsLine1.TextSize = bodyTextSize
	r.statsLine1.Move(fyne.NewPos(leftPadding, 5+r.title.MinSize().Height))

	r.statsLine2.TextSize = bodyTextSize
	r.statsLine2.Move(fyne.NewPos(leftPadding, 5+r.title.MinSize().Height+r.statsLine1.MinSize().Height))

	r.yLabel.TextSize = labelTextSize
	r.yLabel.Move(fyne.NewPos(clamp(leftPadding*0.15, 8, 15), topPadding+graphHeight/2-r.yLabel.MinSize().Height/2))

	r.xLabel.TextSize = labelTextSize
	r.xLabel.Move(fyne.NewPos(leftPadding+graphWidth/2-r.xLabel.MinSize().Width/2, topPadding+graphHeight+clamp(bottomPadding*0.35, 18, 28)))

	numYTicks := int(math.Round(axisMaxPercent / tickStep))
	for i := 0; i <= numYTicks && i < len(r.yTicks); i++ {
		percent := float64(i) * tickStep
		yPos := topPadding + graphHeight - (float32(percent/axisMaxPercent) * graphHeight)

		r.yTicks[i].Move(fyne.NewPos(leftPadding-5, yPos))
		r.yTicks[i].Resize(fyne.NewSize(5, 0))
		r.yTickLabels[i].TextSize = tickTextSize
		r.yTickLabels[i].Move(fyne.NewPos(leftPadding-50, yPos-7))
	}

	numBars := len(outcomes)
	if numBars == 0 {
		r.refreshObjects()
		return
	}

	barSpacing := float32(2)
	totalSpacing := float32(numBars+1) * barSpacing
	barWidth := (graphWidth - totalSpacing) / float32(numBars)
	if barWidth < 1 {
		barWidth = 1
	}

	labelStep := calculateLabelStep(graphWidth, numBars)
	for i, value := range outcomes {
		percentage := stats.Percentages[value]
		barHeight := (float32(percentage) / float32(axisMaxPercent)) * graphHeight
		xPos := leftPadding + barSpacing + float32(i)*(barWidth+barSpacing)

		r.bars[i].Move(fyne.NewPos(xPos, topPadding+graphHeight-barHeight))
		r.bars[i].Resize(fyne.NewSize(barWidth, barHeight))

		r.barLabels[i].TextSize = tickTextSize
		r.barLabels[i].Alignment = fyne.TextAlignCenter
		if i == 0 || i == numBars-1 || (labelStep > 0 && i%labelStep == 0 && i < numBars-labelStep) {
			r.barLabels[i].Text = fmt.Sprintf("%d", value)
			r.barLabels[i].Show()
			r.barLabels[i].Move(fyne.NewPos(xPos+barWidth/2-r.barLabels[i].MinSize().Width/2, topPadding+graphHeight+10))
		} else {
			r.barLabels[i].Hide()
		}
	}

	r.refreshObjects()
}

func (r *barGraphCanvasRenderer) ensureYTicks(count int) {
	for len(r.yTicks) < count {
		tick := canvas.NewLine(color.White)
		tick.StrokeWidth = 1
		label := canvas.NewText("", color.White)
		r.yTicks = append(r.yTicks, tick)
		r.yTickLabels = append(r.yTickLabels, label)
		r.objects = append(r.objects, tick, label)
	}
}

func (r *barGraphCanvasRenderer) ensureBars(count int) {
	for len(r.bars) < count {
		bar := canvas.NewRectangle(color.NRGBA{R: 100, G: 180, B: 255, A: 255})
		label := canvas.NewText("", color.White)
		r.bars = append(r.bars, bar)
		r.barLabels = append(r.barLabels, label)
		r.objects = append(r.objects, bar, label)
	}
}

func (r *barGraphCanvasRenderer) availableGraphWidth() float32 {
	size := r.graph.Size()
	if size.Width == 0 {
		size.Width = 900
	}
	leftPadding := clamp(size.Width*0.12, 60, 100)
	rightPadding := clamp(size.Width*0.03, 16, 24)
	return size.Width - leftPadding - rightPadding
}

func (r *barGraphCanvasRenderer) hideAllButBackground() {
	for _, object := range r.objects {
		object.Hide()
	}
}

func (r *barGraphCanvasRenderer) refreshObjects() {
	for _, object := range r.objects {
		object.Refresh()
	}
}

func (r *barGraphCanvasRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *barGraphCanvasRenderer) Destroy() {
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
	labelWidthEstimate := float32(35)
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

// ShowStatisticsWindow creates and shows a statistics window for the given expression.
func ShowStatisticsWindow(expression string) {
	stats, err := CalculateDiceStatistics(expression)
	if err != nil {
		fmt.Printf("Error calculating statistics: %v\n", err)
		return
	}

	graph := newBarGraphCanvas(stats)

	window := fyne.CurrentApp().NewWindow("Statistics: " + expression)
	window.SetContent(container.NewMax(graph))
	window.Resize(fyne.NewSize(900, 550))
	window.Show()
}
