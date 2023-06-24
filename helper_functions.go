package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"golang.org/x/exp/slices"
)

func helper_function_set_bit(number_to_set uint64, bit_no uint8, bit_value uint8) uint64 {
	if bit_value == 1 {
		number_to_set |= 1 << bit_no
		return number_to_set
	}

	number_to_set &= ^(1 << bit_no)
	return number_to_set

}

func helper_function_set_byte(number_to_set uint64, byte_loc uint8, byte_ uint8) uint64 {

	for i := 0; i < 8; i++ {
		bit_val := byte_ >> i & 0x01
		number_to_set = helper_function_set_bit(number_to_set, byte_loc+uint8(i), bit_val)
	}

	return number_to_set
}

func helper_function_set_bits(number_to_set uint64, bit_loc uint8, bits uint64, bits_count uint8) uint64 {

	for i := 0; i < int(bits_count); i++ {
		bit_val := bits >> i & 0x01
		number_to_set = helper_function_set_bit(number_to_set, bit_loc+uint8(i), uint8(bit_val))
	}

	return number_to_set
}

func helper_function_set_bytes(number_to_set uint64, byte_loc uint8, number_to_save uint64, bytes_size uint8) uint64 {

	for i := 0; i < int(bytes_size); i++ {
		number_to_set = helper_function_set_byte(number_to_set, byte_loc+uint8(i*8), uint8((number_to_save>>i)&0xff))
	}

	return number_to_set
}

func helper_function_find_bin(binary string) (string, error) {

	/* thanks to https://github.com/xxr3376/golspci */

	locations := []string{"/sbin", "/usr/sbin", "/usr/local/sbin", "/usr/bin", "/bin"}

	for _, path := range locations {
		lookup := path + "/" + binary
		fileInfo, err := os.Stat(path + "/" + binary)

		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			return lookup, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Unable to find the '%v' binary", binary))
}

func helper_function_get_command_output(command string, args []string) (stdout []byte, err error) {
	bin, err := helper_function_find_bin(command)
	if err != nil {
		return
	}

	cmd := exec.Command(bin, args...)

	// stdout, err = cmd.Output()
	// if err != nil {
	// 	return
	// }

	out := bytes.Buffer{}
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	stdout = out.Bytes()
	return
}

func helper_function_get_user_groups() (groups []string, err error) {
	curr_user, err := user.Current()
	if err != nil {
		return
	}

	groups_, err := curr_user.GroupIds()
	if err != nil {
		return
	}

	for i := range groups_ {
		group, errd := user.LookupGroupId(groups_[i])
		if errd != nil {
			return
		}

		groups = append(groups, group.Name)

	}

	return

}

func helper_function_is_user_root() (is_root bool, err error) {
	curr_user, err := user.Current()
	if err != nil {
		return
	}

	if curr_user.Gid == "0" && curr_user.Uid == "0" {
		is_root = true
	}

	return
}

func helper_function_is_user_dialout() (err error) {

	is_root, err := helper_function_is_user_root()
	if err != nil {
		return
	}

	if is_root {
		/* root has all access */
		return
	}

	/* if not root then check if user has access to dialout */
	groups, err := helper_function_get_user_groups()
	if err != nil {
		return
	}

	if slices.Contains(groups, "dialout") {
		return
	}
	err = fmt.Errorf("user dont have usb permissions")
	return
}
