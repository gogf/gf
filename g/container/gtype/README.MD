benchmarks:
```shell
john@johnstation:~/Workspace/Go/GOPATH/src/gitee.com/johng/gf/g/container/gtype$ go test *.go -bench=".*"
goos: linux
goarch: amd64
BenchmarkInt_Set-8         	300000000	         5.65 ns/op
BenchmarkInt_Val-8         	2000000000	         0.35 ns/op
BenchmarkInt_Add-8         	300000000	         5.68 ns/op
BenchmarkInt32_Set-8       	300000000	         5.67 ns/op
BenchmarkInt32_Val-8       	2000000000	         0.36 ns/op
BenchmarkInt32_Add-8       	300000000	         5.68 ns/op
BenchmarkInt64_Set-8       	300000000	         5.66 ns/op
BenchmarkInt64_Val-8       	2000000000	         0.35 ns/op
BenchmarkInt64_Add-8       	300000000	         5.67 ns/op
BenchmarkUint_Set-8        	300000000	         5.68 ns/op
BenchmarkUint_Val-8        	2000000000	         0.34 ns/op
BenchmarkUint_Add-8        	300000000	         5.68 ns/op
BenchmarkUint32_Set-8      	300000000	         5.65 ns/op
BenchmarkUint32_Val-8      	2000000000	         0.34 ns/op
BenchmarkUint32_Add-8      	300000000	         5.66 ns/op
BenchmarkUint64_Set-8      	300000000	         5.68 ns/op
BenchmarkUint64_Val-8      	2000000000	         0.34 ns/op
BenchmarkUint64_Add-8      	300000000	         5.68 ns/op
BenchmarkBool_Set-8        	300000000	         5.67 ns/op
BenchmarkBool_Val-8        	2000000000	         0.34 ns/op
BenchmarkString_Set-8      	30000000	        48.5 ns/op
BenchmarkString_Val-8      	100000000	        12.3 ns/op
BenchmarkBytes_Set-8       	50000000	        34.9 ns/op
BenchmarkBytes_Val-8       	100000000	        12.5 ns/op
BenchmarkInterface_Set-8   	50000000	        34.2 ns/op
BenchmarkInterface_Val-8   	100000000	        12.2 ns/op
PASS
ok  	command-line-arguments	43.447s
```


