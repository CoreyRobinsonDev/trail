# trail
Leave a trail of commands

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
```

- Use `-conn {id}` to continue an existing session
```bash
trail -conn 3A
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
go build -o /usr/local/bin
```
- Confirm that the program was built successfully
```bash
trail -v
```
## License
[The Unlicense](./LICENSE)
