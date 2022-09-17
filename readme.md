# CDPM
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
) ![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)

<abbr title='Cisco Discovery Packet Monitor'>CDPM</abbr> is a simple **Windows** tool to get CDP packet data from a Cisco switch. 

***[Npcap](https://npcap.com/#download) is required.***

### Install
- Download and install [Npcap](https://npcap.com/#download).
  - If you've been using WireShark, you likely already have Npcap.
- Download the latest executable binary from [releases](https://github.com/aboxofsox/cdpm/releases).
- Move the executable binary to some folder, probably `C:\Program Files\cdpm`.
- Add that path to your `PATH` environment variable.

*or*

- Run the `install.ps1` script.
  - If you're on Windows 7 for some reason, or using some older version of PowerShell, the installation script might not work.


### Usage:
```
cdpm [command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
help        Help about any command
start       start cdpm

Flags:
-h, --help               help for cdpm
-i, --interface string   Defines the interface to listen to.
-o, --outfile string     Output results to file. (default "results.md")
-t, --toggle             Help message for toggle

Use "cdpm [command] --help" for more information about a command.
```


### Example Usage
Start CDPM and use the default interface, `Ethernet`.
```
cdpm start
```
Start CDPM and define an interface to listen to.
```
cdpm start -i Ethernet
```


