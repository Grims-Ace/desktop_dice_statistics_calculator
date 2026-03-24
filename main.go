package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dice Statistics Calculator")

	var calculations []*calculation
	var historyList *widget.List

	// Dice input bar
	inputTextSize := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText) * 2
	diceInputEntry := newCustomEntry(inputTextSize)
	diceInputEntry.SetPlaceHolder("e.g., 2d20H, 3d6+5")
	diceInputContainer := container.NewThemeOverride(diceInputEntry, newSizeOverrideTheme(fyne.CurrentApp().Settings().Theme(), inputTextSize))

	historyList = widget.NewList(
		func() int {
			return len(calculations)
		},
		func() fyne.CanvasObject {
			return newHistoryItemWithCallback(&calculation{}, func(equation string) {
				diceInputEntry.SetText(equation)
			})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*historyItem).SetCalculation(calculations[i])
		},
	)

	roll := func() {
		diceInput := strings.TrimSpace(diceInputEntry.Text)
		if diceInput == "" {
			return
		}

		result, diceRolls, err := CalculateDice(diceInput)
		if err != nil {
			// TODO: show error to user
		} else {
			id := 0
			if len(calculations) > 0 {
				id = calculations[len(calculations)-1].id + 1
			}
			c := &calculation{
				id:        id,
				equation:  diceInput,
				diceRolls: diceRolls,
				result:    fmt.Sprintf("= %s", strconv.FormatFloat(result, 'g', -1, 64)),
			}
			calculations = append([]*calculation{c}, calculations...)
			historyList.Refresh()
		}
	}

	clearHistory := func() {
		calculations = nil
		diceInputEntry.SetText("")
		historyList.Refresh()
	}

	// Roll button
	rollButton := newCustomButton2WithImportance("ROLL", widget.HighImportance, roll)

	diceInputEntry.OnSubmitted = func(s string) {
		roll()
	}

	graphButton := newCustomButton2WithColors("📊", color.NRGBA{R: 46, G: 160, B: 67, A: 255}, color.White, func() {
		diceInput := strings.TrimSpace(diceInputEntry.Text)
		if diceInput == "" {
			return
		}
		ShowStatisticsWindow(diceInput)
		roll()
	})

	lButton := newCustomButton2("L", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "L")
	})
	hButton := newCustomButton2("H", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "H")
	})
	dxButton := newCustomButton2("dX", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d")
	})
	d4Button := newCustomButton2("d4", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d4")
	})
	d6Button := newCustomButton2("d6", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d6")
	})
	d8Button := newCustomButton2("d8", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d8")
	})
	d10Button := newCustomButton2("d10", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d10")
	})
	d12Button := newCustomButton2("d12", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d12")
	})
	d20Button := newCustomButton2("d20", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d20")
	})
	d100Button := newCustomButton2("d100", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "d100")
	})
	sevenButton := createCalcButton("7", diceInputEntry)
	eightButton := createCalcButton("8", diceInputEntry)
	nineButton := createCalcButton("9", diceInputEntry)
	fourButton := createCalcButton("4", diceInputEntry)
	fiveButton := createCalcButton("5", diceInputEntry)
	sixButton := createCalcButton("6", diceInputEntry)
	oneButton := createCalcButton("1", diceInputEntry)
	twoButton := createCalcButton("2", diceInputEntry)
	threeButton := createCalcButton("3", diceInputEntry)
	zeroButton := createCalcButton("0", diceInputEntry)
	multiplyButton := newCustomButton2("*", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "*")
	})
	backspaceButton := newCustomButton2("BKSP", func() {
		text := diceInputEntry.Text
		if len(text) > 0 {
			diceInputEntry.SetText(text[:len(text)-1])
		}
	})
	divideButton := newCustomButton2("/", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "/")
	})
	clearButton := newCustomButton2("C", func() {
		diceInputEntry.SetText("")
	})
	addButton := newCustomButton2("+", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "+")
	})
	acButton := newCustomButton2WithColors("AC", color.NRGBA{R: 198, G: 40, B: 40, A: 255}, color.White, clearHistory)
	decimalButton := newCustomButton2(".", func() {
		diceInputEntry.SetText(diceInputEntry.Text + ".")
	})
	powerButton := newCustomButton2("^", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "^")
	})
	subtractButton := newCustomButton2("-", func() {
		diceInputEntry.SetText(diceInputEntry.Text + "-")
	})

	keypadButtons := []fyne.CanvasObject{
		lButton, hButton, dxButton, d4Button, d6Button, graphButton,
		d8Button, d10Button, d12Button, d20Button, d100Button,
		sevenButton, eightButton, nineButton, multiplyButton, backspaceButton, rollButton,
		fourButton, fiveButton, sixButton, divideButton, clearButton,
		oneButton, twoButton, threeButton, addButton, acButton,
		decimalButton, zeroButton, powerButton, subtractButton,
	}

	buttonSpans := map[fyne.CanvasObject]gridSpan{
		lButton:         {row: 0, col: 0, rowSpan: 1, colSpan: 1},
		hButton:         {row: 0, col: 1, rowSpan: 1, colSpan: 1},
		dxButton:        {row: 0, col: 2, rowSpan: 1, colSpan: 1},
		d4Button:        {row: 0, col: 3, rowSpan: 1, colSpan: 1},
		d6Button:        {row: 0, col: 4, rowSpan: 1, colSpan: 1},
		graphButton:     {row: 0, col: 5, rowSpan: 2, colSpan: 1},
		d8Button:        {row: 1, col: 0, rowSpan: 1, colSpan: 1},
		d10Button:       {row: 1, col: 1, rowSpan: 1, colSpan: 1},
		d12Button:       {row: 1, col: 2, rowSpan: 1, colSpan: 1},
		d20Button:       {row: 1, col: 3, rowSpan: 1, colSpan: 1},
		d100Button:      {row: 1, col: 4, rowSpan: 1, colSpan: 1},
		sevenButton:     {row: 2, col: 0, rowSpan: 1, colSpan: 1},
		eightButton:     {row: 2, col: 1, rowSpan: 1, colSpan: 1},
		nineButton:      {row: 2, col: 2, rowSpan: 1, colSpan: 1},
		multiplyButton:  {row: 2, col: 3, rowSpan: 1, colSpan: 1},
		backspaceButton: {row: 2, col: 4, rowSpan: 1, colSpan: 1},
		rollButton:      {row: 2, col: 5, rowSpan: 4, colSpan: 1},
		fourButton:      {row: 3, col: 0, rowSpan: 1, colSpan: 1},
		fiveButton:      {row: 3, col: 1, rowSpan: 1, colSpan: 1},
		sixButton:       {row: 3, col: 2, rowSpan: 1, colSpan: 1},
		divideButton:    {row: 3, col: 3, rowSpan: 1, colSpan: 1},
		clearButton:     {row: 3, col: 4, rowSpan: 1, colSpan: 1},
		oneButton:       {row: 4, col: 0, rowSpan: 1, colSpan: 1},
		twoButton:       {row: 4, col: 1, rowSpan: 1, colSpan: 1},
		threeButton:     {row: 4, col: 2, rowSpan: 1, colSpan: 1},
		addButton:       {row: 4, col: 3, rowSpan: 1, colSpan: 1},
		acButton:        {row: 4, col: 4, rowSpan: 2, colSpan: 1},
		decimalButton:   {row: 5, col: 0, rowSpan: 1, colSpan: 1},
		zeroButton:      {row: 5, col: 1, rowSpan: 1, colSpan: 1},
		powerButton:     {row: 5, col: 2, rowSpan: 1, colSpan: 1},
		subtractButton:  {row: 5, col: 3, rowSpan: 1, colSpan: 1},
	}

	buttonsContainer := container.New(newSpanGridLayout(6, 6, buttonSpans))
	for _, button := range keypadButtons {
		buttonsContainer.Add(button)
	}

	topContent := container.NewBorder(
		diceInputContainer,
		nil,
		nil,
		nil,
		historyList,
	)

	split := container.NewVSplit(topContent, buttonsContainer)
	split.Offset = 0.5 // Start with a 50/50 split

	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func createCalcButton(label string, entry *customEntry) *customButton2 {
	return newCustomButton2(label, func() {
		entry.SetText(entry.Text + label)
	})
}
