package main

import (
	"collatz/drawing"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
)

type Entry interface {
	fyne.Disableable
	fyne.Focusable
	fyne.Tappable
	fyne.Widget
	desktop.Mouseable
	desktop.Keyable
	mobile.Keyboardable
	mobile.Touchable
	fyne.Tabbable
}

type Field struct {
	LabelText string
	Entry     Entry
}

func (f Field) Label() *widget.Label {
	return widget.NewLabel(f.LabelText)
}

func NewField(label string, initialValue string, validator func(string) error, updateParam func(string)) Field {
	entry := widget.NewEntry()
	entry.SetText(initialValue)
	updateParam(initialValue)
	entry.OnChanged = updateParam
	entry.Validator = validator
	return Field{
		LabelText: label,
		Entry:     entry,
	}
}

func NewSelectField(label string, initialValue string, validator func(string) error, updateParam func(string), values ...string) Field {
	entry := widget.NewSelectEntry(values)
	entry.SetText(initialValue)
	updateParam(initialValue)
	return Field{
		LabelText: label,
		Entry:     entry,
	}
}

type FieldTab struct {
	Title   string
	Columns int
	Fields  []Field
}

func (t FieldTab) TabItem() *container.TabItem {
	objects := make([]fyne.CanvasObject, 0, 2*len(t.Fields))
	for i := range t.Fields {
		objects = append(objects, t.Fields[i].Label(), t.Fields[i].Entry)
	}
	return container.NewTabItem(t.Title, container.NewGridWithColumns(t.Columns, objects...))
}

type ColorPicker struct {
	LabelText  string
	Button     *widget.Button
	Rectangle  *canvas.Rectangle
	ColorLabel *canvas.Text
}

func (c ColorPicker) Label() *widget.Label {
	return widget.NewLabel(c.LabelText)
}

type ColorTab struct {
	Title        string
	ColorPickers []ColorPicker
}

func (t ColorTab) TabItem() *container.TabItem {
	objects := make([]fyne.CanvasObject, 0, 2*len(t.ColorPickers))
	whiteRectange := canvas.NewRectangle(color.White)
	for i := range t.ColorPickers {
		button := container.New(layout.NewMaxLayout(), whiteRectange, t.ColorPickers[i].Rectangle, t.ColorPickers[i].ColorLabel, t.ColorPickers[i].Button)
		objects = append(objects, t.ColorPickers[i].Label(), button)
	}
	return container.NewTabItem(t.Title, container.NewGridWithColumns(2, objects...))
}

func NewColorPicker(w fyne.Window, label string, initialColor color.Color, updateParam func(color.Color)) ColorPicker {
	colorLabel := canvas.NewText("", color.Black)
	colorLabel.TextStyle = fyne.TextStyle{
		Monospace: true,
	}
	updateLabelFromColor(colorLabel, initialColor)
	rectangle := canvas.NewRectangle(initialColor)
	colorPicker := colorpicker.New(300, colorpicker.StyleHue)
	updateParam(initialColor)
	colorPicker.SetColor(initialColor)
	colorPicker.SetOnChanged(func(c color.Color) {
		updateParam(c)
		rectangle.FillColor = c
		rectangle.Refresh()
		updateLabelFromColor(colorLabel, c)
	})
	return ColorPicker{
		LabelText:  label,
		Rectangle:  rectangle,
		ColorLabel: colorLabel,
		Button: widget.NewButton("", func() {
			dialog.ShowCustom("Pick Color", "Ok", container.NewWithoutLayout(colorPicker), w)
		}),
	}
}

func main() {
	a := app.NewWithID("app")
	w := a.NewWindow("Collatz")
	params := &drawing.Params{}
	var img image.Image
	var err error

	fieldTabs := []FieldTab{
		{
			Title:   "General",
			Columns: 4,
			Fields: []Field{
				NewField("Filename", "poster", isNonEmptyString, func(s string) {
					params.Filename = s
				}),
				NewSelectField("Format", "png", isNonEmptyString, func(s string) {
					params.Format = drawing.ImageFormat(s)
				}, string(drawing.ImageFormatPNG), string(drawing.ImageFormatJPG)),
				NewField("# of Lines", "3000", isInteger, func(s string) {
					params.N, _ = strconv.Atoi(s)
				}),
				NewField("Max Value", "200000000", isInteger, func(s string) {
					params.Max, _ = strconv.Atoi(s)
				}),
				NewField("Width", "15118", isInteger, func(s string) {
					params.Width, _ = strconv.Atoi(s)
				}),
				NewField("Height", "8504", isInteger, func(s string) {
					params.Height, _ = strconv.Atoi(s)
				}),
			},
		},
		{
			Title:   "Tilt",
			Columns: 2,
			Fields: []Field{
				NewField("Angle", "8.0", isFloat, func(s string) {
					params.TiltAngle, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("Angle Even Adjuster", strconv.FormatFloat(math.Log(2), 'f', -1, 64), isFloat, func(s string) {
					params.TiltAngleEvenAdjuster, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("Angle Odd Adjuster", strconv.FormatFloat(math.Log(3)*1.12, 'f', -1, 64), isFloat, func(s string) {
					params.TiltAngleOddAdjuster, _ = strconv.ParseFloat(s, 64)
				}),
			},
		},
		{
			Title:   "Line",
			Columns: 2,
			Fields: []Field{
				NewField("Width", "7", isFloat, func(s string) {
					params.LineWidth, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("Length Multiplier", "260", isFloat, func(s string) {
					params.LineLengthMultiplier, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("Shortening Power", "1.21", isFloat, func(s string) {
					params.LineShorteningPower, _ = strconv.ParseFloat(s, 64)
				}),
			},
		},
		{
			Title:   "Grid",
			Columns: 2,
			Fields: []Field{
				NewField("Line Width", "4", isFloat, func(s string) {
					params.GridLineWidth, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("# Columns", "64", isInteger, func(s string) {
					params.GridColumns, _ = strconv.Atoi(s)
				}),
				NewField("# Rows", "36", isInteger, func(s string) {
					params.GridRows, _ = strconv.Atoi(s)
				}),
			},
		},
		{
			Title:   "Start",
			Columns: 2,
			Fields: []Field{
				NewField("Angle", "", isFloat, func(s string) {
					params.StartAngle, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("X Position", "0.0", isFloat, func(s string) {
					params.StartX, _ = strconv.ParseFloat(s, 64)
				}),
				NewField("Y Position", strconv.FormatFloat(float64(8504)*0.80, 'f', -1, 64), isFloat, func(s string) {
					params.StartY, _ = strconv.ParseFloat(s, 64)
				}),
			},
		},
	}

	colorTabs := []ColorTab{
		{
			Title: "Colors",
			ColorPickers: []ColorPicker{
				NewColorPicker(w, "Line Color", color.NRGBA{255, 230, 191, 60}, func(c color.Color) {
					params.LineColor = c
				}),
				NewColorPicker(w, "Text Color", color.NRGBA{255, 230, 191, 150}, func(c color.Color) {
					params.TextColor = c
				}),
				NewColorPicker(w, "Grid Color", color.NRGBA{7, 67, 110, 255}, func(c color.Color) {
					params.GridColor = c
				}),
				NewColorPicker(w, "Background Color", color.NRGBA{12, 36, 77, 255}, func(c color.Color) {
					params.BackgroundColor = c
				}),
			},
		},
		{
			Title: "Gradient",
			ColorPickers: []ColorPicker{
				NewColorPicker(w, "Color #1", color.NRGBA{12, 36, 77, 255}, func(c color.Color) {
					if len(params.GradientColors) < 1 {
						params.GradientColors = append(params.GradientColors, make([]color.Color, 1-len(params.GradientColors))...)
					}
					params.GradientColors[0] = c
				}),
				NewColorPicker(w, "Color #2", color.NRGBA{0, 0, 0, 255}, func(c color.Color) {
					if len(params.GradientColors) < 2 {
						params.GradientColors = append(params.GradientColors, make([]color.Color, 2-len(params.GradientColors))...)
					}
					params.GradientColors[1] = c
				}),
				NewColorPicker(w, "Color #3", nil, func(c color.Color) {
					if len(params.GradientColors) < 3 {
						params.GradientColors = append(params.GradientColors, make([]color.Color, 3-len(params.GradientColors))...)
					}
					params.GradientColors[2] = c
				}),
				NewColorPicker(w, "Color #4", nil, func(c color.Color) {
					if len(params.GradientColors) < 4 {
						params.GradientColors = append(params.GradientColors, make([]color.Color, 4-len(params.GradientColors))...)
					}
					params.GradientColors[3] = c
				}),
			},
		},
	}

	tabItems := make([]*container.TabItem, 0, len(fieldTabs)+len(colorTabs))
	for i := range fieldTabs {
		tabItems = append(tabItems, fieldTabs[i].TabItem())
	}
	for i := range colorTabs {
		tabItems = append(tabItems, colorTabs[i].TabItem())
	}

	box := container.NewVBox(
		container.NewMax(
			container.NewAppTabs(tabItems...),
		),
	)
	progress := widget.NewProgressBarInfinite()
	progress.Stop()
	progress.Hide()

	imageBox := canvas.NewImageFromImage(img)
	imageBox.FillMode = canvas.ImageFillContain
	imageBox.SetMinSize(fyne.NewSize(900, 500))

	var overwrite bool
	overwriteCheck := widget.NewCheck("Overwrite previously generated image", func(b bool) {
		overwrite = b
	})

	generateButton := widget.NewButton("Generate", func() {})
	generateButton.OnTapped = func() {
		generateButton.Disable()
		generateButton.Hide()
		overwriteCheck.Disable()
		overwriteCheck.Hide()
		progress.Show()
		progress.Start()
		for i := range fieldTabs {
			for j := range fieldTabs[i].Fields {
				fieldTabs[i].Fields[j].Entry.Disable()
			}
		}

		defer func() {
			progress.Stop()
			progress.Hide()
			generateButton.Show()
			generateButton.Enable()
			overwriteCheck.Show()
			overwriteCheck.Enable()
			for i := range fieldTabs {
				for j := range fieldTabs[i].Fields {
					fieldTabs[i].Fields[j].Entry.Enable()
				}
			}
		}()

		imageBox.Image, err = drawing.Draw(params, true, overwrite)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		imageBox.Refresh()
	}
	box.Add(container.NewCenter(container.NewHBox(overwriteCheck, generateButton)))
	box.Add(progress)
	box.Add(container.NewCenter(imageBox))

	w.SetContent(box)
	w.ShowAndRun()
}

func isNonEmptyString(s string) error {
	if s == "" {
		return errors.New("must be a non-empty string")
	}
	return nil
}

func isInteger(s string) error {
	if s != "" {
		_, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("must be an integer")
		}
	}
	return nil
}

func isFloat(s string) error {
	if s != "" {
		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return errors.New("must be a float")
		}
	}
	return nil
}

func updateLabelFromColor(label *canvas.Text, c color.Color) {
	if c != nil {
		rgba := c.(color.NRGBA)
		luminance := (0.299*float64(rgba.R) + 0.587*float64(rgba.G) + 0.114*float64(rgba.B)) / 255.0
		transparency := float64(rgba.A) / 255.0
		if luminance > 0.5 || transparency < 0.5 {
			label.Color = color.Black
		} else {
			label.Color = color.White
		}
		label.Text = fmt.Sprintf(" #%.2X%.2X%.2X%.2X - R: %d, G: %d, B: %d, A: %d", rgba.R, rgba.G, rgba.B, rgba.A, rgba.R, rgba.G, rgba.B, rgba.A)
		label.Refresh()
	}
}
