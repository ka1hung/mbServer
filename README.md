# Simple modbus server

## Quickly start

### install

    go get -u github.com/ka1hung/mbserver 

***

### simple example

    package main

    import "github.com/ka1hung/mbserver"

    func main() {
        mbs := mbserver.NewServer(uint8(1))
        mbs.Start("0.0.0.0:502")
    }

enjoy it :)
