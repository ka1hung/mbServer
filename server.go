package mbserver

import (
	"context"
	"io"
	"log"
	"net"
)

// CommMSG is data struct for communication inspect
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
	host            string
	Datas           []Databank
	ChanCommInspect chan CommMSG
}

// NewServer creates a new Modbus server and numbers for modbus devices (1~254).
func NewServer(num int) *Server {
	if num > 254 || num < 1 {
		num = 1
	}
	s := &Server{}
	s.Datas = make([]Databank, num)
	return s
}

func (s *Server) passCommMsg(source string, id, fn uint8, addr uint16, coils []bool, regs []uint16) {
	if s.ChanCommInspect != nil || len(s.ChanCommInspect) < cap(s.ChanCommInspect) {
		s.ChanCommInspect <- CommMSG{
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
		packet := make([]byte, 65535)
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
		if id > len(s.Datas) || id < 0 {
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

		if len(wbuf) > 0 {
			_, err = conn.Write(wbuf)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}

	}
}
func (s *Server) processorWithContext(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defer log.Println("disconnect from: " + conn.RemoteAddr().String())
	log.Println("connect from: " + conn.RemoteAddr().String())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			packet := make([]byte, 65535)
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
			if id > len(s.Datas) || id < 0 {
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
					log.Println(err.Error())
					return
				}
			}
		}
	}
}

// UseCommInspect ...
func (s *Server) UseCommInspect(size int) {
	s.ChanCommInspect = make(chan CommMSG, size)
}

// ListenCommInspect ...
func (s *Server) ListenCommInspect() CommMSG {
	if s.ChanCommInspect == nil {
		panic("should init CommInspect before use it")
	}
	return <-s.ChanCommInspect
}

// Start the Modbus server listening on "address:port".
func (s *Server) Start(host string) {
	//set server listen
	s.host = host
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
	}
	s.listener = listen
	defer s.listener.Close()

	//start server listening
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connections: %#v\n", err)
		}
		go s.processor(conn)
	}
}

// Start the Modbus server listening on "address:port".
func (s *Server) StartWithContext(ctx context.Context, host string) {
	//set server listen
	s.host = host
	listen, err := net.Listen("tcp", host)
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
				if conn == nil {
					return
				}
			}
			go s.processorWithContext(ctx, conn)
		}
	}()

	for range ctx.Done() {
		return
	}
}
