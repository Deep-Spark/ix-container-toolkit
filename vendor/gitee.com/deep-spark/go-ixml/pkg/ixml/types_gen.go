// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs types.go

package ixml

type Device struct {
	Handle *_Ctype_struct_nvmlDevice_st
}

type EventSet struct {
	Handle *_Ctype_struct_nvmlEventSet_st
}

type Memory struct {
	Total uint64
	Free  uint64
	Used  uint64
}

type Memory_v2 struct {
	Version  uint32
	Total    uint64
	Reserved uint64
	Free     uint64
	Used     uint64
}

type Utilization struct {
	Gpu    uint32
	Memory uint32
}

type ProcessInfo struct {
	Pid                      uint32
	UsedGpuMemory            uint64
	GpuInstanceId            uint32
	ComputeInstanceId        uint32
	UsedGpuCcProtectedMemory uint64
}

type ProcessInfo_v1 struct {
	Pid           uint32
	UsedGpuMemory uint64
}

type PciInfo struct {
	BusIdLegacy    [16]int8
	Domain         uint32
	Bus            uint32
	Device         uint32
	PciDeviceId    uint32
	PciSubSystemId uint32
	BusId          [32]int8
}
