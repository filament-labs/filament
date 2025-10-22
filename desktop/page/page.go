package page

import (
	"gioui.org/layout"
	"github.com/filament-labs/filament/desktop/component"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page represents a screen or route in the app.
//
// Pages implement methods to identify themselves, provide optional navigation items,
// handle events, and render their layout.
//
// The NavItem method allows pages to optionally participate in the navigation drawer:
//   - If a page has a nav drawer entry, NavItem should return a pointer to its NavItem
//     struct containing the label, ID, icon, etc.
//   - If a page does not have a nav drawer entry, NavItem should return nil.
//
// This design lets the app register all pages for routing and event handling, but only
// adds pages with non-nil NavItem values to the navigation drawer. This keeps navigation
// consistent while allowing some pages (like splash or back pages) to be invisible in
// the drawer.
type Page interface {
	// ID returns a unique identifier for the page.
	ID() string

	// NavItem returns a navigation drawer entry for this page, or nil if the page
	// should not appear in the drawer.
	NavItem() (navItem *component.NavItem, isNavLayout bool)

	// HandleEvents processes any events relevant to the page.
	HandleEvents()

	// Layout renders the page contents and returns its layout dimensions.
	Layout(C) D
}

// Common page IDs
const (
	BackPageID  = "back" // generic ID for a back navigation page
	StartPageID = ""     // ID for the starting page of the app
)
