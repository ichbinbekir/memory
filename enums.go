package memory

import "golang.org/x/sys/windows"

var (
	ntdll    = windows.NewLazyDLL("ntdll.dll")
	kernel32 = windows.NewLazyDLL("kernel32.dll")
)

var (
	ntSuspendProcess = ntdll.NewProc("NtSuspendProcess")
	ntResumeProcess  = ntdll.NewProc("NtResumeProcess")
)

var (
	writeProcessMemory = kernel32.NewProc("WriteProcessMemory")
	readProcessMemory  = kernel32.NewProc("ReadProcessMemory")
	createRemoteThread = kernel32.NewProc("CreateRemoteThread")
)

const (
	processAllAccess = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xFFFF
)
