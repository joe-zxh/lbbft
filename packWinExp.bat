@echo off

echo 创建windows版本的实验文件...
cd ..

echo %cd%

del bftWin.zip
rd bftWin /s /q

md bftWin
md bftWin\keys
md bftWin\scripts

echo 拷贝Sublime包
copy .\Sublime.zip .\bftWin\

echo 拷贝keys
xcopy .\hotstuff\keys\* .\bftWin\keys\ /s

echo 编译hotstuff...
cd .\hotstuff
del .\hotstuffserver.exe 
del .\hotstuffclient.exe 
go build -o .\hotstuffserver.exe .\cmd\hotstuffserver\main.go
go build -o .\hotstuffclient.exe .\cmd\hotstuffclient\main.go
cd ..

echo 编译pbft...
cd .\pbft
del .\pbftserver.exe 
del .\pbftclient.exe 
go build -o .\pbftserver.exe .\cmd\pbftserver\main.go
go build -o .\pbftclient.exe .\cmd\pbftclient\main.go
cd ..

echo 编译pbftlinear...
cd .\pbftlinear
del .\pbftlinearserver.exe 
del .\pbftlinearclient.exe 
go build -o .\pbftlinearserver.exe .\cmd\pbftlinearserver\main.go
go build -o .\pbftlinearclient.exe .\cmd\pbftlinearclient\main.go
cd ..

echo 编译lbbft...
cd .\lbbft
del .\lbbftserver.exe 
del .\lbbftclient.exe 
go build -o .\lbbftserver.exe .\cmd\lbbftserver\main.go
go build -o .\lbbftclient.exe .\cmd\lbbftclient\main.go
cd ..

echo 拷贝hotstuff...
copy .\hotstuff\scripts\run_hs_server.bat .\bftWin\scripts
copy .\hotstuff\hotstuffclient.exe .\bftWin
copy .\hotstuff\hotstuffserver.exe .\bftWin

echo 拷贝pbft...
copy .\pbft\scripts\run_pbft_server.bat .\bftWin\scripts
copy .\pbft\pbftclient.exe .\bftWin
copy .\pbft\pbftserver.exe .\bftWin

echo 拷贝pbftlinear...
copy .\pbftlinear\scripts\run_pbftlinear_server.bat .\bftWin\scripts
copy .\pbftlinear\pbftlinearclient.exe .\bftWin
copy .\pbftlinear\pbftlinearserver.exe .\bftWin

echo 拷贝lbbft...
copy .\lbbft\scripts\run_lbbft_server.bat .\bftWin\scripts
copy .\lbbft\lbbftclient.exe .\bftWin
copy .\lbbft\lbbftserver.exe .\bftWin

echo 打包成bftWin.zip...
Bandizip.exe  a -r -l:5 bftWin.zip bftWin\*

echo 打包完成...

