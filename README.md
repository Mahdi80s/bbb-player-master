Download public BigBlueButton recordings and play them offline.
Playback is to a great extent based on the BigBlueButton frontend called [bbb-playback](https://github.com/bigbluebutton/bbb-playback).

## Quickstart

Must have **Python3.6** or later (with **pip**)

download [python3](https://www.python.org/downloads/) here.


### in Windows platforms


run `setup.bat`

After setup completed, press Enter and exit setup

run `run.bat` and open [localhost:5000](http://localhost:5000).

## Go server (experimental)

An experimental Go implementation with PostgreSQL based login is located in
`go-server`. It requires Go 1.20+.

### Setup

1. Create a PostgreSQL database and apply `go-server/schema.sql`.
2. Set the environment variable `DATABASE_URL` or adjust the default connection
   string in `go-server/main.go`.
3. Start the server with:

```bash
go run ./go-server
```

The server will listen on `http://localhost:8080` and exposes simple register,
login and admin pages.
