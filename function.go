package mbserver

func (s *Server) f1(source string, rbuf []byte) []byte {
	start := uint16(rbuf[8])<<8 | uint16(rbuf[9])
	lens := uint16(rbuf[10])<<8 | uint16(rbuf[11])
	id := rbuf[6] - 1
	fn := rbuf[7]
	index := id
	if id == 255 {
		index = 0
	}
	data, err := s.Datas[index].ReadCoil(start, lens)
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
	index := id
	if id == 255 {
		index = 0
	}
	data, err := s.Datas[index].ReadCoilIn(start, lens)
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
	index := id
	if id == 255 {
		index = 0
	}
	data, err := s.Datas[index].ReadReg(start, lens)
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
	index := id
	if id == 255 {
		index = 0
	}
	data, err := s.Datas[index].ReadRegIn(start, lens)
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

	index := id
	if id == 255 {
		index = 0
	}
	err := s.Datas[index].WriteCoil(start, []bool{coil})
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

	index := id
	if id == 255 {
		index = 0
	}
	err := s.Datas[index].WriteReg(start, []uint16{data})
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

	index := id
	if id == 255 {
		index = 0
	}
	err := s.Datas[index].WriteCoil(start, data)
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

	index := id
	if id == 255 {
		index = 0
	}
	err := s.Datas[index].WriteReg(start, data)
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
