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

// output tc and compute F
func (s *StateEngine) outputTC(msg *message.ConMessage) (err error) {
	if s.OutputNum >= 2*util.MaxFaultyNode {
		fmt.Println("late!")
		return
	}

	fmt.Printf("\n===>[Output]Current Output Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from outputTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from outputTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	result := new(string)
	if err := json.Unmarshal(msg.Payload, result); err != nil {
		panic(fmt.Errorf("===>[ERROR from outputTC]Invalid[%s] Output message[%s]", err, msg))
	}

	// new output message
	// send output message to primary node
	if *result == s.Result.String() {
		s.OutputNum++
		// if s.OutputNum == util.TotalNodeNum-1 {
		if s.OutputNum == 2*util.MaxFaultyNode {
			fmt.Println("===>[Output]new output is", s.Result)
			s.GlobalTimer.tack()
			util.WriteResult(s.Result.String())
			s.stage = Collect
			s.quit <- true

			time.Sleep(5 * time.Second)
			config.WriteOutput(s.Result.String())

			outputNum++
			if outputNum > 1 {
				currentTime := time.Now()
				timeArray = append(timeArray, float64(currentTime.Sub(lastTime).Seconds()))
			}
			lastTime = time.Now()
			if outputNum == 11 {
				writeDataFile(timeArray)
			}

			// fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
			// buf := make([]byte, 64*1024)
			// runtime.Stack(buf, true)
			// fmt.Println("the detail of goroutines: ", string(buf))
		}
	}

	return
}

func writeDataFile(timeArray []float64) {
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
	for index, _ := range timeArray {
		data = append(data, string(strconv.FormatFloat(timeArray[index], 'f', 5, 32)))
	}

	// write all rows at once
	writer.Write(data)
}
