package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/filament-labs/filament/pkg/theme"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type NavItem struct {
	ID   string
	Text string
	Icon *widget.Icon
}

// renderNavItem holds both basic nav item state and the interaction
// state for that item.
type renderNavItem struct {
	NavItem
	hovering bool
	selected bool
	widget.Clickable
}

func (n *renderNavItem) Clicked(gtx C) bool {
	return n.Clickable.Clicked(gtx)
}

func (n *renderNavItem) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	for {
		event, ok := gtx.Event(pointer.Filter{
			Target: n,
			Kinds:  pointer.Enter | pointer.Leave,
		})
		if !ok {
			break
		}
		switch event := event.(type) {
		case pointer.Event:
			switch event.Kind {
			case pointer.Enter:
				n.hovering = true
			case pointer.Leave, pointer.Cancel:
				n.hovering = false
			}
		}
	}
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{
		Max: gtx.Constraints.Max,
	}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, n)
	return layout.Inset{
		Top:    unit.Dp(2),
		Bottom: unit.Dp(2),
		Left:   unit.Dp(4),
		Right:  unit.Dp(4),
	}.Layout(gtx, func(gtx C) D {
		return material.Clickable(gtx, &n.Clickable, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D { return n.layoutBackground(gtx, th) }),
				layout.Stacked(func(gtx C) D { return n.layoutContent(gtx, th) }),
			)
		})
	})
}

func (n *renderNavItem) layoutContent(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	contentColor := th.NavText
	if n.selected {
		//contentColor = th.Base.ContrastBg
	}
	return layout.Inset{
		Left:  unit.Dp(8),
		Right: unit.Dp(8),
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if n.NavItem.Icon == nil {
					return layout.Dimensions{}
				}
				return layout.Inset{Right: unit.Dp(40)}.Layout(gtx,
					func(gtx C) D {
						iconSize := gtx.Dp(unit.Dp(24))
						gtx.Constraints = layout.Exact(image.Pt(iconSize, iconSize))
						return n.NavItem.Icon.Layout(gtx, contentColor)
					})
			}),
			layout.Rigid(func(gtx C) D {
				label := material.Label(th.Base, unit.Sp(14), n.Text)
				label.Color = contentColor
				label.Font.Weight = font.Bold
				return layout.Center.Layout(gtx, label.Layout)
			}),
		)
	})
}

func (n *renderNavItem) layoutBackground(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if !n.selected && !n.hovering {
		return layout.Dimensions{}
	}

	var fill color.NRGBA
	if n.hovering || n.selected {
		fill = th.Background
	}

	rr := gtx.Dp(unit.Dp(2))
	defer clip.RRect{
		Rect: image.Rectangle{
			Max: gtx.Constraints.Max,
		},
		NE: rr,
		SE: rr,
		NW: rr,
		SW: rr,
	}.Push(gtx.Ops).Pop()
	theme.PaintRect(gtx, gtx.Constraints.Max, fill)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// NavDrawer implements the Material Design Navigation Drawer
// described here: https://material.io/components/navigation-drawer
type NavDrawer struct {
	selectedItem    int
	selectedChanged bool // selected item changed during the last frame
	items           []renderNavItem
	navigate        func(string)

	navList layout.List
}

// NewNav configures a navigation drawer
func NewNavDrawer(navigatorFunc func(pageID string)) NavDrawer {
	return NavDrawer{
		navigate: navigatorFunc,
	}
}

// AddNavItem inserts a navigation target into the drawer. This should be
// invoked only from the layout thread to avoid nasty race conditions.
func (m *NavDrawer) AddNavItem(item NavItem, selected bool) {
	m.items = append(m.items, renderNavItem{
		NavItem:  item,
		selected: selected,
	})
	if len(m.items) == 1 {
		m.items[0].selected = true
	}
}

func (m *NavDrawer) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	m.selectedChanged = false
	gtx.Constraints.Min.Y = 0
	m.navList.Axis = layout.Vertical
	return m.navList.Layout(gtx, len(m.items), func(gtx C, index int) D {
		gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(48))
		gtx.Constraints.Min = gtx.Constraints.Max
		if m.items[index].Clicked(gtx) {
			m.changeSelected(index)
		}
		dimensions := m.items[index].Layout(gtx, th)
		return dimensions
	})
}

func (m *NavDrawer) changeSelected(newIndex int) {
	if newIndex == m.selectedItem && m.items[m.selectedItem].selected {
		return
	}

	m.items[m.selectedItem].selected = false
	m.selectedItem = newIndex
	m.items[m.selectedItem].selected = true
	m.navigate(m.items[m.selectedItem].ID)
}
