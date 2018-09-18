package au3master

import (
	"encoding/json"
	"fmt"
	"time"
)

func (s *Server) sendcommand(cmd *Command) {
	select {
	case s.tosend <- cmd:
	default:
		// channel is full
		if s.Stdout != nil {
			s.Stdout.Write([]byte("[autoit master] channel is full!\r\n"))
		}
		for i := 0; i < 3; i++ {
			select {
			case droppedcmd := <-s.tosend:
				if s.Stdout != nil {
					s.Stdout.Write([]byte(fmt.Sprintf("dropped command {ID:%v Name:%v} \r\n", droppedcmd.ID, droppedcmd.Name)))
				}
			default:
			}
		}
		s.tosend <- cmd
	}
}

// Ping will send a ping command to test autoit connection
func (s *Server) Ping(timeout time.Duration) (time.Duration, error) {
	t0 := time.Now()
	cmd := newCommand("_ping_")
	s.sendcommand(cmd)
	rcvchan := s.waitchan(cmd.ID)
	select {
	case <-rcvchan:
		return time.Now().Sub(t0), nil
	case <-time.After(timeout):
		//return 0, fmt.Errorf("timeout")
	}
	return 0, fmt.Errorf("timeout")
}

// AutoItSetOption - Changes the operation of various AutoIt functions/parameters.
//
// https://www.autoitscript.com/autoit3/docs/functions/AutoItSetOption.htm
func (s *Server) AutoItSetOption(option string, param interface{}) []byte {
	cmd := newCommand("Opt")
	cmd.SetParams(option, param)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	return result.Value
}

// TraySetIcon Loads/Sets a specified tray icon.
//
// https://www.autoitscript.com/autoit3/docs/functions/TraySetIcon.htm
func (s *Server) TraySetIcon(filename string, iconID int) {
	cmd := newCommand("TraySetIcon")
	cmd.SetParams(filename, iconID)
	s.sendcommand(cmd)
	s.wait(cmd.ID)
}

// TrayTip Displays a balloon tip from the AutoIt Icon.
//
// title: Text appears in bold at the top of the balloon tip. (63 characters maximum)
// test: Message the balloon tip will display. (255 characters maximum)
// timeout: A rough estimate of the time (in seconds) the balloon tip should be displayed.
// (Windows has a min and max of about 10-30 seconds but does not always honor a time in that range.)
// option: [optional]
//    (0) = No icon (default)
//    (1) = Info icon
//    (2) = Warning icon
//    (3) = Error icon
//    (16) = Disable sound
func (s *Server) TrayTip(title, text string, timeout, option int) {
	cmd := newCommand("TrayTip")
	cmd.SetParams(title, text, timeout, option)
	s.sendcommand(cmd)
	s.wait(cmd.ID)
}

// TraySetState Sets the state of the tray icon.
//
// [optional] A combination of the following:
//     $TRAY_ICONSTATE_SHOW (1) = Shows the tray icon (default)
//     $TRAY_ICONSTATE_HIDE (2) = Destroys/Hides the tray icon
//     $TRAY_ICONSTATE_FLASH (4) = Flashes the tray icon
//     $TRAY_ICONSTATE_STOPFLASH (8) = Stops tray icon flashing
//     $TRAY_ICONSTATE_RESET (16) = Resets the icon to the defaults (no flashing, default tip text)
//
// Constants are defined in consts/tray.go
//
// https://www.autoitscript.com/autoit3/docs/functions/TraySetState.htm
func (s *Server) TraySetState(flag int) {
	cmd := newCommand("TraySetState")
	cmd.SetParams(flag)
	s.sendcommand(cmd)
	s.wait(cmd.ID)
}

// TrayCreateItemInput is the param set of TrayCreateItem
type TrayCreateItemInput struct {
	// The text of the control.
	Text string
	// [optional] Allows you to create a submenu in the referenced menu. If equal -1 it will be added 'behind' the last created item (default setting).
	MenuID *int
	// [optional] Allows you to define the entry number to be created. The entries are numbered starting at 0. If equal -1 it will be added 'behind' the last created entry (default setting).
	MenuEntry *int
	// [optional] (0) = (default) create a normal menuitem. (1) = create a menuradioitem.
	MenuRadioItem int
}

// TrayCreateItem Creates a menuitem control for the tray.
//
// text          ->	The text of the control.
// menuID        ->	[optional] Allows you to create a submenu in the referenced menu. If equal -1 it will be added 'behind' the last created item (default setting).
// menuentry     ->	[optional] Allows you to define the entry number to be created. The entries are numbered starting at 0. If equal -1 it will be added 'behind' the last created entry (default setting).
// menuradioitem ->	[optional]
//     $TRAY_ITEM_NORMAL (0) = (default) create a normal menuitem.
//     $TRAY_ITEM_RADIO (1) = create a menuradioitem.
//
// Constants are defined in TrayConstants.au3.
func (s *Server) TrayCreateItem(input TrayCreateItemInput) (<-chan bool, error) {
	if input.MenuEntry == nil {
		input.MenuEntry = Int(-1)
	}
	if input.MenuID == nil {
		input.MenuID = Int(-1)
	}
	s.trayindex++
	trayid := s.trayindex
	s.traychan[trayid] = make(chan bool, 2)
	cmd := newCommand("TrayCreateItem")
	cmd.SetParams(input.Text, *input.MenuID, *input.MenuEntry, input.MenuRadioItem, trayid)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := 0
	json.Unmarshal(result.Value, &rr)
	if rr == 0 {
		return nil, fmt.Errorf("failed to create tray item")
	}
	return s.traychan[trayid], nil
}

// WinGetTitle Retrieves the full title from a window.
//   title - The title/hWnd/class of the window to get the title.
//           See Title special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text  - [optional] The text of the window to get the title.
//           Default is an empty string. See Text special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
func (s *Server) WinGetTitle(title, text string) string {
	cmd := newCommand("WinGetTitle")
	cmd.SetParams(title, text)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := ""
	json.Unmarshal(result.Value, &rr)
	return rr
}

// WinClose Closes a window.
//   title - The title/hWnd/class of the window to get the title.
//           See Title special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text  - [optional] The text of the window to get the title.
//           Default is an empty string. See Text special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//
// https://www.autoitscript.com/autoit3/docs/functions/WinClose.htm
func (s *Server) WinClose(title, text string) int {
	cmd := newCommand("WinClose")
	cmd.SetParams(title, text)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	var rr int
	json.Unmarshal(result.Value, &rr)
	return rr
}

// ControlGetText Retrieves text from a control.
//   title     - The title/hWnd/class of the window to get the title.
//               See Title special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text      - [optional] The text of the window to get the title.
//               Default is an empty string. See Text special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//   controlID - The control to interact with. See Controls:
//               https://www.autoitscript.com/autoit3/docs/intro/controls.htm
func (s *Server) ControlGetText(title, text, controlID string) string {
	cmd := newCommand("ControlGetText")
	cmd.SetParams(title, text, controlID)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := ""
	json.Unmarshal(result.Value, &rr)
	return rr
}

// ControlGetText2 Retrieves text from a control (ID).
//   title     - The title/hWnd/class of the window to get the title.
//               See Title special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text      - [optional] The text of the window to get the title.
//               Default is an empty string. See Text special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//   controlID - The control to interact with. See Controls:
//               https://www.autoitscript.com/autoit3/docs/intro/controls.htm
func (s *Server) ControlGetText2(title, text string, controlID int) string {
	cmd := newCommand("ControlGetText")
	cmd.SetParams(title, text, controlID)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := ""
	json.Unmarshal(result.Value, &rr)
	return rr
}

// ControlSetText Sets text of a control.
//   title     - The title/hWnd/class of the window to get the title.
//               See Title special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text      - [optional] The text of the window to get the title.
//               Default is an empty string. See Text special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//   controlID - The control to interact with. See Controls:
//               https://www.autoitscript.com/autoit3/docs/intro/controls.htm
//   setText   - The text to be applied to the control.
//   flag      - [optional] when different from 0 (default) will force the target window to be redrawn.
func (s *Server) ControlSetText(title, text, controlID, setText string, flag int) bool {
	cmd := newCommand("ControlSetText")
	cmd.SetParams(title, text, controlID, setText, flag)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := false
	json.Unmarshal(result.Value, &rr)
	return rr
}

// ControlSetText2 Sets text of a control.
//   title     - The title/hWnd/class of the window to get the title.
//               See Title special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text      - [optional] The text of the window to get the title.
//               Default is an empty string. See Text special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//   controlID - The control to interact with. See Controls:
//               https://www.autoitscript.com/autoit3/docs/intro/controls.htm
//   setText   - The text to be applied to the control.
//   flag      - [optional] when different from 0 (default) will force the target window to be redrawn.
func (s *Server) ControlSetText2(title, text string, controlID int, setText string, flag int) bool {
	cmd := newCommand("ControlSetText")
	cmd.SetParams(title, text, controlID, setText, flag)
	s.sendcommand(cmd)
	result := s.wait(cmd.ID)
	rr := false
	json.Unmarshal(result.Value, &rr)
	return rr
}
