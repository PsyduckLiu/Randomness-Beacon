package node

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"entropyNode/commitment"
	"entropyNode/config"
	"entropyNode/message"
	"entropyNode/signature"
	"entropyNode/util"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// initialize an entropy node
func StartEntropyNode(id int) {
	fmt.Printf("[Node%d] is running\n", id)

	// get specified curve
	config.ReadConfig()
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey(marshalledCurve)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Key message parse err:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve
	fmt.Printf("Curve is %v\n", curve.Params())

	// generate private key
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Generate private key err:%s", err))
	}
	fmt.Printf("===>My own key is: %v\n", privateKey)

	go WatchConfig(privateKey, id)
}

func WatchConfig(privateKey *ecdsa.PrivateKey, id int) {
	previousOutput := string(config.GetPreviousInput())
	fmt.Println("init output", previousOutput)

	// set config file
	viper.SetConfigFile("../config.yml")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		time.Sleep(1 * time.Second)
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		newOutput := string(config.GetPreviousInput())
		if previousOutput != newOutput {
			fmt.Println("output change")
			fmt.Println("newoutput:", newOutput)

			// calculate VRF result
			previousOutput = newOutput
			vrfResult := calVRF([]byte(newOutput), privateKey)
			vrfResultBinary := util.BytesToBinaryString(vrfResult)
			fmt.Printf("VRF result is:%v\n", util.BytesToBinaryString(vrfResult))
			fmt.Printf("VRF result last bit is:%v\n", vrfResultBinary[len(vrfResultBinary)-1:])

			// match VRF result with difficulty
			difficulty := config.GetDifficulty()
			vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
			if err != nil {
				panic(err)
			}
			if vrfResultTail == difficulty {
				fmt.Println("yes!!!!!!!!!!")
				timeCommitment := commitment.GenerateTimeCommitment()
				sendTCMsg(privateKey, vrfResult, int64(id), timeCommitment.String())
				fmt.Printf("Time commitment(now random number) is:%v\n", timeCommitment)
			}
		}
	})
}

// send time commitment message
func sendTCMsg(sk *ecdsa.PrivateKey, vrfResult []byte, id int64, tc string) {
	// new time commitment message
	// send time commitment message to origin nodes
	marshalledKey, err := x509.MarshalPKIXPublicKey(&sk.PublicKey)
	if err != nil {
		panic(fmt.Errorf("setup conf curve(marshalled Key) failed, err:%s", err))
	}
	tcMsg := &message.EntropyMessage{
		PublicKey:      marshalledKey,
		VRFResult:      vrfResult,
		TimeStamp:      time.Now().Unix(),
		ClientID:       id,
		TimeCommitment: tc,
	}

	// get consensus nodes' information
	nodeConfig := config.GetConsensusNode()
	for i := 0; i < len(nodeConfig); i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: util.EntropyPortByID(i)})
		if err != nil {
			fmt.Printf("dial tcp err:%s\n", err)
			continue
		}

		cMsg := message.CreateConMsg(message.MTCollect, tcMsg, sk, id)
		bs, err := json.Marshal(cMsg)
		if err != nil {
			panic(err)
		}

		// go WriteUDP(conn, rAddr, bs)
		go WriteTCP(conn, bs)
	}
}

func WriteTCP(conn *net.TCPConn, v []byte) {
	_, err := conn.Write(v)
	if err != nil {
		fmt.Printf("===>[ERROR]write to node err:%s\n", err)
		panic(err)
	}
	fmt.Println("Send request success!:=>")
}

// calculate VRF output
func calVRF(previousOutput []byte, sk *ecdsa.PrivateKey) []byte {
	// skD := sk.D
	vrfRes := signature.GenerateSig(previousOutput, sk)

	valid := signature.VerifySig(previousOutput, vrfRes, &sk.PublicKey)
	fmt.Printf("Verify result is:%v\n", valid)
	return vrfRes
}
