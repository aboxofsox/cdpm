$RegKey = 'Registry::HKEY_LOCAL_MACHINE\System\CurrentControlSet\Control\Session Manager\Environment'

$InstallDir = "$($env:ProgramsFiles)\\cdpm"
$Binary = 'https://github.com/aboxofsox/cdpm/releases/download/v1.0.0/cdpm.exe'

wget $Binary -Outfile "$InstallDir\\cdpm.exe"

$CurrentEnvPath = (Get-ItemProperty -Path $RegKey -Name PATH).path
$NewEnvPath = "$CurrentEnvPath;$InstallDir"

Set-ItemProperty -Path $RegKey -Name PATH -Value $NewEnvPath