package memory

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetAllProcesses() ([]*Process, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	var processes []*Process

	for windows.Process32Next(snapshot, &entry) == nil {
		for i := 0; i < 260; i++ {
			if entry.ExeFile[i] == 0 {
				processes = append(processes, &Process{name: syscall.UTF16ToString(entry.ExeFile[:i]), id: entry.ProcessID})

				break
			}
		}
	}

	return processes, nil
}

func GetAllModules(process *Process) ([]*Module, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, process.id)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snapshot)

	var entry windows.ModuleEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	var modules []*Module

	for windows.Module32Next(snapshot, &entry) == nil {
		for i := 0; i < 256; i++ {
			if entry.Module[i] == 0 {
				modules = append(modules, &Module{process, syscall.UTF16ToString(entry.Module[:i]), entry.ModBaseAddr})

				break
			}
		}
	}

	return modules, nil
}

/*func (process Process) CreateRemoteThread(baseAdderss, parameterAddress uintptr) error {
	if r, _, err := createRemoteThread.Call(uintptr(process.handle), 0, 0, baseAdderss, parameterAddress, 0, 0); r == 0 {
		return err
	}

	return nil
}*/
