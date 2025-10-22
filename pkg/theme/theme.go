package theme

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

// Aliases for layout
type (
	C = layout.Context
	D = layout.Dimensions
)

// Mode represents theme mode
type Mode int

const (
	Dark Mode = iota
	Light
)

// Theme holds the colors and material base theme
type Theme struct {
	Base *material.Theme
	mode Mode

	// Core colors
	Background    color.NRGBA // main background
	Surface       color.NRGBA // cards, panels
	Foreground    color.NRGBA // primary text/icons
	SecondaryText color.NRGBA // secondary/subtle text
	ButtonColor   color.NRGBA // primary buttons
	AccentColor   color.NRGBA // highlights, links
	BorderColor   color.NRGBA // borders/dividers
	ErrorColor    color.NRGBA // error messages
	SuccessColor  color.NRGBA // success messages

	// Navigation / interaction
	NavHover        color.NRGBA // hover/active nav items
	NavText         color.NRGBA // nav item text
	NavTextInactive color.NRGBA // inactive nav item text

	// Extras
	HintColor     color.NRGBA // placeholder text / hints
	LinkColor     color.NRGBA // hyperlinks
	DisabledColor color.NRGBA // disabled buttons or inputs
}

// NewTheme creates a theme with the given mode
func NewTheme(mode Mode) *Theme {
	t := &Theme{
		Base: material.NewTheme(),
		mode: mode,
	}

	if t.mode == Light {
		t.useLight()
	} else {
		t.useDark()
	}

	return t
}

// ToggleMode switches between light and dark themes
func (t *Theme) ToggleMode() {
	if t.mode == Light {
		t.useDark()
	} else {
		t.useLight()
	}
}

// useDark sets all colors for dark mode
func (t *Theme) useDark() {
	t.Background = color.NRGBA{R: 37, G: 37, B: 37, A: 255}
	t.Surface = color.NRGBA{R: 54, G: 54, B: 54, A: 255}
	t.Foreground = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	t.SecondaryText = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
	t.ButtonColor = color.NRGBA{R: 30, G: 136, B: 229, A: 255}
	t.AccentColor = color.NRGBA{R: 255, G: 152, B: 0, A: 255}
	t.BorderColor = color.NRGBA{R: 68, G: 68, B: 68, A: 255}
	t.ErrorColor = color.NRGBA{R: 229, G: 57, B: 53, A: 255}
	t.SuccessColor = color.NRGBA{R: 67, G: 160, B: 71, A: 255}
	t.NavHover = color.NRGBA{R: 51, G: 51, B: 51, A: 255}
	t.NavText = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	t.NavTextInactive = color.NRGBA{R: 136, G: 136, B: 136, A: 255}
	t.HintColor = color.NRGBA{R: 117, G: 117, B: 117, A: 255}
	t.LinkColor = color.NRGBA{R: 30, G: 136, B: 229, A: 255}
	t.DisabledColor = color.NRGBA{R: 85, G: 85, B: 85, A: 255}

	t.mode = Dark
}

// useLight sets all colors for light mode
func (t *Theme) useLight() {
	t.Background = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	t.Surface = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	t.Foreground = color.NRGBA{R: 33, G: 33, B: 33, A: 255}
	t.SecondaryText = color.NRGBA{R: 117, G: 117, B: 117, A: 255}
	t.ButtonColor = color.NRGBA{R: 33, G: 150, B: 243, A: 255}
	t.AccentColor = color.NRGBA{R: 255, G: 152, B: 0, A: 255}
	t.BorderColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	t.ErrorColor = color.NRGBA{R: 211, G: 47, B: 47, A: 255}
	t.SuccessColor = color.NRGBA{R: 56, G: 142, B: 60, A: 255}
	t.NavHover = color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	t.NavText = color.NRGBA{R: 33, G: 33, B: 33, A: 255}
	t.NavTextInactive = color.NRGBA{R: 117, G: 117, B: 117, A: 255}
	t.HintColor = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
	t.LinkColor = color.NRGBA{R: 33, G: 150, B: 243, A: 255}
	t.DisabledColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255}

	t.mode = Light
}

// PaintRect paints a rectangle of a given size and color
func PaintRect(gtx layout.Context, size image.Point, col color.NRGBA) layout.Dimensions {
	paint.FillShape(
		gtx.Ops,
		col,
		clip.UniformRRect(
			image.Rectangle{
				Max: size,
			},
			0,
		).Op(gtx.Ops))
	return layout.Dimensions{Size: size}
}

// PaintRectMax paints a rectangle filling the max constraints
func PaintRectMax(gtx C, col color.NRGBA) D {
	size := image.Point{
		X: gtx.Constraints.Max.X,
		Y: gtx.Constraints.Max.Y,
	}

	return PaintRect(gtx, size, col)
}
