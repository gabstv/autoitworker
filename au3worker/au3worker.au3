#NoTrayIcon
#include "include\setup.au3"
#include "include\HTTP.au3"
#include "include\JSON.au3"
#include "include\JSON_Translate.au3"
#include "include\JSONGET.au3"

; Required for the $TRAY_ICONSTATE_SHOW constant.
#include <TrayConstants.au3>
#include <MsgBoxConstants.au3>
#include <InetConstants.au3>
#include <File.au3>

; Hide tray icon!
Opt("TrayIconHide", 1)

Local $nparams = $CmdLine[0]

If $nparams <> 1 Then
    Exit 1
EndIf

Local $basePath = $CmdLine[1]

; perform a test
Local $result = _HTTP_Get($basePath & "/health_check")

If $result <> "1" Then
    Exit 1
EndIf

; connection works
; will endlessly wait for commands

Func MainLoop()
    Local $data = _HTTP_Get($basePath & "/sync")
    Switch ($data)
        Case "0"
            ; do nothing
        Case Else
            $object = _JSONDecode($data)
            Local $action = _JSONGet($object, "action")
            Switch ($action)
                Case "command"
                    ParseCommand($object)
                Case "shutdown"
                    Return 0
                Case Else
                    ; do nothing
            EndSwitch
    EndSwitch
    Return 1
EndFunc

Func ParseCommand(ByRef $object)
    Local $name = _JSONGet($object, "command.name")
    Switch ($name)
        Case "WinGetTitle"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $result = WinGetTitle($title, $text)
            Local $jobj = _JSONObject( _
                'success', True, _
                'type', 'command', _
                'command_id', _JSONGet($object, "command.id"), _
                'value', $result _
            )
            Local $resp = _JSONEncode( $jobj, 'JSON_pack', '', @LF, false)
            $resp = _JSONFixLineBreaks($resp)
            _HTTP_Post($basePath & "/sync", $resp)
        Case "ControlGetText"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $controlID = _JSONGet($object, "command.params.2")
            Local $result = ControlGetText($title, $text, $controlID)
            Local $jobj = _JSONObject( _
                'success', True, _
                'type', 'command', _
                'command_id', _JSONGet($object, "command.id"), _
                'value', $result _
            )
            Local $resp = _JSONEncode( $jobj, 'JSON_pack', '', @LF, false)
            $resp = _JSONFixLineBreaks($resp)
            _HTTP_Post($basePath & "/sync", $resp)
        Case 12 To 17
            true
        Case Else
            false
    EndSwitch
EndFunc

Func Run1()
    Local $lastCode = 1
    While $lastCode = 1
        Sleep(10)
        $lastCode = MainLoop()
    WEnd
    Exit 0
EndFunc

Run1()