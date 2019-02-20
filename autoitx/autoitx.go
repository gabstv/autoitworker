package autoitx

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	autoitdll         = windows.NewLazyDLL("AutoItX3.dll")
	au3ControlGetText = autoitdll.NewProc("AU3_ControlGetText")
	au3ControlSetText = autoitdll.NewProc("AU3_ControlSetText")
	au3WinExists      = autoitdll.NewProc("AU3_WinExists")
	//au3ControlGetText, _ = syscall.GetProcAddress(autoitdll, "AU3_ControlGetText")
	//getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

// AU3_API void WINAPI AU3_ControlGetText(LPCWSTR szTitle, LPCWSTR szText, LPCWSTR szControl, LPWSTR szControlText, int nBufSize);

// ControlGetText Retrieves text from a control.
//   title     - The title/hWnd/class of the window to get the title.
//               See Title special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text      - [optional] The text of the window to get the title.
//               Default is an empty string. See Text special definition:
//               https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
//   controlID - The control to interact with. See Controls:
//               https://www.autoitscript.com/autoit3/docs/intro/controls.htm
func ControlGetText(title, text, controlID string) string {
	buf := make([]uint16, 1024*16)
	ret, _, callErr := au3ControlGetText.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(controlID))),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(1024*16),
	)
	if callErr != nil && !IsErrSuccess(callErr) {
		abort("Call ControlGetText", callErr)
	}
	return syscall.UTF16ToString(buf)
}

// AU3_API int WINAPI AU3_ControlSetText(LPCWSTR szTitle, LPCWSTR szText, LPCWSTR szControl, LPCWSTR szControlText);

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
func ControlSetText(title, text, controlID, setText string) bool {
	ret, _, callErr := au3ControlSetText.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(controlID))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(setText))),
	)
	if callErr != nil && !IsErrSuccess(callErr) {
		abort("Call ControlSetText", callErr)
	}
	return ret == 0
}

// AU3_API int WINAPI AU3_WinExists(LPCWSTR szTitle, /*[in,defaultvalue("")]*/LPCWSTR szText);

func WinExists(title, text string) bool {
	ret, _, callErr := au3WinExists.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
	)
	if callErr != nil && !IsErrSuccess(callErr) {
		abort("Call WinExists", callErr)
	}
	return ret == 1
}
