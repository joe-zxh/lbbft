@echo off

echo ����windows�汾��ʵ���ļ�...
cd ..

echo %cd%

del bftWin.zip
rd bftWin /s /q

md bftWin
md bftWin\keys
md bftWin\scripts

echo ����Sublime��
copy .\Sublime.zip .\bftWin\

echo ����keys
xcopy .\hotstuff\keys\* .\bftWin\keys\ /s

echo ����hotstuff...
cd .\hotstuff
del .\hotstuffserver.exe 
del .\hotstuffclient.exe 
go build -o .\hotstuffserver.exe .\cmd\hotstuffserver\main.go
go build -o .\hotstuffclient.exe .\cmd\hotstuffclient\main.go
cd ..

echo ����pbft...
cd .\pbft
del .\pbftserver.exe 
del .\pbftclient.exe 
go build -o .\pbftserver.exe .\cmd\pbftserver\main.go
go build -o .\pbftclient.exe .\cmd\pbftclient\main.go
cd ..

echo ����pbftlinear...
cd .\pbftlinear
del .\pbftlinearserver.exe 
del .\pbftlinearclient.exe 
go build -o .\pbftlinearserver.exe .\cmd\pbftlinearserver\main.go
go build -o .\pbftlinearclient.exe .\cmd\pbftlinearclient\main.go
cd ..

echo ����lbbft...
cd .\lbbft
del .\lbbftserver.exe 
del .\lbbftclient.exe 
go build -o .\lbbftserver.exe .\cmd\lbbftserver\main.go
go build -o .\lbbftclient.exe .\cmd\lbbftclient\main.go
cd ..

echo ����hotstuff...
copy .\hotstuff\scripts\run_hs_server.bat .\bftWin\scripts
copy .\hotstuff\hotstuffclient.exe .\bftWin
copy .\hotstuff\hotstuffserver.exe .\bftWin

echo ����pbft...
copy .\pbft\scripts\run_pbft_server.bat .\bftWin\scripts
copy .\pbft\pbftclient.exe .\bftWin
copy .\pbft\pbftserver.exe .\bftWin

echo ����pbftlinear...
copy .\pbftlinear\scripts\run_pbftlinear_server.bat .\bftWin\scripts
copy .\pbftlinear\pbftlinearclient.exe .\bftWin
copy .\pbftlinear\pbftlinearserver.exe .\bftWin

echo ����lbbft...
copy .\lbbft\scripts\run_lbbft_server.bat .\bftWin\scripts
copy .\lbbft\lbbftclient.exe .\bftWin
copy .\lbbft\lbbftserver.exe .\bftWin

echo �����bftWin.zip...
Bandizip.exe  a -r -l:5 bftWin.zip bftWin\*

echo ������...

