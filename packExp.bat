
@echo off
echo %1%

cd ..

if "%1%"=="Linux" (
    set exeSuffix=
    set dirName=bftLinux
    
    SET GOOS=linux
    SET GOARCH=amd64

    echo ����Linux�汾��ʵ���ļ�...

) else if "%1%"=="Win" (
    set exeSuffix=.exe
    set dirName=bftWin

    SET GOOS=windows
    SET GOARCH=amd64

    echo ����windows�汾��ʵ���ļ�...  

) else (
    echo ����linux�汾��ʵ���ļ���������: .\packExp.bat Linux
    echo ����Windows�汾��ʵ���ļ���������: .\packExp.bat Win
    exit
)

del %dirName%.zip
rd %dirName% /s /q

md %dirName%
md %dirName%\keys

if "%1%"=="Linux" (
    echo ����Linux�汾�ű�...
    copy .\lbbft\scripts\linux\*.sh .\%dirName%
    copy .\lbbft\cmrecorder .\%dirName%
) else if "%1%"=="Win" (
    md %dirName%\scripts
    echo ����windows�汾�ű�...
    copy .\hotstuff\scripts\run_hotstuff_server.bat .\%dirName%\scripts
    copy .\pbft\scripts\run_pbft_server.bat .\%dirName%\scripts
    copy .\pbftlinear\scripts\run_pbftlinear_server.bat .\%dirName%\scripts
    copy .\lbbft\scripts\run_lbbft_server.bat .\%dirName%\scripts
)

echo ����keys
xcopy .\hotstuff\keys\* .\%dirName%\keys\ /s


echo ����hotstuff...
cd .\hotstuff
del .\hotstuffserver%exeSuffix% 
del .\hotstuffclient%exeSuffix% 
go build -o .\hotstuffserver%exeSuffix% .\cmd\hotstuffserver\main.go
go build -o .\hotstuffclient%exeSuffix% .\cmd\hotstuffclient\main.go
cd ..

echo ����hotstuff...
copy .\hotstuff\hotstuffserver%exeSuffix% .\%dirName%
copy .\hotstuff\hotstuffclient%exeSuffix% .\%dirName%
copy .\hotstuff\hotstuff.toml .\%dirName%


echo ����pbft...
cd .\pbft
del .\pbftserver%exeSuffix% 
del .\pbftclient%exeSuffix% 
del .\pbftvcclient%exeSuffix% 
go build -o .\pbftserver%exeSuffix% .\cmd\pbftserver\main.go
go build -o .\pbftclient%exeSuffix% .\cmd\pbftclient\main.go
go build -o .\pbftvcclient%exeSuffix% .\cmd\pbftvcclient\main.go
cd ..

echo ����pbft...
copy .\pbft\pbftserver%exeSuffix% .\%dirName%
copy .\pbft\pbftclient%exeSuffix% .\%dirName%
copy .\pbft\pbftvcclient%exeSuffix% .\%dirName%
copy .\pbft\pbft.toml .\%dirName%


echo ����pbftlinear...
cd .\pbftlinear
del .\pbftlinearserver%exeSuffix% 
del .\pbftlinearclient%exeSuffix% 
go build -o .\pbftlinearserver%exeSuffix% .\cmd\pbftlinearserver\main.go
go build -o .\pbftlinearclient%exeSuffix% .\cmd\pbftlinearclient\main.go
cd ..

echo ����pbftlinear...
copy .\pbftlinear\pbftlinearserver%exeSuffix% .\%dirName%
copy .\pbftlinear\pbftlinearclient%exeSuffix% .\%dirName%
copy .\pbftlinear\pbftlinear.toml .\%dirName%


echo ����lbbft...
cd .\lbbft
del .\lbbftserver%exeSuffix% 
del .\lbbftclient%exeSuffix% 
go build -o .\lbbftserver%exeSuffix% .\cmd\lbbftserver\main.go
go build -o .\lbbftclient%exeSuffix% .\cmd\lbbftclient\main.go
cd ..

echo ����lbbft...
copy .\lbbft\lbbftserver%exeSuffix% .\%dirName%
copy .\lbbft\lbbftclient%exeSuffix% .\%dirName%
copy .\lbbft\lbbft.toml .\%dirName%


echo %dirName%.zip...
Bandizip.exe  a -r -l:5 %dirName%.zip %dirName%\*

echo ������...

