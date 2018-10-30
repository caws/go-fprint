package main

import (
	"fmt"
	"os"
	"github.com/caws/go-fprint"
)

func main() {
	//This channel is needed so that you can have access to the messages
	//related to the device and fingerprint capture
	messages := make(chan string)

	//Since channels are naturally blocking calls in Go,
	//we create a go routine that 'listens' to whenever a
	//new string is added to the messages channel
	go func(msg string) {
		for {
			select {
			case msg := <-messages:
				fmt.Println("Channel Update: ", msg)
			default:
				//fmt.Println("No message received")
			}
		}
	}("Finished")

	dev := gofprint.GoOpenDevice(messages)

	if dev == nil{
		fmt.Println("Error while opening device")
		os.Exit(1)
	}

	fmt.Println("Opened device. It's now time to enroll your finger.")
	data := gofprint.GoEnroll(messages, dev)
	if data == nil{
		fmt.Println("Fingerprint not properly captured....")
		os.Exit(1)
	}

	//If you'd like to compare with fingerprint data in the original fp_print_data format
	fmt.Println("Comparing with original fp_print_data...", data)
	index := gofprint.GoVerify(messages, dev, data)
	fmt.Println("Index:", index)

	//If you'd like to compare with fingerprint data that is inside a []byte slice
	fprintByteSlice := gofprint.GoFprintDataToByteSlice(data)
	fmt.Println("Comparing with byte slice...")
	check:= gofprint.GoVerify(messages, dev, gofprint.GoByteSliceDataToFprint(fprintByteSlice ))

	fmt.Println("SCAN RESULT: ",check)

	gofprint.GoFreeData(data)
	gofprint.GoCloseDevice(dev)
}