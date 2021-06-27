# Radio Simulator
> Originally created by WeiTing, modified by YanJie.
>
> In this version, only support registration /deregistration procedure & service request procedure (without any PDU session)

## Requirement
To use this simulator, you need to run a MongoDB and modify `configs/rancfg.yaml`
```yaml
dbName: simulator
dbUrl: mongodb://127.0.0.1:27017 # your mongo url
.....
```

And for CLI tool, you can use `--db` option to config the url
```bash
./bin/simctl --db <Mongo URL> [command]

# example
./bin/simctl --db mongodb://127.0.0.1:27017 get ues
```

For more information, please use `--help` option

## Build

```bash
# build simulator & simctl
make
```

## Run

### 1. Run RAN simulator
```bash
./bin/simulator
```
It will create a RAN and connect to AMF, remember to run AMF first before running this simulator.

### 2. Run CLI to connect to RAN simulator
After running the RAN simulator, please open another terminal to run the cli command to control the RAN simulator
```bash
./bin/simctl --help # print the help message
```

For more information about how to use the command, please read the documents in `docs/`

## Generate documents for simctl
> reference: https://github.com/spf13/cobra/blob/master/doc/md_docs.md

1. Uncomment `doc.GenMarkdownTree` function call in main function of `cmd/simctl/main.go`, and comment `rootCmd.Execute()` function call
```go
func main() {
	// Generate command document
	if err := doc.GenMarkdownTree(rootCmd, "./docs"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
```

2. In project root folder, run `go run cmd/simctl/main.go`
