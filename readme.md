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
  -d, --duration int       How long to listen for packets. (default 30)
  -h, --help               help for cdpm
  -i, --interface string   Defines the interface to listen to.
  -l, --log                Log received packets to file.
  -o, --outfile string     Output results to file. (default "results.md")
  -t, --toggle             Help message for toggle

Use "cdpm [command] --help" for more information about a command.
```


### Example Usage
Start CDPM and use the default interface, `Ethernet`.

*Note: depending on how the Cisco switch is configured, it could take some time to receive the first CDP packet.*
```
cdpm start
```
Start CDPM and define an interface to listen to.
```
cdpm start -i Ethernet
```
Start CDPM and log packets received for 30 seconds and write their info to file.
```ps1
cdpm start -l -d 30
```

### Available Data
While there's a lot more to a typical CDP or LLDP packet, the given purpose of this tool only needs the following:

- The device name
- Device IP
- Port
- Native VLAN
- Voice VLAN

For trunked ports, the switch may send an LLDP packet, which isn't *quite* as useful as a CDP switch, given the scope of this tool.

- Port description
- System name
- System description

### Managed vs Unmanaged
If you run CDPM on an unmanaged switch, you will be able to obtain the host name, but not an IP address. It is expected for the `Device IP` field to be empty when listening to a network interface connected to an unmanaged switch.

Data expansion isn't planned, but I do plan on refactoring parts of this tool as separate packages; to be more modular. 


