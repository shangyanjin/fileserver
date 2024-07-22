@echo off
REM Set Go environment variables
set GO111MODULE=on

REM Check input parameters
if "%1"=="" goto menu
if "%2"=="" (
    echo Please provide an output file name.
    echo Usage: build.bat [win|linux] [output_filename]
    goto end
)

set TARGET=%1
set OUTPUT=%2

goto %TARGET%

:menu
echo Please choose the build target:
echo 1. Windows
echo 2. Linux
set /p choice=Enter your choice (1 or 2):

if "%choice%"=="1" set TARGET=windows
if "%choice%"=="2" set TARGET=linux
if not defined TARGET (
    echo Invalid choice, please run the script again and select 1 or 2.
    goto end
)

set /p OUTPUT=Enter the output file name:

:%TARGET%
if "%TARGET%"=="windows" goto windows
if "%TARGET%"=="linux" goto linux

:windows
echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o %OUTPUT%.exe main.go
if %errorlevel% neq 0 (
    echo Windows build failed!
    exit /b 1
)
echo Windows build succeeded.
goto end

:linux
echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-s -w" -o %OUTPUT% main.go
if %errorlevel% neq 0 (
    echo Linux build failed!
    exit /b 1
)
echo Linux build succeeded.
goto end

:end
echo Build process completed.
pause
