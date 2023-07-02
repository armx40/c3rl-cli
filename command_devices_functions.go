package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	// "github.com/google/gousb"
	"github.com/fatih/color"
	"github.com/google/gousb"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/block"
	"golang.org/x/exp/slices"
)

var csv_file *os.File
var csv_writer *csv.Writer
var csv_header_set bool
var csv_delimiter string

var sqlite3_db *sql.DB
var sqlite3_db_statement *sql.Stmt
var sqlite3_db_tx *sql.Tx

var command_devices_functions_device_symmetric_key []byte
var command_devices_functions_valid_products = []string{"ATNode-R"}

func command_devices_functions_find_c3rl_device(print_out bool) (devices []*gousb.Device, err error) {

	/* check if user has permissions */

	err = helper_functions_is_user_dialout()
	if err != nil {
		print_text := color.New(color.FgRed)
		print_text.Println("User doesn't have permissions to access USB devices.")
		err = nil
		return
	}
	/**/

	ctx := gousb.NewContext()
	defer ctx.Close()

	// Debugging can be turned on; this shows some of the inner workings of the libusb package.
	ctx.Debug(*debug)

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Vendor == 0 && desc.Product == 1 {
			return true
		}
		return false
	})

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	if err != nil {
		return
	}

	for _, dev := range devs {
		manu, errd := dev.Manufacturer()
		if errd != nil {
			err = errd
			return
		}

		prod, errd := dev.Product()
		if errd != nil {
			err = errd
			return
		}

		serial, errd := dev.SerialNumber()
		if errd != nil {
			err = errd
			return
		}

		if manu == "c3rl opc pvt ltd" {
			/* check product */
			if slices.Contains(command_devices_functions_valid_products, prod) {
				/* check serial */
				devices = append(devices, dev)
				if print_out {
					fmt.Printf("%s [%s]\n", prod, serial)
				}

			}

		}
		time.Sleep(2 * time.Second)

		intf, done, err := dev.DefaultInterface()
		if err != nil {
			log.Fatalf("%s.DefaultInterface(): %v", dev, err)
		}
		log.Println(intf)
		defer done()
		ep, err := intf.OutEndpoint(1)
		if err != nil {
			log.Fatalf("%s.OutEndpoint(0x82): %v", intf, err)
		}

		epin, err := intf.InEndpoint(0x82)
		if err != nil {
			log.Fatalf("%s.OutEndpoint(0x82): %v", intf, err)
		}

		data := make([]byte, 5)
		for i := range data {
			data[i] = byte(i + 1)
		}

		numBytes, err := ep.Write(data)
		if numBytes != 5 {
			log.Fatalf("%s.Read([5]): only %d bytes read, returned error is %v", ep, numBytes, err)
		}
		log.Println(data)

		time.Sleep(1 * time.Second)
		log.Println("reading")

		numBytes, err = epin.Read(data)
		if numBytes != 5 {
			log.Fatalf("%s.Read([5]): only %d bytes read, returned error is %v", ep, numBytes, err)
		}
		log.Println(data)

		// log.Println("writing")
		// // Write data to the USB device.
		// numBytes, err := ep.Write(data)
		// if numBytes != 5 {
		// 	log.Fatalf("%s.Write([5]): only %d bytes written, returned error is %v", ep, numBytes, err)
		// }

		// time.Sleep(1 * time.Second)
		// log.Println("writing")

		// // Write data to the USB device.
		// numBytes, err = ep.Write(data)
		// if numBytes != 5 {
		// 	log.Fatalf("%s.Write([5]): only %d bytes written, returned error is %v", ep, numBytes, err)
		// }

	}

	return
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

func command_devices_functions_get_all_log_files(folder string) (dirs []fs.DirEntry, err error) {
	dirs_, err := os.ReadDir(folder)

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

func command_devices_functions_get_all_log_files_sorted(folder string, start_index int, num int) (files []fs.FileInfo, err error) {

	/* process num */

	all_log_files, err := command_devices_functions_get_all_log_files(folder)
	if err != nil {
		return
	}

	if num > len(all_log_files) {
		num = len(all_log_files)
	}

	curr_index := start_index

	for len(files) < num {

		filename := filepath.Join(folder, fmt.Sprintf("data_log_%d.data", curr_index))
		curr_index = (curr_index - 1) % 400000
		file_info, err := os.Stat(filename)
		if err != nil {
			continue
		}
		files = append(files, file_info)
	}

	return
}

func command_devices_functions_request_device_symmetric_key(device_uid string) (key []byte, err error) {

	if len(command_devices_functions_device_symmetric_key) > 0 {
		key = command_devices_functions_device_symmetric_key
		err = nil
		return
	}

	/* get public key bytes */
	public_key := crypto_ecc_get_session_key().PublicKey
	public_key_x := public_key.X.Bytes()
	public_key_y := public_key.Y.Bytes()
	/**/

	request_data := make(map[string]string)
	request_data["x"] = hex.EncodeToString(public_key_x)
	request_data["y"] = hex.EncodeToString(public_key_y)
	request_data["uid"] = device_uid

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	// params := make(map[string]string)
	// params["g"] = "lgn"

	headers := make(map[string]string)
	headers["Authorization"] = auth_data.Token

	resp, err := network_request(API_HOST+"devices?g=dsk", nil, headers, request_data)

	if err != nil {
		return
	}

	var response generalPayloadV2

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if response.Status != "success" {
		err = fmt.Errorf("failed to get device symmetric key")
		return
	}

	/* decode data */

	marshaled_bytes, err := json.Marshal(response.Payload)
	if err != nil {
		return
	}

	type data_ struct {
		ecdhPayload
		Key string `json:"k"`
	}

	var ecdh_data data_

	err = json.Unmarshal(marshaled_bytes, &ecdh_data)
	if err != nil {
		return
	}

	/**/
	/* get x and y bytes */
	x, y, err := ecdh_data.get_bytes()
	if err != nil {
		return
	}
	/**/
	/* calculate secret */
	secret, err := crypto_ecdh_perform_ecdh(x, y)
	if err != nil {
		return
	}
	/**/

	/* get key bytes */
	key_bytes, err := hex.DecodeString(ecdh_data.Key)
	if err != nil {
		return
	}
	/**/

	/* decrypt key */
	key, err = crypto_aes_decrypt(key_bytes[16:], secret[:16], key_bytes[:16])
	if err != nil {
		return
	}

	command_devices_functions_device_symmetric_key = key
	return
}
