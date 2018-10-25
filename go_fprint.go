package gofprint
/*
#include <stdio.h>
#include <stdlib.h>
#include <libfprint/fprint.h>
#include <unistd.h>
#cgo LDFLAGS: -lfprint
*/
import "C"
import (
"fmt"
"os"
"unsafe"
)

func GoDiscoverDevice(discovered_devs  **C.struct_fp_dscv_dev) *C.struct_fp_dscv_dev {
	ddev := discovered_devs
	if ddev == nil{
		os.Exit(1)
	}

	//The two lines below discover the name of the device
	//drv := C.fp_dscv_dev_get_driver(*ddev)
	//fmt.Println(fmt.Sprintf("Found device claimed by %s driver", test))

	return *ddev
}

//The following method frees the data generated during the enroll process
func GoFreeData(data *C.struct_fp_print_data){
	C.fp_print_data_free(data)
}

//The following method closes the device
func GoCloseDevice(dev *C.struct_fp_dev){
	C.fp_dev_close(dev)
}

//The following function deinitialises libfprint
func GoFpExit(){
	C.fp_exit()
}

//The following method enrolls one's fingerprint and returns the fingerprint data
func GoEnroll(dev *C.struct_fp_dev ) *C.struct_fp_print_data  {
	C.fp_dev_get_nr_enroll_stages(dev)
	var enrolled_print *C.struct_fp_print_data

	for i := 0; i < 5; i++ {
		r := C.fp_enroll_finger(dev, &enrolled_print)
		if r < 0 {
			fmt.Println("Enroll failed with error ", r)
			os.Exit(9)
		}
		switch r {
		case C.FP_ENROLL_COMPLETE:
			fmt.Println("Enroll complete!")
			break
		case C.FP_ENROLL_FAIL:
			fmt.Println("Enroll failed, something wen't wrong :(")
			break
		case C.FP_ENROLL_PASS:
			fmt.Println("Enroll stage passed. Yay!")
			break
		case C.FP_ENROLL_RETRY:
			fmt.Println("Didn't quite catch that. Please try again")
			break
		case C.FP_ENROLL_RETRY_TOO_SHORT:
			fmt.Println("Your swipe was too short, please try again.")
			break
		case C.FP_ENROLL_RETRY_CENTER_FINGER:
			fmt.Println("Didn't catch that, please center your finger on the sensor and try again")
			break
		case C.FP_ENROLL_RETRY_REMOVE_FINGER:
			fmt.Println("Scan failed, please remove your finger and then try again")
			break
		}
	}
	return enrolled_print
}

//The following method checks one's fingerprint against fingerprint data passed as parameter
//and returns an integer.
// If the integer returned is -2, the verification failed.
// If the integer returned is -1, the fingerprint tested didn't match any of the fingerprints passed in gallery.
func GoVerify(dev *C.struct_fp_dev, data *C.struct_fp_print_data) int {

	for {
		fmt.Println("Scan your finger now.")
		r := C.fp_verify_finger(dev, data)
		if r < 0 {
			fmt.Println("verification failed with error :", r)
			return -2
		}
		switch r {
		case C.FP_VERIFY_NO_MATCH:
			fmt.Println("NO MATCH")
			return -1
		case C.FP_VERIFY_MATCH:
			fmt.Println("MATCH")
			return 1
		case C.FP_VERIFY_RETRY:
			fmt.Println("Scan didn't quite work. Please try again")
			return 0
		case C.FP_VERIFY_RETRY_TOO_SHORT:
			fmt.Println("Swipe was too short, please try again.")
			return 0
		case C.FP_VERIFY_RETRY_CENTER_FINGER:
			fmt.Println("Please center your finger on the sensor and try again.")
			return 0
		case C.FP_VERIFY_RETRY_REMOVE_FINGER:
			fmt.Println("Please remove finger from the sensor and try again.")
			return 0
		}
	}
	return -1
}

//The following method checks one's fingerprint against an array of fingerprint data passed as parameter
//and returns an integer.
// If the integer returned is -2, the verification failed.
// If the integer returned is -1, the fingerprint tested didn't match any of the fingerprints passed in gallery.
func GoIdentify(dev *C.struct_fp_dev, gallery []*C.struct_fp_print_data) int {

	var match_offset C.size_t
	match_offset = 0

	for {
		fmt.Println("Scan your finger now.")
		r := C.fp_identify_finger(dev, &gallery[0], &match_offset)
		if r < 0 {
			fmt.Println("verification failed with error :", r)
			return -2
		}
		switch r {
		case C.FP_VERIFY_NO_MATCH:
			fmt.Println("NO MATCH")
			return -1
		case C.FP_VERIFY_MATCH:
			fmt.Println("MATCH")
			fmt.Println("Index in the slice:", match_offset)
			return (int)(match_offset)
		case C.FP_VERIFY_RETRY:
			fmt.Println("Scan didn't quite work. Please try again")
			//return 0
		case C.FP_VERIFY_RETRY_TOO_SHORT:
			fmt.Println("Swipe was too short, please try again.")
			//return 0
		case C.FP_VERIFY_RETRY_CENTER_FINGER:
			fmt.Println("Please center your finger on the sensor and try again.")
			//return 0
		case C.FP_VERIFY_RETRY_REMOVE_FINGER:
			fmt.Println("Please remove finger from the sensor and try again.")
			//return 0
		}
	}
	return -1
}

//The following method opens the device
func GoOpenDevice() *C.struct_fp_dev {
	i := C.fp_init()

	if i < 0 {
		fmt.Println("Failed to initialize libfprint")
	}

	discovered_devs := C.fp_discover_devs()

	if discovered_devs == nil {
		fmt.Println("Could not discover devices")
		return nil
	}

	ddev := GoDiscoverDevice(discovered_devs)

	if ddev == nil {
		fmt.Println("No devices detected.")
		return nil
	}

	dev := C.fp_dev_open(ddev)
	C.fp_dscv_devs_free(discovered_devs)
	if dev == nil {
		fmt.Println("Could not open device.")
		return nil
	}

	return dev
}

//The following method converts fingerprint data to Go's byte slice
func GoFprintDataToByteSlice(data *C.struct_fp_print_data) []byte {
	var ret *C.uchar

	//Converting fprint_data_struct to file
	//size is of type size_t
	size := C.fp_print_data_get_data(data, &ret)

	//Converting size_t to C.int
	final_fprint_size := C.int(size)

	//Converting ret (fprint_file) to go []byte
	byteSlice := C.GoBytes(unsafe.Pointer(ret), final_fprint_size)

	return byteSlice
}

//The following method converts Go's byte slice to fingerprint data
func GoByteSliceDataToFprint(byteSlice []byte) *C.struct_fp_print_data {

	//This way I don't have to receive the size_t
	sizeSlice := cap(byteSlice)
	x := (*C.size_t)(unsafe.Pointer(&sizeSlice))

	//Converts mySlice to array of unsigned char
	cstr := (*C.uchar)(unsafe.Pointer(&byteSlice[0]))
	sliced_data := C.fp_print_data_from_data(cstr, *x)
	return sliced_data
}

//The following method converts Go's byte bidimensional slice to a struct_fp_print_data array and sends it
//as a parameter to be identified
func GoIdentifyFingerprints(dev *C.struct_fp_dev, byteSlices [][]byte) int {
	var print_gallery []*C.struct_fp_print_data
	for _, byteSlice := range byteSlices{
		print_gallery = append(print_gallery, GoByteSliceDataToFprint(byteSlice))
	}
	print_gallery = append(print_gallery, nil)

	index := GoIdentify(dev, print_gallery)
	fmt.Println("Index matched:", index)
	return index
}