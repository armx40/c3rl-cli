package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/block"
)

var csv_file *os.File
var csv_writer *csv.Writer
var csv_header_set bool
var csv_delimiter string

var sqlite3_db *sql.DB
var sqlite3_db_statement *sql.Stmt
var sqlite3_db_tx *sql.Tx

func command_devices_functions_find_c3rl_device() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Debugging can be turned on; this shows some of the inner workings of the libusb package.
	ctx.Debug(*debug)

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// The usbid package can be used to print out human readable information.

		fmt.Printf("%03d.%03d %s:%s %s\n", desc.Bus, desc.Address, desc.Vendor, desc.Product, usbid.Describe(desc))
		fmt.Printf("  Protocol: %s\n", usbid.Classify(desc))

		// The configurations can be examined from the DeviceDesc, though they can only
		// be set once the device is opened.  All configuration references must be closed,
		// to free up the memory in libusb.
		// for _, cfg := range desc.Configs {
		// 	// This loop just uses more of the built-in and usbid pretty printing to list
		// 	// the USB devices.
		// 	fmt.Printf("  %s:\n", cfg)
		// 	for _, intf := range cfg.Interfaces {
		// 		fmt.Printf("    --------------\n")
		// 		for _, ifSetting := range intf.AltSettings {
		// 			fmt.Printf("    %s\n", ifSetting)
		// 			fmt.Printf("      %s\n", usbid.Classify(ifSetting))
		// 			for _, end := range ifSetting.Endpoints {
		// 				fmt.Printf("      %s\n", end)
		// 			}
		// 		}
		// 	}
		// 	fmt.Printf("    --------------\n")
		// }

		// After inspecting the descriptor, return true or false depending on whether
		// the device is "interesting" or not.  Any descriptor for which true is returned
		// opens a Device which is retuned in a slice (and must be subsequently closed).
		return false
	})

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	// OpenDevices can occasionally fail, so be sure to check its return value.
	if err != nil {
		log.Fatalf("list: %s", err)
	}

	for _, dev := range devs {
		// Once the device has been selected from OpenDevices, it is opened
		// and can be interacted with.
		_ = dev
	}
}

func command_devices_functions_find_sdcard_device() (devices []*block.Partition, err error) {
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	for _, disk := range block.Disks {
		if disk.IsRemovable {
			for _, part := range disk.Partitions {
				devices = append(devices, part)
			}

		}
	}
	return
}

func command_devices_functions_read_user_storage_device() (device *block.Partition, err error) {
	devices, err := command_devices_functions_find_sdcard_device()
	if err != nil {
		return
	}

	if len(devices) == 0 {
		return device, fmt.Errorf("no storage present")
	}
	for i, device := range devices {
		devices = append(devices, device)
		fmt.Printf("[%d]:  %v\n", i, device)
	}

	fmt.Print("Select storage device: ")

	var device_idx int
	fmt.Scan(&device_idx)

	if device_idx > len(devices)-1 {
		return device, fmt.Errorf("invalid storage device")
	}

	device = devices[device_idx]

	return
}

func command_devices_functions_get_all_log_files(device *block.Partition) (dirs []fs.DirEntry, err error) {
	dirs_, err := os.ReadDir(device.MountPoint)

	if err != nil {
		return
	}

	for i := range dirs_ {

		name := dirs_[i].Name()

		if len(name) < 14 {
			continue
		}
		if name[:9] != "data_log_" {
			continue
		}

		if name[len(name)-5:] != ".data" {
			continue
		}

		if name[9:13] == "main" {
			continue
		}

		dirs = append(dirs, dirs_[i])

	}
	return
}

func command_devices_functions_get_all_log_files_sorted(device *block.Partition, start_index int, num int) (files []fs.FileInfo, err error) {

	/* process num */

	all_log_files, err := command_devices_functions_get_all_log_files(device)
	if err != nil {
		return
	}

	if num > len(all_log_files) {
		num = len(all_log_files)
	}

	curr_index := start_index

	for len(files) < num {

		filename := filepath.Join(device.MountPoint, fmt.Sprintf("data_log_%d.data", curr_index))
		curr_index = (curr_index - 1) % 400000
		file_info, err := os.Stat(filename)
		if err != nil {
			continue
		}
		files = append(files, file_info)

	}

	return
}
