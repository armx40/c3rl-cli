package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v3/cpu"
)

/****************************** device id **************************************/

type host_device_payloads_information_data_device_id_t struct {
	MachineID string
}

func (h *host_device_payloads_information_data_device_id_t) get() (err error) {

	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	h.MachineID = id

	return
}

/****************************** baseboard ******************************************/

type host_device_payloads_information_data_baseboard_t struct {
	AssetTag     string
	SerialNumber string
	Vendor       string
	Version      string
	Product      string
}

func (h *host_device_payloads_information_data_baseboard_t) get() (err error) {

	baseboard, err := ghw.Baseboard()
	if err != nil {
		fmt.Printf("Error getting baseboard info: %v", err)
		return
	}

	h.AssetTag = baseboard.AssetTag
	h.SerialNumber = baseboard.SerialNumber
	h.Vendor = baseboard.Vendor
	h.Version = baseboard.Version
	h.Product = baseboard.Product

	return
}

/****************************** bios ******************************************/

type host_device_payloads_information_data_bios_t struct {
	Vendor  string
	Version string
	Date    string
}

func (h *host_device_payloads_information_data_bios_t) get() (err error) {

	bios, err := ghw.BIOS()
	if err != nil {
		fmt.Printf("Error getting bios info: %v", err)
		return
	}

	h.Date = bios.Date
	h.Vendor = bios.Vendor
	h.Version = bios.Version

	return
}

/****************************** gpu ******************************************/

type host_device_payloads_information_data_gpu_t struct {
	Devices string // for now just store it in string
}

func (h *host_device_payloads_information_data_gpu_t) get() (err error) {

	gpu, err := ghw.GPU()
	if err != nil {
		fmt.Printf("Error getting gpu info: %v", err)
		return
	}
	h.Devices = gpu.JSONString(false)

	return

}

/****************************** network ******************************************/

type host_device_payloads_information_data_network_t struct {
	Devices string // for now just store it in string
}

func (h *host_device_payloads_information_data_network_t) get() (err error) {

	network, err := ghw.Network()
	if err != nil {
		fmt.Printf("Error getting network info: %v", err)
		return
	}

	h.Devices = network.JSONString(false)

	return

}

/****************************** pci ******************************************/
type host_device_payloads_information_data_pci_t struct {
	Devices string // for now just store it in string
}

func (h *host_device_payloads_information_data_pci_t) get() (err error) {

	if runtime.GOOS == "darwin" {

	} else {
		pci, err := ghw.PCI()
		if err != nil {
			fmt.Printf("Error getting pci info: %v", err)
			return err
		}

		h.Devices = pci.JSONString(false)

	}

	return

}

/****************************** cpu ******************************************/

type host_device_payloads_information_data_cpu_core_t struct {
	NumThreads        uint32
	LogicalProcessors []int
}
type host_device_payloads_information_data_cpu_single_block_t struct {
	NumCores     uint32
	NumThreads   uint32
	Vendor       string
	Model        string
	Capabilities []string
	Cores        []host_device_payloads_information_data_cpu_core_t
}

type host_device_payloads_information_data_cpu_t struct {
	CPUs         []host_device_payloads_information_data_cpu_single_block_t
	TotalCores   uint32
	TotalThreads uint32
}

func (h *host_device_payloads_information_data_cpu_t) get() (err error) {

	if runtime.GOOS == "darwin" {

		cpu_count, err := cpu.Counts(false)
		if err != nil {
			return err
		}

		h.TotalCores = uint32(cpu_count)

		cpu_info, err := cpu.Info()
		if err != nil {
			return err
		}
		h.TotalThreads = uint32(len(cpu_info))

		for i := range cpu_info {

			h.CPUs = append(h.CPUs, host_device_payloads_information_data_cpu_single_block_t{
				NumCores:     uint32(cpu_info[i].Cores),
				NumThreads:   uint32(cpu_info[i].Cores),
				Vendor:       cpu_info[i].VendorID,
				Model:        cpu_info[i].Model,
				Capabilities: cpu_info[i].Flags,
			})

		}

		return err

	} else {
		cpu, err := ghw.CPU()
		if err != nil {
			fmt.Printf("Error getting cpu info: %v", err)
			return err
		}

		h.TotalCores = cpu.TotalCores
		h.TotalThreads = cpu.TotalThreads

		for i := range cpu.Processors {

			cores := []host_device_payloads_information_data_cpu_core_t{}

			for j := range cpu.Processors[i].Cores {
				cores = append(cores, host_device_payloads_information_data_cpu_core_t{
					NumThreads:        cpu.Processors[i].Cores[j].NumThreads,
					LogicalProcessors: cpu.Processors[i].Cores[j].LogicalProcessors,
				})
			}

			h.CPUs = append(h.CPUs, host_device_payloads_information_data_cpu_single_block_t{
				NumCores:     cpu.Processors[i].NumCores,
				NumThreads:   cpu.Processors[i].NumThreads,
				Vendor:       cpu.Processors[i].Vendor,
				Model:        cpu.Processors[i].Model,
				Capabilities: cpu.Processors[i].Capabilities,
			})
		}

	}

	return

}

/****************************** memory ******************************************/

type host_device_payloads_information_data_memory_single_block_t struct {
	Size         int64
	Label        string
	Vendor       string
	Location     string
	SerialNumber string
}

type host_device_payloads_information_data_memory_t struct {
	Devices     []host_device_payloads_information_data_memory_single_block_t
	TotalSize   int64
	TotalUsable int64
}

func (h *host_device_payloads_information_data_memory_t) get() (err error) {

	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
		return
	}

	h.TotalSize = memory.TotalPhysicalBytes
	h.TotalSize = memory.TotalUsableBytes

	for i := range memory.Modules {
		h.Devices = append(h.Devices, host_device_payloads_information_data_memory_single_block_t{
			Size:         memory.Modules[i].SizeBytes,
			Label:        memory.Modules[i].Label,
			Vendor:       memory.Modules[i].Vendor,
			Location:     memory.Modules[i].Location,
			SerialNumber: memory.Modules[i].SerialNumber,
		})
	}

	return

}

/****************************** lsusb ******************************************/

// type host_device_payloads_information_data_lsusb_t struct {
// 	Devices []string
// }

// func (h *host_device_payloads_information_data_lsusb_t) get() (err error) {

// 	blocks, err := ghw.PCI()
// 	if err != nil {
// 		fmt.Printf("Error getting memory info: %v", err)
// 	}

// 	cmd_out, err := helper_function_get_command_output("lsusb", []string{})
// 	if err != nil {
// 		return
// 	}

// 	devices := strings.Split(string(cmd_out), "\n")

// 	for i := range devices {
// 		devices[i] = strings.TrimSpace(devices[i])
// 	}

// 	h.Devices = devices
// 	return
// }

/****************************** blocks ******************************************/
type host_device_payloads_information_data_block_single_partition_t struct {
	Name        string
	Label       string
	Size        uint64
	Mountpoint  string
	IsRemovable bool
}

type host_device_payloads_information_data_block_single_disk_t struct {
	Name         string
	Size         uint64
	IsRemovable  bool
	Vendor       string
	Model        string
	SerialNumber string
	Partitions   []host_device_payloads_information_data_block_single_partition_t
}

type host_device_payloads_information_data_block_t struct {
	Disks []host_device_payloads_information_data_block_single_disk_t
}

func (h *host_device_payloads_information_data_block_t) get() (err error) {

	blocks, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block info: %v", err)
		return
	}

	for i := range blocks.Disks {

		/* get partitions */

		tmp_partitions := []host_device_payloads_information_data_block_single_partition_t{}

		for j := range blocks.Disks[i].Partitions {
			tmp_partitions = append(tmp_partitions, host_device_payloads_information_data_block_single_partition_t{
				Name:        blocks.Disks[i].Partitions[j].Name,
				Label:       blocks.Disks[i].Partitions[j].Label,
				Size:        blocks.Disks[i].Partitions[j].SizeBytes,
				Mountpoint:  blocks.Disks[i].Partitions[j].MountPoint,
				IsRemovable: blocks.Disks[i].Partitions[j].Disk.IsRemovable,
			})

		}

		h.Disks = append(h.Disks, host_device_payloads_information_data_block_single_disk_t{
			Name:         blocks.Disks[i].Name,
			Size:         blocks.Disks[i].SizeBytes,
			IsRemovable:  blocks.Disks[i].IsRemovable,
			Vendor:       blocks.Disks[i].Vendor,
			Model:        blocks.Disks[i].Model,
			SerialNumber: blocks.Disks[i].SerialNumber,
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

	// cmd_out, err = helper_function_get_command_output("uname", []string{"-i"})
	// if err != nil {
	// 	return
	// }
	// h.HardwarePlatform = string(cmd_out)

	cmd_out, err = helper_function_get_command_output("uname", []string{"-o"})
	if err != nil {
		return
	}
	h.OperatingSystem = string(cmd_out)

	return
}

/************************************************************************/

type Host_device_payloads_information_data_t struct {
	Uname     host_device_payloads_information_data_uname_t
	PCI       host_device_payloads_information_data_pci_t
	DeviceID  host_device_payloads_information_data_device_id_t
	Blocks    host_device_payloads_information_data_block_t
	Memory    host_device_payloads_information_data_memory_t
	CPU       host_device_payloads_information_data_cpu_t
	Network   host_device_payloads_information_data_network_t
	Baseboard host_device_payloads_information_data_baseboard_t
	GPU       host_device_payloads_information_data_gpu_t
	BIOS      host_device_payloads_information_data_bios_t
}

func (h *Host_device_payloads_information_data_t) Get() (err error) {

	/* uname */
	err = h.Uname.get()
	if err != nil {
		return
	}
	/**/

	/* usb */
	// err = h.USB.get()
	// if err != nil {
	// 	return
	// }
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

	err = h.Baseboard.get()
	if err != nil {
		return
	}

	err = h.BIOS.get()
	if err != nil {
		return
	}

	err = h.Network.get()
	if err != nil {
		return
	}

	err = h.GPU.get()
	if err != nil {
		return
	}
	return
}
