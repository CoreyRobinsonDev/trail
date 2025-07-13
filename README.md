# trail
Leave a trail of commands
```bash
c to add a comment • r to remove a command • t to set a session title • esc quit
0 ping 10.129.134.166
1 nmap -sV 10.129.134.166
  • 135/tcp  open  msrpc         Microsoft Windows RPC
    139/tcp  open  netbios-ssn   Microsoft Windows netbios-ssn
    445/tcp  open  microsoft-ds?
    3389/tcp open  ms-wbt-server Microsoft Terminal Services
    5985/tcp open  http          Microsoft HTTPAPI httpd 2.0 (SSDP/UPnP)

2 xfreerdp3 /v:10.129.134.166 /cert:ignore /u:Administrator
  • hit
  • flag: 951fa96d7830c451b536be5a6be008a0
```

<br>

[Usage](#Usage) <span>&nbsp;•&nbsp;</span> [Install](#Install)

## Usage
- Run `trail` with no flags to start a new session
- Use `-ls` to list past sessions
```bash
> trail -ls
3ATO4ATNwADMuAzK        2025-07-01 22:31:33.963043556 -0500 EST

run trail -conn {id} to continue the session
run trail -rm {id|"*"} to remove a session
run trail -x {id} to export a session
```

- Use `-conn {id}` to continue an existing session
```bash
trail -conn 3A
```

- Use `-x {id}` to export an existing session to a markdown file
```bash
trail -x 3A
```

- Use `-rm {id}` to remove an existing session
- Use `-rm "*"` to remove all past sessions
```bash
trail -rm 3A
```
## Install
Download pre-built binary for your system here [Releases](https://github.com/CoreyRobinsonDev/trail/releases).

### Compiling from Source
- Clone this repository
```bash
git clone https://github.com/CoreyRobinsonDev/trail.git
```
- Create **trail** binary
```bash
cd trail
sudo go build -o /usr/local/bin
```
- Confirm that the program was built successfully
```bash
trail -v
```
## License
[The Unlicense](./LICENSE)
