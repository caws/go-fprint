package main

import (
	"fmt"
	"os"
	"github.com/caws/go-fprint"
)

func main() {
	dev := gofprint.GoOpenDevice()

	if dev == nil{
		fmt.Println("Error while opening device")
		os.Exit(1)
	}

	fmt.Println("Opened device. \n It's now time to enroll your FIRST fingerprint.")
	data := gofprint.GoEnroll(dev)
	if data == nil{
		fmt.Println("Fingerprint not properly captured....")
		os.Exit(1)
	}

	fmt.Println("It's now time to enroll your SECOND fingerprint.")
	data2 := gofprint.GoEnroll(dev)
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
	gofprint.GoIdentifyFingerprints(dev, fingerprintArray)

	gofprint.GoFreeData(data)
	gofprint.GoCloseDevice(dev)
}