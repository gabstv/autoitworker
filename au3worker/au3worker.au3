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

Local $trayKeys[64] = [0]
Local $trayValues[64] = [0]

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

Func SendResult(ByRef $jsonobject)
    Local $resp = _JSONEncode( $jsonobject, 'JSON_pack', '', @LF, false)
    $resp = _JSONFixLineBreaks($resp)
    _HTTP_Post($basePath & "/sync", $resp)
EndFunc

Func SendCommandResult($id, $result)
    Local $jobj = _JSONObject( _
        'success', True, _
        'type', 'command', _
        'command_id', $id, _
        'value', $result _
    )
    SendResult($jobj)
EndFunc

Func ParseCommand(ByRef $object)
    Local $name = _JSONGet($object, "command.name")
    Switch ($name)
        Case "Opt"
            Local $option = _JSONGet($object, "command.params.0")
            Local $value = _JSONGet($object, "command.params.1")
            Local $result = Opt($option, $value)
            SendCommandResult(_JSONGet($object, "command.id"), $result)
        Case "TraySetIcon"
            Local $filename = _JSONGet($object, "command.params.0")
            Local $iconid = _JSONGet($object, "command.params.1")
            TraySetIcon($filename, $iconid)
            SendCommandResult(_JSONGet($object, "command.id"), True)
        Case "TrayTip"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $timeout = _JSONGet($object, "command.params.2")
            Local $option = _JSONGet($object, "command.params.3")
            TrayTip($title, $text, $timeout, $option)
            SendCommandResult(_JSONGet($object, "command.id"), True)
        Case "TraySetState"
            Local $flag = _JSONGet($object, "command.params.0")
            TraySetState($flag)
            SendCommandResult(_JSONGet($object, "command.id"), True)
        Case "TrayCreateItem"
            ; input.Text, *input.MenuID, *input.MenuEntry, input.MenuRadioItem, trayid
            Local $text = _JSONGet($object, "command.params.0")
            Local $menuid = _JSONGet($object, "command.params.1")
            Local $menuentry = _JSONGet($object, "command.params.2")
            Local $radioitem = _JSONGet($object, "command.params.3")
            Local $trayid = _JSONGet($object, "command.params.4")
            Local $result = TrayCreateItem($text, $menuid, $menuentry, $radioitem)
            If $result = 0 Then
                SendCommandResult(_JSONGet($object, "command.id"), 0)
            Else
                ; add to tray array
                $trayKeys[$trayKeys[0]+1] = $result
                $trayValues[$trayValues[0]+1] = $trayid
                $trayKeys[0] = $trayKeys[0]+1
                $trayValues[0] = $trayValues[0]+1
                SendCommandResult(_JSONGet($object, "command.id"), 1)
            EndIf
        Case "WinGetTitle"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $result = WinGetTitle($title, $text)
            SendCommandResult(_JSONGet($object, "command.id"), $result)
        Case "ControlGetText"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $controlID = _JSONGet($object, "command.params.2")
            Local $result = ControlGetText($title, $text, $controlID)
            SendCommandResult(_JSONGet($object, "command.id"), $result)
        Case "ControlSetText"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $controlID = _JSONGet($object, "command.params.2")
            Local $setText = _JSONGet($object, "command.params.3")
            Local $flag = _JSONGet($object, "command.params.4")
            Local $result = ControlSetText($title, $text, $controlID, $setText, $flag)
            SendCommandResult(_JSONGet($object, "command.id"), $result)
        Case 12 To 17
            true
        Case Else
            false
    EndSwitch
EndFunc

Func Run1()
    Local $lastCode = 1
    While $lastCode = 1
        Local $trayevt = TrayGetMsg()
        If $trayevt <> 0 Then
            For $i = 0 To $trayKeys[0] Step 1
                If $trayevt = $trayKeys[$i+1] Then
                    _HTTP_Get($basePath & "/tray/" & $trayValues[$i+1])
                EndIf
            Next
        EndIf
        Sleep(10)
        $lastCode = MainLoop()
    WEnd
    Exit 0
EndFunc

Run1()