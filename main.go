package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

type Calc struct {
	window   fyne.Window
	display  *widget.RichText
	expr     *widget.RichText
	history  *widget.List
	histData []string
	current  string
	fresh    bool
	mem      float64
	lastOp   string
	lastRight string
}

func (c *Calc) updateDisplay() {
	if c.current == T("Error") || c.current == T("CannotDivideByZero") {
		c.display.ParseMarkdown(`<span style="color: #ff5555"># ` + c.current + `</span>`)
	} else if c.current == "" {
		c.display.ParseMarkdown("# 0")
	} else {
		screenStr := strings.ReplaceAll(c.current, "*", "×")
		screenStr = strings.ReplaceAll(screenStr, "/", "÷")
		c.display.ParseMarkdown("# " + screenStr)
	}
}

func (c *Calc) updateExpr(t string) {
	if t == "" {
		c.expr.ParseMarkdown(" ")
	} else {
		c.expr.ParseMarkdown("### _" + t + "_")
	}
}

func (c *Calc) copyToClipboard(text string) {
	if text == "" || text == T("Error") || text == T("CannotDivideByZero") {
		return
	}
	c.window.Clipboard().SetContent(text)
	c.updateExpr(T("Copied") + ": " + text)

	go func() {
		time.Sleep(1 * time.Second)
		fyne.Do(func() { c.updateExpr("") })
	}()
}

func (c *Calc) pasteFromClipboard() {
	clipText := strings.TrimSpace(c.window.Clipboard().Content())
	if clipText == "" {
		return
	}
	if c.fresh || c.current == "0" {
		c.current = clipText
	} else {
		c.current += clipText
	}
	c.fresh = false
	c.updateDisplay()
}

func (c *Calc) logToHistory(expression string, result string) {
	exprNice := strings.ReplaceAll(expression, "*", "×")
	exprNice = strings.ReplaceAll(exprNice, "/", "÷")
	entry := exprNice + " = " + result

	if len(c.histData) == 1 && c.histData[0] == T("EmptyHistory") {
		c.histData = []string{entry}
	} else {
		c.histData = append([]string{entry}, c.histData...)
	}
	c.history.Refresh()
}

func (c *Calc) appendChar(char string) {
	if c.fresh {
		if char == "+" || char == "-" || char == "*" || char == "/" {
			c.fresh = false
		} else {
			c.current = ""
			c.fresh = false
		}
	}
	if c.current == "0" && char != "." && char != "+" && char != "-" && char != "*" && char != "/" {
		c.current = char
	} else {
		c.current += char
	}
	c.updateDisplay()
}

func (c *Calc) setupKeyboard() {
	c.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyEnter, fyne.KeyReturn:
			c.calc()
		case fyne.KeyBackspace:
			c.backspace()
		case fyne.KeyEscape:
			c.clearAll()
		case fyne.Key0, fyne.Key1, fyne.Key2, fyne.Key3, fyne.Key4,
			fyne.Key5, fyne.Key6, fyne.Key7, fyne.Key8, fyne.Key9:
			c.appendChar(string(key.Name))
		case fyne.KeyPeriod, fyne.KeyComma:
			c.appendChar(".")
		case fyne.KeyPlus:
			c.appendChar("+")
		case fyne.KeyMinus:
			c.appendChar("-")
		case fyne.KeyAsterisk:
			c.appendChar("*")
		case fyne.KeySlash:
			c.appendChar("/")
		case fyne.KeyEqual:
			c.calc()
		}
	})
}

func (c *Calc) evalCurrent() float64 {
	if c.current == "" || c.current == T("Error") || c.current == T("CannotDivideByZero") {
		return 0
	}
	program, err := expr.Compile(c.current)
	if err != nil {
		return 0
	}
	result, err := expr.Run(program, nil)
	if err != nil {
		return 0
	}
	if val, ok := result.(float64); ok {
		return val
	}
	if valInt, ok := result.(int); ok {
		return float64(valInt)
	}
	return 0
}

func (c *Calc) calc() {
	if c.current == "" || c.current == T("Error") || c.current == T("CannotDivideByZero") {
		return
	}

	if !c.fresh {
		opIdx := findLastBinaryOp(c.current)
		if opIdx != -1 {
			right := strings.TrimSpace(c.current[opIdx+1:])
			if right != "" {
				c.lastOp = string(c.current[opIdx])
				c.lastRight = right
			}
		}
	}

	if c.fresh && c.lastOp != "" {
		if _, err := strconv.ParseFloat(c.current, 64); err == nil {
			c.current = c.current + c.lastOp + c.lastRight
		}
	}

	tree, err := parser.Parse(c.current)
	if err != nil {
		c.current = T("Error")
		c.fresh = true
		c.updateDisplay()
		return
	}

	visitor := &ZeroDivVisitor{}
	ast.Walk(&tree.Node, visitor)

	if visitor.err != nil {
		c.current = T("CannotDivideByZero")
		c.fresh = true
		c.updateDisplay()
		return
	}

	program, err := expr.Compile(c.current)
	if err != nil {
		c.current = T("Error")
		c.fresh = true
		c.updateDisplay()
		return
	}

	result, err := expr.Run(program, nil)
	if err != nil {
		c.current = T("Error")
		c.fresh = true
		c.updateDisplay()
		return
	}

	var val float64
	if resFloat, ok := result.(float64); ok {
		val = resFloat
	} else if resInt, ok := result.(int); ok {
		val = float64(resInt)
	} else {
		c.current = T("Error")
		c.fresh = true
		c.updateDisplay()
		return
	}

	resStr := fmtResult(val)
	c.logToHistory(c.current, resStr)

	c.current = resStr
	c.fresh = true
	c.updateDisplay()
}

func (c *Calc) clearAll() {
	c.current = "0"
	c.fresh = true
	c.updateExpr("")
	c.updateDisplay()
}

func (c *Calc) backspace() {
	if c.fresh || len(c.current) <= 1 {
		c.current = "0"
		c.fresh = true
	} else {
		c.current = c.current[:len(c.current)-1]
	}
	c.updateDisplay()
}

func (c *Calc) toggleSign() {
	if c.current == "0" || c.fresh {
		return
	}
	if _, err := strconv.ParseFloat(c.current, 64); err == nil {
		if strings.HasPrefix(c.current, "-") {
			c.current = c.current[1:]
		} else {
			c.current = "-" + c.current
		}
		c.updateDisplay()
	}
}

func findLastBinaryOp(s string) int {
	depth := 0
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case ')':
			depth++
		case '(':
			depth--
		case '+':
			if depth == 0 {
				return i
			}
		case '-':
			if depth == 0 && i > 0 && s[i-1] != '(' && s[i-1] != '+' && s[i-1] != '-' && s[i-1] != '*' && s[i-1] != '/' {
				return i
			}
		case '*', '/':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func evalStr(s string) (float64, error) {
	program, err := expr.Compile(s)
	if err != nil {
		return 0, err
	}
	result, err := expr.Run(program, nil)
	if err != nil {
		return 0, err
	}
	switch v := result.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	}
	return 0, fmt.Errorf("not a number")
}

func (c *Calc) percent() {
	if c.current == "" || c.current == T("Error") || c.current == T("CannotDivideByZero") {
		return
	}

	originalExpr := c.current
	opIdx := findLastBinaryOp(c.current)

	if opIdx == -1 {
		val := c.evalCurrent()
		c.current = fmtResult(val / 100)
	} else {
		op := string(c.current[opIdx])
		leftPart := c.current[:opIdx]
		rightStr := strings.TrimSpace(c.current[opIdx+1:])

		rightVal, err := strconv.ParseFloat(rightStr, 64)
		if err != nil {
			val := c.evalCurrent()
			c.current = fmtResult(val / 100)
		} else if op == "+" || op == "-" {
			leftVal, err := evalStr(leftPart)
			if err != nil {
				val := c.evalCurrent()
				c.current = fmtResult(val / 100)
			} else {
				pct := leftVal * rightVal / 100
				c.current = leftPart + op + fmtResult(pct)
			}
		} else {
			c.current = leftPart + op + fmtResult(rightVal/100)
		}
	}

	c.logToHistory(originalExpr+" %", c.current)
	c.fresh = true
	c.updateDisplay()
}

func (c *Calc) singleOp(op string) {
	val := c.evalCurrent()
	var res float64
	switch op {
	case "1/x":
		if val == 0 {
			c.current = T("CannotDivideByZero")
			c.fresh = true
			c.updateDisplay()
			return
		}
		res = 1 / val
	case "x²":
		res = val * val
	case "√":
		if val < 0 {
			c.current = T("Error")
			c.fresh = true
			c.updateDisplay()
			return
		}
		res = math.Sqrt(val)
	}
	c.current = fmtResult(res)
	c.fresh = true
	c.updateDisplay()
}

func (c *Calc) mc()   { c.mem = 0; c.updateExpr(T("MemoryCleared")) }
func (c *Calc) mr()   { c.current = fmtResult(c.mem); c.fresh = true; c.updateDisplay() }
func (c *Calc) mAdd() { c.mem += c.evalCurrent(); c.fresh = true; c.updateExpr(T("AddedToMemory")) }
func (c *Calc) mSub() { c.mem -= c.evalCurrent(); c.fresh = true; c.updateExpr(T("SubtractedFromMemory")) }

func fmtResult(v float64) string {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return T("Error")
	}
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if len(s) > 16 {
		s = strconv.FormatFloat(v, 'g', 10, 64)
	}
	return s
}

func main() {
	initLocale()

	a := app.New()
	a.Settings().SetTheme(&MaterialTheme{})

	w := a.NewWindow(T("WindowTitle"))
	w.Resize(fyne.NewSize(680, 680))
	w.SetIcon(theme.ComputerIcon())

	c := &Calc{
		window:   w,
		display:  widget.NewRichTextFromMarkdown("# 0"),
		expr:     widget.NewRichTextFromMarkdown(" "),
		histData: []string{T("EmptyHistory")},
		fresh:    true,
	}

	c.history = widget.NewList(
		func() int { return len(c.histData) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(c.histData[i]) },
	)

	c.history.OnSelected = func(id widget.ListItemID) {
		defer c.history.Unselect(id)
		str := c.histData[id]
		if str == T("EmptyHistory") {
			return
		}
		parts := strings.Split(str, "=")
		if len(parts) > 1 {
			expr := strings.TrimSpace(parts[0])
			expr = strings.ReplaceAll(expr, "×", "*")
			expr = strings.ReplaceAll(expr, "÷", "/")
			c.current = expr
			c.fresh = false
			c.updateDisplay()
		}
	}

	clearHistBtn := widget.NewButton(T("Clear"), func() {
		c.histData = []string{T("EmptyHistory")}
		c.history.Refresh()
	})
	clearHistBtn.Importance = widget.LowImportance

	pasteBtn := widget.NewButton(T("Paste")+" 📋", c.pasteFromClipboard)
	pasteBtn.Importance = widget.LowImportance

	historyButtons := container.NewGridWithColumns(2, pasteBtn, clearHistBtn)
	historyPanel := container.NewBorder(nil, historyButtons, nil, nil, c.history)

	displayButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		c.copyToClipboard(c.current)
	})
	displayButton.Importance = widget.LowImportance

	top := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), c.expr),
		container.NewBorder(nil, nil, nil, displayButton, container.NewHBox(layout.NewSpacer(), c.display)),
	)

	createBtn := func(txt string, importance widget.Importance, fn func()) *widget.Button {
		b := widget.NewButton(txt, fn)
		b.Importance = importance
		return b
	}

	grid := container.NewGridWithColumns(5)
	layoutPlan := []struct {
		txt string
		imp widget.Importance
		fn  func()
	}{
		{"MC", widget.LowImportance, c.mc},
		{"MR", widget.LowImportance, c.mr},
		{"M+", widget.LowImportance, c.mAdd},
		{"M-", widget.LowImportance, c.mSub},
		{"C", widget.LowImportance, c.clearAll},

		{"(", widget.LowImportance, func() { c.appendChar("(") }},
		{")", widget.LowImportance, func() { c.appendChar(")") }},
		{"1/x", widget.MediumImportance, func() { c.singleOp("1/x") }},
		{"x²", widget.MediumImportance, func() { c.singleOp("x²") }},
		{"√x", widget.MediumImportance, func() { c.singleOp("√") }},

		{"7", widget.MediumImportance, func() { c.appendChar("7") }},
		{"8", widget.MediumImportance, func() { c.appendChar("8") }},
		{"9", widget.MediumImportance, func() { c.appendChar("9") }},
		{"÷", widget.LowImportance, func() { c.appendChar("/") }},
		{"⌫", widget.LowImportance, c.backspace},

		{"4", widget.MediumImportance, func() { c.appendChar("4") }},
		{"5", widget.MediumImportance, func() { c.appendChar("5") }},
		{"6", widget.MediumImportance, func() { c.appendChar("6") }},
		{"×", widget.LowImportance, func() { c.appendChar("*") }},
		{"%", widget.MediumImportance, c.percent},

		{"1", widget.MediumImportance, func() { c.appendChar("1") }},
		{"2", widget.MediumImportance, func() { c.appendChar("2") }},
		{"3", widget.MediumImportance, func() { c.appendChar("3") }},
		{"-", widget.LowImportance, func() { c.appendChar("-") }},
		{"π", widget.MediumImportance, func() { c.appendChar("3.1415926535") }},

		{"0", widget.MediumImportance, func() { c.appendChar("0") }},
		{".", widget.MediumImportance, func() { c.appendChar(".") }},
		{"+/-", widget.MediumImportance, c.toggleSign},
		{"+", widget.LowImportance, func() { c.appendChar("+") }},
		{"=", widget.HighImportance, func() { c.calc() }},
	}

	for _, r := range layoutPlan {
		if r.txt == "" {
			grid.Add(widget.NewLabel(""))
			continue
		}
		grid.Add(createBtn(r.txt, r.imp, r.fn))
	}

	calcContainer := container.NewBorder(top, nil, nil, nil, grid)
	mainLayout := container.NewHSplit(calcContainer, historyPanel)
	mainLayout.Offset = 0.65

	w.SetContent(mainLayout)
	c.setupKeyboard()

	w.ShowAndRun()
}