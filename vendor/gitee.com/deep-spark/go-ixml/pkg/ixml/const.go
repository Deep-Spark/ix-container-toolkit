// WARNING: This file has automatically been generated on Tue, 29 Oct 2024 14:18:20 CST.
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package ixml

/*
#cgo LDFLAGS: -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#cgo CFLAGS: -DNVML_NO_UNVERSIONED_FUNC_DEFS=1
#include "api.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

const (
	// DEVICE_UUID_BUFFER_SIZE as defined in ixml/api.h:231
	DEVICE_UUID_BUFFER_SIZE = 80
	// SYSTEM_DRIVER_VERSION_BUFFER_SIZE as defined in ixml/api.h:236
	SYSTEM_DRIVER_VERSION_BUFFER_SIZE = 80
	// DEVICE_NAME_BUFFER_SIZE as defined in ixml/api.h:241
	DEVICE_NAME_BUFFER_SIZE = 64
	// DEVICE_NAME_V2_BUFFER_SIZE as defined in ixml/api.h:246
	DEVICE_NAME_V2_BUFFER_SIZE = 96
	// DEVICE_PCI_BUS_ID_BUFFER_SIZE as defined in ixml/api.h:251
	DEVICE_PCI_BUS_ID_BUFFER_SIZE = 32
	// DEVICE_PCI_BUS_ID_BUFFER_V2_SIZE as defined in ixml/api.h:256
	DEVICE_PCI_BUS_ID_BUFFER_V2_SIZE = 16
	// NO_UNVERSIONED_FUNC_DEFS as defined in go-ixml/<predefined>:349
	NO_UNVERSIONED_FUNC_DEFS = 1
)

// Return as declared in ixml/api.h:128
type Return int32

// Return enumeration from ixml/api.h:128
const (
	SUCCESS                         Return = iota
	ERROR_UNINITIALIZED             Return = 1
	ERROR_INVALID_ARGUMENT          Return = 2
	ERROR_NOT_SUPPORTED             Return = 3
	ERROR_NO_PERMISSION             Return = 4
	ERROR_ALREADY_INITIALIZED       Return = 5
	ERROR_NOT_FOUND                 Return = 6
	ERROR_INSUFFICIENT_SIZE         Return = 7
	ERROR_INSUFFICIENT_POWER        Return = 8
	ERROR_DRIVER_NOT_LOADED         Return = 9
	ERROR_TIMEOUT                   Return = 10
	ERROR_IRQ_ISSUE                 Return = 11
	ERROR_LIBRARY_NOT_FOUND         Return = 12
	ERROR_FUNCTION_NOT_FOUND        Return = 13
	ERROR_CORRUPTED_INFOROM         Return = 14
	ERROR_GPU_IS_LOST               Return = 15
	ERROR_RESET_REQUIRED            Return = 16
	ERROR_OPERATING_SYSTEM          Return = 17
	ERROR_LIB_RM_VERSION_MISMATCH   Return = 18
	ERROR_IN_USE                    Return = 19
	ERROR_MEMORY                    Return = 20
	ERROR_NO_DATA                   Return = 21
	ERROR_VGPU_ECC_NOT_SUPPORTED    Return = 22
	ERROR_INSUFFICIENT_RESOURCES    Return = 23
	ERROR_FREQ_NOT_SUPPORTED        Return = 24
	ERROR_ARGUMENT_VERSION_MISMATCH Return = 25
	ERROR_UNKNOWN                   Return = 999
)

// TemperatureSensors as declared in ixml/api.h:175
type TemperatureSensors int32

// TemperatureSensors enumeration from ixml/api.h:175
const (
	TEMPERATURE_GPU   TemperatureSensors = iota
	TEMPERATURE_COUNT TemperatureSensors = 1
)

// ClockType as declared in ixml/api.h:191
type ClockType int32

// ClockType enumeration from ixml/api.h:191
const (
	CLOCK_GRAPHICS ClockType = iota
	CLOCK_SM       ClockType = 1
	CLOCK_MEM      ClockType = 2
	CLOCK_VIDEO    ClockType = 3
	CLOCK_COUNT    ClockType = 4
)