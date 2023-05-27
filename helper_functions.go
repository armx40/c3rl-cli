package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

	locations := []string{"/sbin", "/usr/sbin", "/usr/local/sbin", "/usr/bin"}

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
