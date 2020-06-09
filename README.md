# Source `RCON` Protocol

>This project is a work in progress.

The Source `RCON` Protocol is a `TCP/IP` based communication protocol used by [Source Dedicated Server](https://developer.valvesoftware.com/wiki/Source_Dedicated_Server), which allows console commands to be issued to the server via a remote console.

### Reference

* [`RCON` protocol](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#See_also)

## How to use it?

If you want to use this protocol implementation in your own projects, you only need to execute the next command in your project root directory:

```bash
$ go get -u github.com/csgoservers/source-rcon/pkg/protocol
```

Also, if you only want to execute some commands over a *Source Dedicated Servers*, then you can execute the `rcon-cli` application. To use it you just need to clone this repository and execute the `make build` directive. You can change some flags to configure your server settings. See the table below:

| Name 	| Default value 	| Description                                           	|
|------	|---------------	|-------------------------------------------------------	|
| `-H` 	| `127.0.0.1`   	| Host where the server is running                      	|
| `-p` 	| `27015`       	| Port where server is listening for connections. *TCP* 	|
| `-s` 	|               	| Password of the server if any                         	|
| `-c` 	|               	| Command to be executed                                	|

>Note that the `-c` flag is required in order to run the *cli*.

For example, if you want to execute the `cvarlist` command then you will use the following command:

```bash
$ ./rcon-cli -H 1.2.3.4 -p 27025 -s 1234 -c cvarlist
```

## Test

If you want to execute all tests from this repository then execute `make test`.
