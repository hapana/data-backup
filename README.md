# Data-backup

Simple binary to manage data backups when you need to wrap commands around it.
The reason for building this is I have a number of Dockerised apps I run for
home and would like to backup data directories to my NAS but first need to make
sure the Docker app is down before doing the copy.

Therefore, I wanted something I could run as a cron and point to a directory
where I could read a config file and operate on paths relative to that
directory.

## Installing

You can run `go get github.com/hapana/data-backup` to install this binary

# Usage

The config for orchestrating your backups lives within a `.backup.json` file.
This file is read in automatically when you pass a path to the binary like so:

`data-backup -p ./test/`

# Config

The config model is as follows:

```
{
  "directories": {
    "data": "./data",
    "backup": "./backup"
  }
}
```

The data field contains the relative back to the directory you'd like to backup
and the backup field contains the directory you'd like to backup to"

# Contact

Please raise issues on this repo for any addiitional features or bugs that may exist.
