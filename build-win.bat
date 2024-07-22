@echo off
REM Set Go environment variables
set GO111MODULE=on

REM Check input parameters
if "%1"=="" (
    echo Please provide an output file name.
    echo Usage: win.bat [output_filename]
    goto end
)

set OUTPUT=%1

echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o %OUTPUT%.exe main.go
if %errorlevel% neq 0 (
    echo Windows build failed!
    exit /b 1
)
echo Windows build succeeded.

:end
echo Build process completed.
