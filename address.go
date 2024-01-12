package memory

import (
	"unsafe"
)

type Address[t comparable] struct {
	Process *Process
	Data    uintptr
}

func NewAddress[t comparable](process *Process, address uintptr) *Address[t] {
	return &Address[t]{process, address}
}

func (address *Address[t]) Write(data t) error {
	if r, _, err := writeProcessMemory.Call(uintptr(address.Process.handle), address.Data, uintptr(unsafe.Pointer(&data)), unsafe.Sizeof(data), 0); r == 0 {
		return err
	}

	return nil
}

func (address *Address[t]) Read() (t, error) {
	var data t

	if r, _, err := readProcessMemory.Call(uintptr(address.Process.handle), address.Data, uintptr(unsafe.Pointer(&data)), unsafe.Sizeof(data), 0); r == 0 {
		return data, err
	}

	return data, nil
}
