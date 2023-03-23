package main

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
