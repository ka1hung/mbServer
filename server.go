package mbserver

import (
	"context"
	"io"
	"log"
	"net"
)

//CommMSG is data struct for communication inspect
type CommMSG struct {
	Source   string
	ID       uint8
	Fn       uint8
	Addr     uint16
	RegData  []uint16
	CoilData []bool
}

// Server is a Modbus slave with allocated memory for discrete inputs, coils, etc.
type Server struct {
	listener        net.Listener
	port            int
	Datas           []Databank
	done            chan bool
	chanCommInspect chan CommMSG
}

// NewServer creates a new Modbus server.
func NewServer(num uint8) *Server {
	s := &Server{}
	s.Datas = make([]Databank, num)
	s.done = make(chan bool)
	return s
}

// Close stops listening to TCP/IP port.
func (s *Server) Close() {
	s.done <- true
}

func (s *Server) f1(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	data, err := s.Datas[id].ReadCoil(start, lens)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	plens := (lens/8 + 3)
	if lens%8 != 0 {
		plens++
	}
	wbuf = append(wbuf, byte(plens>>8))
	wbuf = append(wbuf, byte(plens))
	//id func
	wbuf = append(wbuf, rbuf[6:8]...)
	//Byte Count
	ByteCount := (lens + 7) / 8
	wbuf = append(wbuf, byte(ByteCount))

	//append data
	buf := make([]byte, ByteCount)
	for i := 0; i < len(data); i++ {
		index := uint16(i / 8)
		if data[i] == true {
			shift := (uint(i) % 8)
			buf[index] |= byte(1 << shift)
		}
	}
	wbuf = append(wbuf, buf...)
	s.passCommMsg(source, id, fn, start, data, []uint16{})
	return wbuf
}

func (s *Server) f2(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	data, err := s.Datas[id].ReadCoilIn(start, lens)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	plens := (lens/8 + 3)
	if lens%8 != 0 {
		plens++
	}
	wbuf = append(wbuf, byte(plens>>8))
	wbuf = append(wbuf, byte(plens))
	//id func
	wbuf = append(wbuf, rbuf[6:8]...)
	//Byte Count
	ByteCount := (lens + 7) / 8
	wbuf = append(wbuf, byte(ByteCount))

	//append data
	buf := make([]byte, ByteCount)
	for i := 0; i < len(data); i++ {
		index := uint16(i / 8)
		if data[i] == true {
			shift := (uint(i) % 8)
			buf[index] |= byte(1 << shift)
		}
	}
	wbuf = append(wbuf, buf...)
	s.passCommMsg(source, id, fn, start, data, []uint16{})
	return wbuf
}

func (s *Server) f3(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	data, err := s.Datas[id].ReadReg(start, lens)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	plens := (lens*2 + 3)
	wbuf = append(wbuf, byte(plens>>8))
	wbuf = append(wbuf, byte(plens))
	//id func
	wbuf = append(wbuf, rbuf[6:8]...)
	//data lens
	wbuf = append(wbuf, byte(lens*2))
	//append data
	for i := 0; i < len(data); i++ {
		wbuf = append(wbuf, byte(data[i]>>8))
		wbuf = append(wbuf, byte(data[i]))
	}
	s.passCommMsg(source, id, fn, start, []bool{}, data)
	return wbuf
}

func (s *Server) f4(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	data, err := s.Datas[id].ReadRegIn(start, lens)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	plens := (lens*2 + 3)
	wbuf = append(wbuf, byte(plens>>8))
	wbuf = append(wbuf, byte(plens))
	//id func
	wbuf = append(wbuf, rbuf[6:8]...)
	//data lens
	wbuf = append(wbuf, byte(lens*2))
	//append data
	for i := 0; i < len(data); i++ {
		wbuf = append(wbuf, byte(data[i]>>8))
		wbuf = append(wbuf, byte(data[i]))
	}
	s.passCommMsg(source, id, fn, start, []bool{}, data)
	return wbuf
}

func (s *Server) f5(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	data := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]

	var coil bool
	if data == 0 {
		coil = false
	} else {
		coil = true
	}

	err := s.Datas[id].WriteCoil(start, []bool{coil})
	if err != nil {
		return nil
	}
	s.passCommMsg(source, id, fn, start, []bool{coil}, []uint16{})
	return rbuf
}

func (s *Server) f6(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	data := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	err := s.Datas[id].WriteReg(start, []uint16{data})
	if err != nil {
		return nil
	}

	s.passCommMsg(source, id, fn, start, []bool{}, []uint16{data})
	return rbuf
}

func (s *Server) f15(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	//data
	data := make([]bool, lens)
	for i := 0; i < len(data); i++ {
		B := rbuf[13+(i/8)]
		bit := byte(1) << byte(i%8)
		if B&bit == 0 {
			data[i] = false
		} else {
			data[i] = true
		}
	}
	err := s.Datas[id].WriteCoil(start, data)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	wbuf = append(wbuf, []byte{0, 6}...)
	//package
	wbuf = append(wbuf, rbuf[6:12]...)
	s.passCommMsg(source, id, fn, start, data, []uint16{})
	return wbuf
}

func (s *Server) f16(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	//data
	data := []uint16{}
	for i := 0; i < int(lens); i++ {
		value := uint16(rbuf[13+i*2])<<8 + uint16(rbuf[14+i*2])
		data = append(data, value)
	}
	err := s.Datas[id].WriteReg(start, data)
	if err != nil {
		return nil
	}
	//tid pid
	wbuf := rbuf[:4]
	//package lens
	wbuf = append(wbuf, []byte{0, 6}...)
	//package
	wbuf = append(wbuf, rbuf[6:12]...)
	s.passCommMsg(source, id, fn, start, []bool{}, data)
	return wbuf
}
func (s *Server) passCommMsg(source string, id, fn uint8, addr uint16, coils []bool, regs []uint16) {
	if s.chanCommInspect != nil || len(s.chanCommInspect) < cap(s.chanCommInspect) {
		s.chanCommInspect <- CommMSG{
			Source:   source,
			ID:       id,
			Fn:       fn,
			Addr:     addr,
			CoilData: coils,
			RegData:  regs,
		}
	}
}
func (s *Server) processor(conn net.Conn) {
	defer conn.Close()
	defer log.Println("disconnect from: " + conn.RemoteAddr().String())
	log.Println("connect from: " + conn.RemoteAddr().String())

	for {
		packet := make([]byte, 512)
		bytesRead, err := conn.Read(packet)
		if err != nil {
			if err != io.EOF {
				log.Printf("read error %v\n", err)
			}
			return
		}
		// Set the length of the packet to the number of read bytes.
		packet = packet[:bytesRead]
		source := conn.RemoteAddr().String()

		id := int(packet[6])
		if id > len(s.Datas) {
			return
		}
		funcCode := int(packet[7])
		wbuf := []byte{}
		switch funcCode {
		case 1:
			wbuf = s.f1(source, packet)
		case 2:
			wbuf = s.f2(source, packet)
		case 3:
			wbuf = s.f3(source, packet)
		case 4:
			wbuf = s.f4(source, packet)
		case 5:
			wbuf = s.f5(source, packet)
		case 6:
			wbuf = s.f6(source, packet)
		case 15:
			wbuf = s.f15(source, packet)
		case 16:
			wbuf = s.f16(source, packet)
		default:
			return
		}

		if len(wbuf) != 0 {
			_, err = conn.Write(wbuf)
			if err != nil {
				log.Printf(err.Error())
				return
			}
		}

	}
}

// UseCommInspect ...
func (s *Server) UseCommInspect(size int) {
	s.chanCommInspect = make(chan CommMSG, size)
	// return Cfd
}

// ListenCommInspect ...
func (s *Server) ListenCommInspect() CommMSG {
	if s.chanCommInspect == nil {
		panic("should init CommInspect before use it")
	}
	return <-s.chanCommInspect
}

// Start the Modbus server listening on "address:port".
func (s *Server) Start(addressPort string) {
	//set server listen
	listen, err := net.Listen("tcp", addressPort)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
	}
	s.listener = listen
	defer s.listener.Close()

	//start server listening
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("Unable to accept connections: %#v\n", err)
			}
			go s.processor(conn)
		}
	}()

	for {
		select {
		case <-s.done:
			break
		}
	}
}

// StartWithContext start the Modbus server listening on "address:port" with context.
func (s *Server) StartWithContext(ctx context.Context, addressPort string) {
	//set server listen
	listen, err := net.Listen("tcp", addressPort)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
	}
	s.listener = listen
	defer s.listener.Close()

	//start server listening
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("Unable to accept connections: %#v\n", err)
			}
			go s.processor(conn)
		}
	}()

	for {
		select {
		case <-s.done:
			return
		case <-ctx.Done():
			return
		}
	}
}
