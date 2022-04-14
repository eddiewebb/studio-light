package blync

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/eddiewebb/blync-studio-light/hid"
)

const (
	blyncVendorId  = 0x0E53
	blyncProductId = 0x2517
)

const (
	BlinkOff    = 0x00
	BlinkFast   = 0x46
	BlinkMedium = 0x64
	BlinkSlow   = 0x96
)

var Red = [3]byte{0xFF, 0x00, 0x00}
var Green = [3]byte{0x00, 0xFF, 0x00}
var Blue = [3]byte{0x00, 0x00, 0xFF}

type BlyncLight struct {
	devices []hid.Device
	bytes   []byte
}

func NewBlyncLight() (blync BlyncLight) {
	blync.devices = findDevices()
	blync.bytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0xFF}
	return
}

func findDevices() []hid.Device {
	devices := []hid.Device{}
	allDeviceInfos := hid.Devices()
	blyncInfos := []hid.DeviceInfo{}
	for {
		info, more := <-allDeviceInfos
		if more {
			if !isBlyncDevice(*info) {
				//fmt.Printf("%s %s is not a BlyncLight device.\n", info.Manufacturer, info.Product)
			} else {
				blyncInfos = append(blyncInfos, *info)
				//fmt.Printf("%s %s is a BlyncLight device.\n", info.Manufacturer, info.Product)
				//fmt.Printf("%s\n",info)
			}
		} else {
			break
		}
	}
	if len(blyncInfos) == 0 {
		fmt.Println("No BlyncLights found.")
		os.Exit(1)
	}

	sort.Sort(byLocation(blyncInfos))

	for i := 0; i < len(blyncInfos); i++ {
		device, error := blyncInfos[i].Open()
		if error != nil {
			fmt.Println(error)
		} else {
			fmt.Printf("added usb with path %s as index:%d, ID: %d\n", blyncInfos[i].Path, i, i+1)
		}

		devices = append(devices, device)

	}
	return devices
}

func isBlyncDevice(deviceInfo hid.DeviceInfo) bool {
	// from forums: "Blync creates 2 HID devices and the only way to find out the right device is the MaxFeatureReportLength = 0"
	if deviceInfo.VendorId == blyncVendorId && deviceInfo.ProductId == blyncProductId && deviceInfo.FeatureReportLength == 0 {
		return true
	}
	return false
}

func (b BlyncLight) sendFeatureReport(id int) {
	if id == 0 {
		for _, device := range b.devices {
			b.write(device)
		}
	} else {
		index := id - 1
		device := b.devices[index]
		b.write(device)
	}
}

func (b BlyncLight) write(device hid.Device) {
	error := device.Write(b.bytes)
	if error != nil {
		fmt.Println(error)
	}
}

func (b BlyncLight) Close(id int) {
	b.Reset(id)
	for _, device := range b.devices {
		device.Close()
	}
}

func (b BlyncLight) Reset(id int) {
	b.bytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0xFF}
	b.sendFeatureReport(id)
}

// color[0] = r
// color[1] = g
// color[2] = b
func (b BlyncLight) SetColor(color [3]byte, id int) {
	b.bytes[1] = color[0]
	b.bytes[2] = color[2] // They reverse g and b
	b.bytes[3] = color[1]
	b.sendFeatureReport(id)
}

func (b BlyncLight) SetBlinkRate(rate byte, id int) {
	b.bytes[4] = rate
	b.sendFeatureReport(id)
}

func (b BlyncLight) FlashOrder() {
	for i := 0; i < len(b.devices); i++ {
		id := i + 1
		b.SetBlinkRate(BlinkMedium, id)
		b.SetColor(Red, id)
		b.sendFeatureReport(id)
		fmt.Printf("Flashing blync ID: %d\n", id)
		time.Sleep(time.Second * 2)
		b.Reset(id)
	}
}

//16-30 play a tune single time
//49-59 plays never ending versions of the tunes
func (b BlyncLight) Play(mp3 byte, id int) {
	b.bytes[5] = mp3
	b.sendFeatureReport(id)
}

func (b BlyncLight) StopPlay(id int) {
	b.bytes[5] = 0x00
	b.sendFeatureReport(id)
}

// Methods needed for sorting
type byLocation []hid.DeviceInfo

func (a byLocation) Len() int {
	return len(a)
}

func (a byLocation) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byLocation) Less(i, j int) bool {
	return a[i].Path < a[j].Path
}
