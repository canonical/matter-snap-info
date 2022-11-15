# EdgeX Snap Info

Install:
```
go install github.com/canonical/edgex-snap-info
```

Run:
```
edgex-snap-info --help
```

Example:
![image](https://user-images.githubusercontent.com/11150423/201926961-0212e1d3-9228-4b50-91c2-e9ee9282afda.png)


By default, the application fetches the config file from the repository. 

Build and run from source:
```
go run . --conf=./config.json
```
