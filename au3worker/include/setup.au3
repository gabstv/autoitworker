#include <FileConstants.au3>

Func _SetupDefaults()
    AutoItSetOption("MouseCoordMode", 2)
    AutoItSetOption("SendKeyDelay", 15)
    AutoItSetOption("SendKeyDownDelay", 15)
EndFunc

Func _TruncateWrite($path, $data)
    FileDelete($path)
    Local $file = FileOpen($path, $FO_APPEND + $FO_CREATEPATH)
    FileWrite($file, $data)
    FileClose($file)
EndFunc

Func _ReadAllFile($path)
    Local $t0 = FileOpen($bonusFilePath, $FO_READ)
    Local $result = FileRead($t0)
    FileClose($t0)
    Return $result
EndFunc

Global $logging = false

Func _Log($message)
    If $logging Then ; global variable where you can globally switch OFF/ON logging
     FileWriteLine(@ScriptDir & '\file.log', @YEAR & "-" & @MON & "-" & @MDAY & " --> " & $message)
    EndIf
EndFunc

Func _JSONFixLineBreaks($input)
    Return StringReplace($input, @CRLF, '\r\n')
EndFunc