# serve

> Quickly serve a file through firewalld

## Usage

Run the program with the file name and optionally a port. When the program
exists, the firewall port will be closed. The default port is `8080`.

```
$ serve <file> [port]
Serving `<file>` on port [port].
```

## Details

The program will first check that `firewalld` is running, but continues without
changing firewall rules if it's not. Regardless of the status of the port prior
to serving, it will be closed upon exit. Only tested to run on Linux,
specifically Arch.
