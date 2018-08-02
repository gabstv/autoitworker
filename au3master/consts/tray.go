package consts

//     $TRAY_ICONSTATE_SHOW (1) = Shows the tray icon (default)
//     $TRAY_ICONSTATE_HIDE (2) = Destroys/Hides the tray icon
//     $TRAY_ICONSTATE_FLASH (4) = Flashes the tray icon
//     $TRAY_ICONSTATE_STOPFLASH (8) = Stops tray icon flashing
//     $TRAY_ICONSTATE_RESET (16) = Resets the icon to the defaults (no flashing, default tip text)

const (
	// TrayIconstateShow - Shows the tray icon (default)
	TrayIconstateShow int = 1
	// TrayIconstateHide - Destroys/Hides the tray icon
	TrayIconstateHide int = 1
	// TrayIconstateFlash - Flashes the tray icon
	TrayIconstateFlash int = 4
	// TrayIconstateStopFlash - Stops tray icon flashing
	TrayIconstateStopFlash int = 8
	// TrayIconstateReset - Resets the icon to the defaults (no flashing, default tip text)
	TrayIconstateReset int = 16
)
