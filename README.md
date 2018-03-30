### pushfile

### Introduction

Pushfile is a tool to synchronize one file to another server by ssh, you just need provide ssh config && localFile && remoteDir. It also can be synchronized in real time.

### Usage

```
go build -o pushfile
```

```
Usage of pushfile:
  -d string
    	remote dir (default "/tmp")
  -f string
    	local file path
  -h string
    	remoteHost:port (default "127.0.0.1:22")
  -u string
    	user:password (default "root:123456")
```

### LICENSE

MIT
