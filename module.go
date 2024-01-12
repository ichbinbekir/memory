package memory

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Module struct {
	process *Process
	name    string
	address uintptr
}

func NewModule(process *Process, name string) (*Module, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, process.id)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snapshot)

	var entry windows.ModuleEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	for windows.Module32Next(snapshot, &entry) == nil {
		for i := 0; i < 256; i++ {
			if entry.Module[i] == 0 {
				if syscall.UTF16ToString(entry.Module[:i]) == name {
					return &Module{process, name, entry.ModBaseAddr}, nil
				}
			}
		}
	}

	return nil, errors.New("module not found")
}

func (module *Module) GetProcess() *Process {
	return module.process
}

func (module *Module) GetName() string {
	return module.name
}

func (module *Module) GetAddress() uintptr {
	return module.address
}
