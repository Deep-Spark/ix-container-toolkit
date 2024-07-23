/**
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package modifier

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
	"gitee.com/deep-spark/ix-container-runtime/internal/config/image"
	"gitee.com/deep-spark/ix-container-runtime/internal/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

var (
	visibleDevicesEnvvar = "IX_VISIBLE_DEVICES"
	deviceName           = "iluvatar"
	devicePath           = "/dev"

	wildcardDevice = "a"
	blockDevice    = "b"
	charDevice     = "c"
	fifoDevice     = "p"

	ErrNotADevice = errors.New("not a device node")
)

type graphicsModifier struct {
	addDevice []specs.LinuxDevice
}

type IndexDevice struct {
	specs.LinuxDevice
	ixml.Device
	Index uint
}

func (g graphicsModifier) Modify(spec *specs.Spec) error {
	var tmpPtr *specs.LinuxDevice
	if g.addDevice == nil {
		return nil
	}
	if len(g.addDevice) == 0 {
		return nil
	}

	for _, d := range g.addDevice {
		spec.Linux.Devices = append(spec.Linux.Devices, d)
		tmpPtr = new(specs.LinuxDevice)
		*tmpPtr = d
		newDeviceCgroup := specs.LinuxDeviceCgroup{
			Allow:  true,
			Type:   tmpPtr.Type,
			Major:  &tmpPtr.Major,
			Minor:  &tmpPtr.Minor,
			Access: "rwm",
		}
		spec.Linux.Resources.Devices = append(spec.Linux.Resources.Devices, newDeviceCgroup)
	}
	return nil
}

func searchDevice() map[int]specs.LinuxDevice {
	ret := make(map[int]specs.LinuxDevice)
	libRegEx, e := regexp.Compile(deviceName + "[0-9]")
	if e != nil {
		log.Fatal(e)
	}
	e = filepath.Walk(devicePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			absolutePath := path
			var stat unix.Stat_t
			if err := unix.Lstat(path, &stat); err != nil {
				return err
			}

			var (
				devNumber = uint64(stat.Rdev)
				major     = unix.Major(devNumber)
				minor     = unix.Minor(devNumber)
				mode      = stat.Mode
				devType   string
			)

			switch mode & unix.S_IFMT {
			case unix.S_IFBLK:
				devType = blockDevice
			case unix.S_IFCHR:
				devType = charDevice
			case unix.S_IFIFO:
				devType = fifoDevice
			default:
				return ErrNotADevice
			}
			fm := os.FileMode(mode &^ unix.S_IFMT)
			dev := specs.LinuxDevice{
				Type:     devType,
				Path:     absolutePath,
				Major:    int64(major),
				Minor:    int64(minor),
				FileMode: &fm,
				UID:      &stat.Uid,
				GID:      &stat.Gid,
			}

			ret[int(minor)] = dev
		}
		return nil
	})

	if e != nil {
		log.Printf("error for walk:%v\n", e)
		return nil
	}

	return ret
}

func generate_dev_from_string(devmap map[uint]IndexDevice, val string, mountIdx int) *specs.LinuxDevice {
	var ret specs.LinuxDevice
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Can't transfer %v to int type", val)
		return nil
	}

	dev, ok := devmap[uint(i)]
	if !ok {
		log.Printf("Wrong parameter: %v", val)
		return nil
	}
	ret = dev.LinuxDevice
	strIdx := strconv.Itoa(mountIdx)
	ret.Path = devicePath + "/" + deviceName + strIdx
	return &ret
}

func getdevice(devmap map[uint]IndexDevice, cudaImage image.CUDA) []specs.LinuxDevice {
	var ret []specs.LinuxDevice
	devices := cudaImage.DevicesFromEnvvars(visibleDevicesEnvvar)
	if len(devices.List()) == 0 {
		return nil
	} else if len(devices.List()) == 1 {
		val := devices.List()[0]
		switch val {
		case "all":
			for _, dev := range devmap {
				ret = append(ret, dev.LinuxDevice)
			}
			return ret
		case "void":
			return nil
		case "none":
			return nil
		default:
			dev := generate_dev_from_string(devmap, val, 0)
			if dev != nil {
				ret = append(ret, *dev)
				return ret
			} else {
				return nil
			}
		}
	}

	for mountIdx, v := range devices.List() {
		dev := generate_dev_from_string(devmap, v, mountIdx)
		if dev != nil {
			ret = append(ret, *dev)
		}
	}
	return ret
}

func buildMountDevice(index int, dev specs.LinuxDevice) specs.LinuxDevice {
	devIdx := strconv.Itoa(index)
	mountPath := devicePath + "/" + deviceName + devIdx
	return specs.LinuxDevice{Type: dev.Type,
		Path:     mountPath,
		Major:    dev.Major,
		Minor:    dev.Minor,
		FileMode: dev.FileMode,
		UID:      dev.UID,
		GID:      dev.GID,
	}
}

// requiresGraphicsModifier determines whether a graphics modifier is required.
func buildMap(librarypath string) map[uint]IndexDevice {
	IndexMap := make(map[uint]IndexDevice)
	var ret ixml.Return

	devs := searchDevice()

	if librarypath != "" {
		ret = ixml.AbsInit(librarypath)
	} else {
		ret = ixml.Init()
	}
	if ret != ixml.SUCCESS {
		log.Printf("Unable to initialize IXML:%v\n", ret)
		log.Printf("librarypath:%v", librarypath)
		return nil
	}
	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Printf("failed to get count:%v\n", ret)
		return nil
	}

	log.Printf("count: %d\n", count)

	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)

		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}

		MinorID, ret := device.GetMinorNumber()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get id of device at index %v: %v", MinorID, ret)
		}
		IndexMap[i] = IndexDevice{
			Index:       i,
			LinuxDevice: devs[MinorID],
			Device:      device,
		}
	}

	return IndexMap
}

func NewGraphicsModifier(image image.CUDA) oci.SpecModifier {
	devMap := buildMap(image.Cfg.LibraryPath)
	if devMap == nil {
		log.Printf("No graphics modifier required\n")
		return nil
	}

	ret := graphicsModifier{
		addDevice: getdevice(devMap, image),
	}

	return ret
}
