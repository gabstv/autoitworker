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

Opt("TrayMenuMode", 3)
; Hide tray icon!
Opt("TrayIconHide", 1)


; tray items storage
Local $trayKeys[64] = [0]
Local $trayValues[64] = [0]

; watch PID (if not 0)
Local $watchprocess = ""
Local $watchprocesscount = 0

; configuration
Local $httpErrorCount = 0
Local $httpMaxErrors = 100

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
    If (@error) Then
        $httpErrorCount = $httpErrorCount + 1
        If $httpMaxErrors > 0 And $httpErrorCount > $httpMaxErrors Then
            Return 0
        EndIf
        Return 1
    EndIf
    $httpErrorCount = 0
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
        Case "WinClose"
            Local $title = _JSONGet($object, "command.params.0")
            Local $text = _JSONGet($object, "command.params.1")
            Local $result = WinClose($title, $text)
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
        Case "_ping_"
            SendCommandResult(_JSONGet($object, "command.id"), "pong")
        Case "_set_config_"
            Local $cfgn = _JSONGet($object, "command.params.0")
            If $cfgn = "httpMaxErrors" Then
                $httpMaxErrors = _JSONGet($object, "command.params.1")
            ElseIf $cfgn = "watchprocess" Then
                $watchprocess = _JSONGet($object, "command.params.1")
            EndIf
            SendCommandResult(_JSONGet($object, "command.id"), "ok")
        Case Else
            SendCommandResult(_JSONGet($object, "command.id"), "unknown command")
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
        $watchprocesscount = $watchprocesscount + 1
        If $watchprocesscount > 1000 Then
            $watchprocesscount = 0
            If $watchprocess <> "" Then
                If Not ProcessExists($watchprocess) Then
                    Exit 2
                EndIf
            EndIf
        EndIf
        $lastCode = MainLoop()
    WEnd
    Exit 0
EndFunc

Run1()