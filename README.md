# EdgeX Snap Info

Install:
```
go install github.com/canonical/edgex-snap-info
```

Run:
```
edgex-snap-info
```

By default, the application fetches the config file from the repository. 
Set `--help` flag for more details.

Build and run from source:
```
go run . --conf=./config.json
```
