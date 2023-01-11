package mbserver

import (
	"errors"
	"sync"
)

// Databank for store four data memory
type Databank struct {
	Coil      [0xffff]bool
	CoilIn    [0xffff]bool
	Reg       [0xffff]uint16
	RegIn     [0xffff]uint16
	muxCoil   sync.Mutex
	muxCoilIn sync.Mutex
	muxReg    sync.Mutex
	muxRegIn  sync.Mutex
}

func checkRange(start, lens int) bool {
	return start+lens >= 65536
}

// ReadCoil for read data mutexable
func (dd *Databank) ReadCoil(addr, lens uint16) ([]bool, error) {
	if checkRange(int(addr), int(lens)) {
		return nil, errors.New("access over range")
	}
	dd.muxCoil.Lock()
	result := dd.Coil[addr:(addr + lens)]
	dd.muxCoil.Unlock()
	return result, nil
}

// ReadCoilIn for read data mutexable
func (dd *Databank) ReadCoilIn(addr, lens uint16) ([]bool, error) {
	if checkRange(int(addr), int(lens)) {
		return nil, errors.New("access over range")
	}
	dd.muxCoilIn.Lock()
	result := dd.CoilIn[addr:(addr + lens)]
	dd.muxCoilIn.Unlock()
	return result, nil
}

// ReadReg for read data mutexable
func (dd *Databank) ReadReg(addr, lens uint16) ([]uint16, error) {
	if checkRange(int(addr), int(lens)) {
		return nil, errors.New("access over range")
	}
	dd.muxReg.Lock()
	result := dd.Reg[addr:(addr + lens)]
	dd.muxReg.Unlock()
	return result, nil
}

// ReadRegIn for read data mutexable
func (dd *Databank) ReadRegIn(addr, lens uint16) ([]uint16, error) {
	if checkRange(int(addr), int(lens)) {
		return nil, errors.New("access over range")
	}
	dd.muxRegIn.Lock()
	result := dd.RegIn[addr:(addr + lens)]
	dd.muxRegIn.Unlock()
	return result, nil
}

// WriteCoil for write data mutexable
func (dd *Databank) WriteCoil(addr uint16, vals []bool) error {
	if checkRange(int(addr), len(vals)) {
		return errors.New("access over range")
	}
	dd.muxCoil.Lock()
	copy(dd.Coil[addr:], vals)
	dd.muxCoil.Unlock()
	return nil
}

// WriteCoilIn for write data mutexable
func (dd *Databank) WriteCoilIn(addr uint16, vals []bool) error {
	if checkRange(int(addr), len(vals)) {
		return errors.New("access over range")
	}
	dd.muxCoilIn.Lock()
	copy(dd.CoilIn[addr:], vals)
	dd.muxCoilIn.Unlock()
	return nil
}

// WriteReg for write data mutexable
func (dd *Databank) WriteReg(addr uint16, vals []uint16) error {
	if checkRange(int(addr), len(vals)) {
		return errors.New("access over range")
	}
	dd.muxReg.Lock()
	copy(dd.Reg[addr:], vals)
	dd.muxReg.Unlock()
	return nil
}

// WriteRegIn for write data mutexable
func (dd *Databank) WriteRegIn(addr uint16, vals []uint16) error {
	if checkRange(int(addr), len(vals)) {
		return errors.New("access over range")
	}
	dd.muxRegIn.Lock()
	copy(dd.RegIn[addr:], vals)
	dd.muxRegIn.Unlock()
	return nil
}
