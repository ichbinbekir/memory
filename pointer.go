package memory

import "unsafe"

type Pointer[t comparable] struct {
	Module  *Module
	Address uintptr
	Offsets []uintptr
}

func NewPointer[t comparable](module *Module, address uintptr, offsets []uintptr) *Pointer[t] {
	return &Pointer[t]{module, address, offsets}
}

func (pointer *Pointer[t]) CalculateOffsets() (*Address[t], error) {
	address := pointer.Module.address + pointer.Address

	for _, offset := range pointer.Offsets {
		var temp uintptr

		if r, _, err := readProcessMemory.Call(uintptr(pointer.Module.process.handle), address, uintptr(unsafe.Pointer(&temp)), 4, 0); r == 0 {
			return nil, err
		}

		address = temp + offset
	}

	return NewAddress[t](pointer.Module.process, address), nil
}
