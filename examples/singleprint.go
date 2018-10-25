package main

import (
	"fmt"
	"os"
	"gofprint"
)

func main() {
	dev := gofprint.GoOpenDevice()

	if dev == nil{
		fmt.Println("Error while opening device")
		os.Exit(1)
	}

	fmt.Println("Opened device. It's now time to enroll your finger.")
	data := gofprint.GoEnroll(dev)
	if data == nil{
		fmt.Println("Houston....")
		gofprint.GoCloseDevice(dev)
	}

	//If you'd like to compare with fingerprint data in the original fp_print_data format
	fmt.Println("Comparing with original fp_print_data...", data)
	index := gofprint.GoVerify(dev, data)
	fmt.Println("Index:", index)

	//If you'd like to compare with fingerprint data that is inside a []byte slice
	fprintByteSlice := gofprint.GoFprintDataToByteSlice(data)
	fmt.Println("Comparing with byte slice...")
	check:= gofprint.GoVerify(dev, gofprint.GoByteSliceDataToFprint(fprintByteSlice ))

	fmt.Println("SCAN RESULT: ",check)

	gofprint.GoFreeData(data)
	gofprint.GoCloseDevice(dev)
}