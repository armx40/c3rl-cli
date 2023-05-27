package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/****************************** device id **************************************/

type host_device_payloads_information_data_device_id_t struct {
	DBUSMachineID string
}

func (h *host_device_payloads_information_data_device_id_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("cat", []string{"/var/lib/dbus/machine-id"})
	if err != nil {
		return
	}

	h.DBUSMachineID = strings.TrimSpace(string(cmd_out))
	return
}

/****************************** lspci ******************************************/
type host_device_payloads_information_data_lspci_t struct {
	Devices []string
}

func (h *host_device_payloads_information_data_lspci_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("lspci", []string{})
	if err != nil {
		return
	}

	devices := strings.Split(string(cmd_out), "\n")

	for i := range devices {
		devices[i] = strings.TrimSpace(devices[i])
	}

	h.Devices = devices
	return

}

/****************************** lscpu ******************************************/
type host_device_payloads_information_data_lscpu_single_block_t struct {
	Range       string
	Size        uint64
	State       string
	IsRemovable bool
	Block       string
	Node        string
}

type host_device_payloads_information_data_lscpu_t struct {
	CPUs          []host_device_payloads_information_data_lscpu_single_block_t
	CPUName       string
	Architecture  string
	ByteOrder     string
	NumberOfCores uint16
	Vendor        string
}

func (h *host_device_payloads_information_data_lscpu_t) get() (err error) {

	/* get cpu name */
	cmd_out, err := helper_function_get_command_output("lscpu", []string{"-J"})
	if err != nil {
		return
	}
	/**/

	lscpu_data := make(map[string][](map[string]string))

	err = json.Unmarshal(cmd_out, &lscpu_data)

	if err != nil {
		return
	}

	lscpu_key, ok := lscpu_data["lscpu"]
	if !ok {
		log.Println(lscpu_data["lscpu"])
		err = fmt.Errorf("invalid response from lscpu")
		return
	}

	for i := range lscpu_key {

		if (lscpu_key[i]["field"]) == "Architecture:" {
			h.Architecture = lscpu_key[i]["data"]
		}

		if (lscpu_key[i]["field"]) == "Byte Order:" {
			h.ByteOrder = lscpu_key[i]["data"]
		}

		if (lscpu_key[i]["field"]) == "Model name:" {
			h.CPUName = lscpu_key[i]["data"]
		}

		if (lscpu_key[i]["field"]) == "Vendor ID:" {
			h.Vendor = lscpu_key[i]["data"]
		}

		if (lscpu_key[i]["field"]) == "CPU(s):" {

			number_of_cores, err := strconv.ParseUint(lscpu_key[i]["data"], 10, 16)

			if err != nil {
				err = fmt.Errorf("failed to get number of cores")
				return err
			}

			h.NumberOfCores = uint16(number_of_cores)
		}
	}

	/* get core wise data */

	// cmd_out, err := helper_function_get_command_output("lsmem", []string{"-r", "-n", "-b", "-o", "RANGE,SIZE,STATE,REMOVABLE,BLOCK,NODE"})
	// if err != nil {
	// 	return
	// }

	// devices_lines := strings.Split(string(cmd_out), "\n")

	// for i := range devices_lines {

	// 	line := strings.TrimSpace(devices_lines[i])
	// 	if line == "" {
	// 		continue
	// 	}
	// 	/* get values from single line */
	// 	line_values := strings.Split(line, " ")

	// 	if len(line_values) != 6 {
	// 		err = fmt.Errorf("invalid number values for mem")
	// 		return
	// 	}

	// 	size_bytes, err := strconv.ParseUint(line_values[1], 10, 64)

	// 	if err != nil {
	// 		err = fmt.Errorf("failed to get blk size")
	// 		return err
	// 	}

	// 	h.Devices = append(h.Devices, host_device_payloads_information_data_lsmem_single_block_t{
	// 		Range:       line_values[0],
	// 		Size:        size_bytes,
	// 		State:       line_values[2],
	// 		Block:       line_values[4],
	// 		Node:        line_values[5],
	// 		IsRemovable: line_values[3] == "yes",
	// 	})
	// }

	return

}

/****************************** lsmem ******************************************/

type host_device_payloads_information_data_lsmem_single_block_t struct {
	Range       string
	Size        uint64
	State       string
	IsRemovable bool
	Block       string
	Node        string
}

type host_device_payloads_information_data_lsmem_t struct {
	Devices []host_device_payloads_information_data_lsmem_single_block_t
}

func (h *host_device_payloads_information_data_lsmem_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("lsmem", []string{"-r", "-n", "-b", "-o", "RANGE,SIZE,STATE,REMOVABLE,BLOCK,NODE"})
	if err != nil {
		return
	}

	devices_lines := strings.Split(string(cmd_out), "\n")

	for i := range devices_lines {

		line := strings.TrimSpace(devices_lines[i])
		if line == "" {
			continue
		}
		/* get values from single line */
		line_values := strings.Split(line, " ")

		if len(line_values) != 6 {
			err = fmt.Errorf("invalid number values for mem")
			return
		}

		size_bytes, err := strconv.ParseUint(line_values[1], 10, 64)

		if err != nil {
			err = fmt.Errorf("failed to get blk size")
			return err
		}

		h.Devices = append(h.Devices, host_device_payloads_information_data_lsmem_single_block_t{
			Range:       line_values[0],
			Size:        size_bytes,
			State:       line_values[2],
			Block:       line_values[4],
			Node:        line_values[5],
			IsRemovable: line_values[3] == "yes",
		})
	}

	return

}

/****************************** lsusb ******************************************/

type host_device_payloads_information_data_lsusb_t struct {
	Devices []string
}

func (h *host_device_payloads_information_data_lsusb_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("lsusb", []string{})
	if err != nil {
		return
	}

	devices := strings.Split(string(cmd_out), "\n")

	for i := range devices {
		devices[i] = strings.TrimSpace(devices[i])
	}

	h.Devices = devices
	return
}

/****************************** lsblk ******************************************/
type host_device_payloads_information_data_lsblk_single_block_t struct {
	Name        string
	Size        uint64
	Mountpoint  string
	IsRemovable bool
}

type host_device_payloads_information_data_lsblk_t struct {
	Devices []host_device_payloads_information_data_lsblk_single_block_t
}

func (h *host_device_payloads_information_data_lsblk_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("lsblk", []string{"-r", "-n", "-b", "-o", "NAME,SIZE,MOUNTPOINT,RM"})
	if err != nil {
		return
	}

	devices_lines := strings.Split(string(cmd_out), "\n")

	for i := range devices_lines {

		line := strings.TrimSpace(devices_lines[i])
		if line == "" {
			continue
		}
		/* get values from single line */
		line_values := strings.Split(line, " ")

		if len(line_values) != 4 {
			err = fmt.Errorf("invalid number values for blk")
			return
		}

		size_bytes, err := strconv.ParseUint(line_values[1], 10, 64)

		if err != nil {
			err = fmt.Errorf("failed to get blk size")
			return err
		}

		isremovable_int, err := strconv.ParseUint(line_values[3], 10, 8)

		if err != nil {
			err = fmt.Errorf("failed to get if device is removable")
			return err
		}

		h.Devices = append(h.Devices, host_device_payloads_information_data_lsblk_single_block_t{
			Name:        line_values[0],
			Size:        size_bytes,
			Mountpoint:  line_values[2],
			IsRemovable: isremovable_int == 1,
		})
	}

	return

}

/****************************** uname ******************************************/

type host_device_payloads_information_data_uname_t struct {
	KernelName       string
	NodeName         string
	KernelRelease    string
	KernelVersion    string
	Machine          string
	Processor        string
	HardwarePlatform string
	OperatingSystem  string
}

func (h *host_device_payloads_information_data_uname_t) get() (err error) {

	cmd_out, err := helper_function_get_command_output("uname", []string{"-s"})
	if err != nil {
		return
	}
	h.KernelName = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-n"})
	if err != nil {
		return
	}
	h.NodeName = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-r"})
	if err != nil {
		return
	}
	h.KernelRelease = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-v"})
	if err != nil {
		return
	}
	h.KernelVersion = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-m"})
	if err != nil {
		return
	}
	h.Machine = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-p"})
	if err != nil {
		return
	}
	h.Processor = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-i"})
	if err != nil {
		return
	}
	h.HardwarePlatform = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-o"})
	if err != nil {
		return
	}
	h.OperatingSystem = string(cmd_out)

	return
}

/************************************************************************/

type host_device_payloads_information_data_t struct {
	Uname    host_device_payloads_information_data_uname_t
	USB      host_device_payloads_information_data_lsusb_t
	PCI      host_device_payloads_information_data_lspci_t
	DeviceID host_device_payloads_information_data_device_id_t
	Blocks   host_device_payloads_information_data_lsblk_t
	Memory   host_device_payloads_information_data_lsmem_t
	CPU      host_device_payloads_information_data_lscpu_t
}

func (h *host_device_payloads_information_data_t) get() (err error) {

	/* uname */
	err = h.Uname.get()
	if err != nil {
		return
	}
	/**/

	/* usb */
	err = h.USB.get()
	if err != nil {
		return
	}
	/**/

	/* pci */
	err = h.PCI.get()
	if err != nil {
		return
	}
	/**/

	/* device id */
	err = h.DeviceID.get()
	if err != nil {
		return
	}
	/**/

	/* blocks */
	err = h.Blocks.get()
	if err != nil {
		return
	}
	/**/

	/* memory */
	err = h.Memory.get()
	if err != nil {
		return
	}
	/**/

	/* cpu */
	err = h.CPU.get()
	if err != nil {
		return
	}
	/**/

	return
}
