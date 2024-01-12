package memory

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Process struct {
	name   string
	id     uint32
	handle windows.Handle
}

func NewProcessByName(name string) (*Process, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	for windows.Process32Next(snapshot, &entry) == nil {
		for i := 0; i < 260; i++ {
			if entry.ExeFile[i] == 0 {
				if syscall.UTF16ToString(entry.ExeFile[:i]) == name {
					return &Process{name: name, id: entry.ProcessID}, nil
				}
			}
		}
	}

	return nil, errors.New("process not found")
}

func NewProcessById(id uint32) (*Process, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	for windows.Process32Next(snapshot, &entry) == nil {
		for i := 0; i < 260; i++ {
			if entry.ExeFile[i] == 0 {
				if entry.ProcessID == id {
					return &Process{name: syscall.UTF16ToString(entry.ExeFile[:i]), id: entry.ProcessID}, nil
				}
			}
		}
	}

	return nil, errors.New("process not found")
}

func (process *Process) GetName() string {
	return process.name
}

func (process *Process) GetId() uint32 {
	return process.id
}

func (process *Process) Open() error {
	handle, err := windows.OpenProcess(processAllAccess, false, process.id)
	if err != nil {
		return err
	}

	process.handle = handle

	return nil
}

func (process *Process) Close() error {
	err := windows.CloseHandle(process.handle)
	if err != nil {
		return err
	}

	process.handle = 0

	return nil
}

func (process Process) Suspend() error {
	if r, _, err := ntSuspendProcess.Call(uintptr(process.handle)); r == 0 {
		return err
	}

	return nil
}

func (process Process) Resume() error {
	if r, _, err := ntResumeProcess.Call(uintptr(process.handle)); r == 0 {
		return err
	}

	return nil
}
