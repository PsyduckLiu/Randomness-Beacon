package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Reveal tc and compute F
func (s *StateEngine) RevealTC(msg *message.ConMessage) (err error) {
	if s.RevealNum >= 2*util.MaxFaultyNode {
		fmt.Println("late!")
		return
	}

	fmt.Printf("\n===>[Reveal]Current Reveal Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from RevealTC]Parse elliptic curve error:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve

	// unmarshal public key
	x, y := elliptic.Unmarshal(curve, []byte(nodeConfig[msg.From].Pk))
	newPublicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// verify signature
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		panic(fmt.Errorf("===>[ERROR from RevealTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	result := new(string)
	if err := json.Unmarshal(msg.Payload, result); err != nil {
		panic(fmt.Errorf("===>[ERROR from RevealTC]Invalid[%s] Reveal message[%s]", err, msg))
	}

	// new Reveal message
	// send Reveal message to primary node
	if *result == s.Result.String() {
		s.RevealNum++
		// if s.RevealNum == util.TotalNodeNum-1 {
		if s.RevealNum == 2*util.MaxFaultyNode {
			fmt.Println("===>[Reveal]new Reveal is", s.Result)
			s.GlobalTimer.tack()
			util.WriteResult(s.Result.String())
			s.stage = CommitFromEntropy
			s.quit <- true

			RevealNum++
			if RevealNum >= 2 {
				currentTime := time.Now()
				totalTimeArray = append(totalTimeArray, float64(currentTime.Sub(lastTime).Seconds()))
				verifyTimeArray = append(verifyTimeArray, s.VerifyTime)
			}
			lastTime = time.Now()
			if RevealNum == 12 {
				writeTotalTimeFile(totalTimeArray)
				writeVerifyTimeFile(verifyTimeArray)
			}
			if RevealNum == 22 {
				writeTotalTimeFile(totalTimeArray)
				writeVerifyTimeFile(verifyTimeArray)
			}
			if RevealNum == 32 {
				writeTotalTimeFile(totalTimeArray)
				writeVerifyTimeFile(verifyTimeArray)
			}
			if RevealNum == 42 {
				writeTotalTimeFile(totalTimeArray)
				writeVerifyTimeFile(verifyTimeArray)
			}
			if RevealNum == 52 {
				writeTotalTimeFile(totalTimeArray)
				writeVerifyTimeFile(verifyTimeArray)
			}

			time.Sleep(5 * time.Second)
			config.WriteOutput(s.Result.String())
			fmt.Println("\n===>[Reveal]Reveal time is", time.Now())
			// fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
			// buf := make([]byte, 64*1024)
			// runtime.Stack(buf, true)
			// fmt.Println("the detail of goroutines: ", string(buf))
		}
	}

	return
}

func writeTotalTimeFile(totalTimeArray []float64) {
	// create a file
	file, err := os.Create("time.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data []string
	for index, _ := range totalTimeArray {
		data = append(data, string(strconv.FormatFloat(totalTimeArray[index], 'f', 5, 32)))
	}

	// write all rows at once
	writer.Write(data)
}

func writeVerifyTimeFile(verifyTimeArray []float64) {
	// create a file
	file, err := os.Create("verifyTime.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data []string
	for index, _ := range verifyTimeArray {
		data = append(data, string(strconv.FormatFloat(verifyTimeArray[index], 'f', 5, 32)))
	}

	// write all rows at once
	writer.Write(data)
}
