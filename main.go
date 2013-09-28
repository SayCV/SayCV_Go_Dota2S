// Copyright 2013 The GoDota2S Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main
package main

import (
	"log"
	"time"
	"strconv"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tarm/goserial"
)

var isSpecialMode = walk.NewMutableCondition()

type MyMainWindow struct {
	*walk.MainWindow
}

func main() {
	MustRegisterCondition("isSpecialMode", isSpecialMode)
	
	mw := new(MyMainWindow)

	var openAction, showAboutBoxAction *walk.Action
	var recentMenu *walk.Menu
	var toggleSpecialModePB *walk.PushButton

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "sayTerm",
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Image:       "img/open.png",
						Enabled:     Bind("enabledCB.Checked"),
						Visible:     Bind("openVisibleCB.Checked"),
						Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: mw.openAction_Triggered,
					},
					Menu{
						AssignTo: &recentMenu,
						Text:     "Recent",
					},
					Separator{},
					Action{
						Text:        "&Save",
						OnTriggered: mw.saveAction_Triggered,
					},
					Separator{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},
		ToolBarItems: []MenuItem{
			ActionRef{&openAction},
			Menu{
				Text:  "New A",
				Image: "img/document-new.png",
				Items: []MenuItem{
					Action{
						Text:        "A",
						OnTriggered: mw.newAction_Triggered,
					},
					Action{
						Text:        "B",
						OnTriggered: mw.newAction_Triggered,
					},
					Action{
						Text:        "C",
						OnTriggered: mw.newAction_Triggered,
					},
				},
				OnTriggered: mw.newAction_Triggered,
			},
			Separator{},
			Menu{
				Text:  "View",
				Image: "img/document-properties.png",
				Items: []MenuItem{
					Action{
						Text:        "X",
						OnTriggered: mw.changeViewAction_Triggered,
					},
					Action{
						Text:        "Y",
						OnTriggered: mw.changeViewAction_Triggered,
					},
					Action{
						Text:        "Z",
						OnTriggered: mw.changeViewAction_Triggered,
					},
				},
			},
			Separator{},
			Action{
				Text:        "Special",
				Image:       "img/system-shutdown.png",
				Enabled:     Bind("isSpecialMode && enabledCB.Checked"),
				OnTriggered: mw.specialAction_Triggered,
			},
		},
		ContextMenuItems: []MenuItem{
			ActionRef{&showAboutBoxAction},
		},
		MinSize:  Size{600, 400},
		Size:     Size{800, 600},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			CheckBox{
				Name:    "enabledCB",
				Text:    "Open / Special Enabled",
				Checked: true,
			},
			CheckBox{
				Name:    "openVisibleCB",
				Text:    "Open Visible",
				Checked: true,
			},
			PushButton{
				AssignTo: &toggleSpecialModePB,
				Text:     "Enable Special Mode",
				OnClicked: func() {
					isSpecialMode.SetSatisfied(!isSpecialMode.Satisfied())

					if isSpecialMode.Satisfied() {
						toggleSpecialModePB.SetText("Disable Special Mode")
					} else {
						toggleSpecialModePB.SetText("Enable Special Mode")
					}
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	
	// We load our icon from a file.
	icon, err := walk.NewIconFromFile("img/ni.ico")
	if err != nil {
		log.Fatal(err)
	}

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			"sayTerm",
			"sayTerm."); err != nil {

			log.Fatal(err)
		}
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := ni.ShowInfo("sayTerm", "Click the icon to show again."); err != nil {
		log.Fatal(err)
	}
	
	// 
	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	// lv.PostAppendText("")
	log.SetOutput(lv)
/*	
	ser_cfg := &serial.Config{Name: "COM4", Baud: 115200}
	
	ser, err := serial.OpenPort(ser_cfg)
	if err != nil {
		log.Fatal(err)
	}
	
	ser_buf := make([]byte, 1024)

	go func() {
		for i := 0; true; i++ {
			time.Sleep(100 * time.Millisecond)
			_, err := ser.Read(ser_buf)
			if err != nil {
				log.Fatal(err)
			}
			
			log.Print(strconv.Atoi(string(ser_buf)))
		}
	}()
*/
	mw.Run()
}

func (mw *MyMainWindow) openAction_Triggered() {
	//walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
    ser_cfg := &serial.Config{Name: "COM4", Baud: 115200}
	
	ser, err := serial.OpenPort(ser_cfg)
	if err != nil {
		//log.Fatal(err)
		walk.MsgBox(mw, "Open", "error", walk.MsgBoxIconInformation)
	}
	
	ser_buf := make([]byte, 1024)

	go func() {
		for i := 0; true; i++ {
			time.Sleep(100 * time.Millisecond)
			_, err := ser.Read(ser_buf)
			if err != nil {
				//log.Fatal(err)
			}
			
			log.Print(strconv.Atoi(string(ser_buf)))
		}
	}()
}

func (mw *MyMainWindow) saveAction_Triggered() {
	walk.MsgBox(mw, "Save", "Pretend to save a file...", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) newAction_Triggered() {
	walk.MsgBox(mw, "New", "Newing something up... or not.", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) changeViewAction_Triggered() {
	walk.MsgBox(mw, "Change View", "By now you may have guessed it. Nothing changed.", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About", "Walk Actions Example", walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) specialAction_Triggered() {
	walk.MsgBox(mw, "Special", "Nothing to see here.", walk.MsgBoxIconInformation)
}
