@echo off
REM Set Go environment variables
set GO111MODULE=on

REM Check input parameters
if "%1"=="" (
    echo Please provide an output file name.
    echo Usage: linux.bat [output_filename]
    goto end
)

set OUTPUT=%1

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

:end
echo Build process completed.
