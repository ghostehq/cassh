//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

// Show app in dock (regular app)
void showInDock() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    });
}

// Hide from dock (accessory/menu bar only)
void hideFromDock() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}

// Check if currently showing in dock
int isShowingInDock() {
    return [NSApp activationPolicy] == NSApplicationActivationPolicyRegular ? 1 : 0;
}
*/
import "C"
import (
	"log"

	"github.com/getlantern/systray"
	"github.com/shawntz/cassh/internal/config"
)

var (
	menuShowInDock *systray.MenuItem
)

// setupVisibilityMenu adds the visibility toggle menu items
func setupVisibilityMenu() *systray.MenuItem {
	// Create "Appearance" submenu
	menuAppearance := systray.AddMenuItem("Appearance", "App visibility options")

	menuShowInDock = menuAppearance.AddSubMenuItemCheckbox("Show in Dock", "Show cassh icon in the Dock", cfg.User.ShowInDock)

	return menuAppearance
}

// handleShowInDockToggle toggles dock visibility
func handleShowInDockToggle() {
	if menuShowInDock.Checked() {
		// Currently checked, uncheck it (hide from dock)
		menuShowInDock.Uncheck()
		hideFromDock()
		cfg.User.ShowInDock = false
	} else {
		// Currently unchecked, check it (show in dock)
		menuShowInDock.Check()
		showInDock()
		cfg.User.ShowInDock = true
	}

	// Save preference
	if err := config.SaveUserConfig(&cfg.User); err != nil {
		log.Printf("Failed to save dock visibility preference: %v", err)
	}
}

// applyVisibilitySettings applies saved visibility settings on startup
func applyVisibilitySettings() {
	if cfg.User.ShowInDock {
		showInDock()
	} else {
		hideFromDock()
	}
}

// Helper functions exposed to Go
func showInDock() {
	C.showInDock()
}

func hideFromDock() {
	C.hideFromDock()
}

func isShowingInDock() bool {
	return C.isShowingInDock() == 1
}
