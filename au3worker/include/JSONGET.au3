#include <Array.au3>
#include "JSON.au3"
#include "JSON_Translate.au3" ; examples of translator functions, includes JSON_pack and JSON_unpack

; https://www.autoitscript.com/forum/topic/104150-json-udf-library-fully-rfc4627-compliant/?do=findComment&comment=1030327
Func _JSONGet($json, $path, $seperator = ".")
    Local $seperatorPos,$current,$next,$l
    $seperatorPos = StringInStr($path, $seperator)
    If $seperatorPos > 0 Then
        $current = StringLeft($path, $seperatorPos - 1)
        $next = StringTrimLeft($path, $seperatorPos + StringLen($seperator) - 1)
    Else
        $current = $path
        $next = ""
    EndIf
    If _JSONIsObject($json) Then
        $l = UBound($json, 1)
        For $i = 0 To $l - 1
            If $json[$i][0] == $current Then
            If $next == "" Then
            return $json[$i][1]
            Else
            return _JSONGet($json[$i][1], $next, $seperator)
            EndIf
            EndIf
        Next
    ElseIf IsArray($json) And UBound($json, 0) == 1 And UBound($json, 1) > $current Then
        If $next == "" Then
            return $json[$current]
        Else
            return _JSONGet($json[$current], $next, $seperator)
        EndIf
    EndIf

    return $_JSONNull
EndFunc