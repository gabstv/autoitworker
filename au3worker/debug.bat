ECHO Au3Worker debug mode.
set /p connStr="Digite o URL do servidor (Go): "
AutoIt3.exe au3worker.au3 %connStr%
set /p temp="Aperte qualquer tecla para encerrar"