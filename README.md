# Simple modbus server

Simple Modbus Server support modbus function codes: 
 + Read coils (FC1) 
 + Read input discretes (FC2) 
 + Read multiple registers (FC3) 
 + Read input registers (FC4) 
 + Write single coil (FC5) 
 + Write single register (FC6)
 + Write multiple coils (FC15)  
 + Write multiple registers (FC16)

## Quickly start

### install

    go get -u github.com/ka1hung/mbserver

***

### simple example
```go
package main

import "github.com/ka1hung/mbserver"

func main() {
    // num = 1 for use one modbus device(ID1)
    // range 1 ~ 254
    mbs := mbserver.NewServer(1)
    mbs.Start("0.0.0.0:502")
}
```

### Data handle
```go
package main

import (
    "fmt"

    "github.com/ka1hung/mbserver"
)

func main() {
    // num2 for use 2 modbus device(ID1 and ID2)
    mbs := mbserver.NewServer(2)

    // mbs.Datas[0] for handle Device ID1
    fmt.Println(mbs.Datas[0].WriteCoil(0, []bool{true,false,true}))
    fmt.Println(mbs.Datas[0].ReadCoil(0, 3))

    fmt.Println(mbs.Datas[0].WriteCoilIn(0, []bool{true,false,true}))
    fmt.Println(mbs.Datas[0].ReadCoilIn(0, 3))

    fmt.Println(mbs.Datas[0].WriteReg(0, []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
    fmt.Println(mbs.Datas[0].ReadReg(0, 10))

    fmt.Println(mbs.Datas[0].WriteRegIn(0, []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
    fmt.Println(mbs.Datas[0].ReadRegIn(0, 10))


    // mbs.Datas[1] for handle Device ID2
    fmt.Println(mbs.Datas[1].WriteReg(0, []uint16{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}))
    fmt.Println(mbs.Datas[1].ReadReg(0, 10))

    mbs.Start("0.0.0.0:502")
}
```

### listen communications
```go
package main

import (
    "fmt"
    
    "github.com/ka1hung/mbserver"
)

func main() {     
    mbs := mbserver.NewServer(1)
    mbs.UseCommInspect(100) //set buff 100
    go mbs.Start("0.0.0.0:502")
    for{
        msg:= mbs.ListenCommInspect()
        fmt.Printf("%+v\n",msg)
    }
}
```
enjoy it :)

### LICENSE
[MIT](https://github.com/ka1hung/mbserver/blob/master/LICENSE)
