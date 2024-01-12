package memory

import "testing"

func TestMemory(t *testing.T) {
	process, err := NewProcessByName("ac_client.exe")
	if err != nil {
		t.Error(err)
	}

	if err := process.Open(); err != nil {
		t.Error(err)
	}

	defer process.Close()

	module, err := NewModule(process, "ac_client.exe")
	if err != nil {
		t.Error(err)
	}

	ammoPointer := NewPointer[uint32](module, 0x0016F338, []uintptr{0x8, 0x2C, 0x140})

	ammoAddress, err := ammoPointer.CalculateOffsets()
	if err != nil {
		t.Error(err)
	}

	ammo, err := ammoAddress.Read()
	if err != nil {
		t.Error(err)
	}

	if err := ammoAddress.Write(54); err != nil {
		t.Error(err)
	}

	t.Log(ammo)
}
