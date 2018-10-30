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

	fmt.Println("Opened device. \n It's now time to enroll your FIRST fingerprint.")
	data := gofprint.GoEnroll(messages, dev)
	if data == nil{
		fmt.Println("Fingerprint not properly captured....")
		os.Exit(1)
	}

	fmt.Println("It's now time to enroll your SECOND fingerprint.")
	data2 := gofprint.GoEnroll(messages, dev)
	if data2 == nil{
		fmt.Println("Fingerprint not properly captured....")
		os.Exit(1)
	}

	//If you'd like to check one's fingerprint against many fingerprints in a [][]byte slice

	//If you'd like to compare with a bunch of fingerprints in a []byte slice
	myFingerprint1:= gofprint.GoFprintDataToByteSlice(data)
	myFingerprint2:= gofprint.GoFprintDataToByteSlice(data2)

	//You have to compose a bidimensional byte slice...
	var fingerprintArray [][]byte

	//...and append your []byte slices containing fingerprint data to it...
	fingerprintArray = append(fingerprintArray, myFingerprint1)
	fingerprintArray = append(fingerprintArray, myFingerprint2)

	//...and pass it as a parameter
	index := gofprint.GoIdentifyFingerprints(messages, dev, fingerprintArray)

	fmt.Println("Index matched in the slice: ", index)

	gofprint.GoFreeData(data)
	gofprint.GoCloseDevice(dev)
}