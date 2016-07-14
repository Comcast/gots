# goTS

gots (Go Transport Streams) is a library for working with MPEG transport streams. It provides abstractions for reading packet information and program specific information (psi)

## Bug / Feature Reporting
Add requests to github issues. To submit a PR see [CONTRIBUTING](./CONTRIBUTING)
## Tests
The included Makefile has a single target called `test`
```bash
make test
```
## License 
See [LICENSE](./LICENSE)
## Examples
This is a simple example that extracts all PIDs from a ts file and prints them. [CLI example parser can be found here](cli/parsefile.go)
```go
func main() {
	pidSet := make(map[uint16]bool, 5)
	filename := "./scenario1.ts"
	file, err := os.Open(filename)
	if err == nil {
		pkt := make([]byte, packet.PacketSize)
		for read, err := file.Read(pkt); read > 0 && err == nil; read, err = file.Read(pkt) {
                        if err != nil {
                                println(err)
                                return
                        }
			pid, err := packet.Pid(pkt)
			if err != nil {
				println(err)
				continue
			}
			pidSet[pid] = true
		}

        	for v := range pidSet {
	        	fmt.Printf("Found pid %d\n", v)
	        }
	} else {
		fmt.Printf("Unable to open file [%s] due to [%s]\n", filename, err.Error())
}
```
