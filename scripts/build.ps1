$Package = 'cdpm'

$Platforms = 'amd64/windows', '386/windows'

Write-Host 'Building CDPM ⏳'
ForEach ($Platform in $Platforms) {
    $PlatformSplit = $Platform.Split('/')

    $GOARCH = $PlatformSplit[0]
    $GOOS = $PlatformSplit[1]

    $Env:GOOS = $GOOS
    $Env:GOARCH = $GOARCH

    $OutName = "./bin/$Package-$GOOS-$GOARCH"

    If ($GOOS -eq 'windows') {
        $OutName += '.exe'
        cmd.exe /c "go build -o $OutName ./bin/$Package"
    }
}

Write-Host 'Build complete ✅' -ForegroundColor  Green