package node

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"entropyNode/commitment"
	"entropyNode/config"
	"entropyNode/message"
	"entropyNode/util"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/protocol"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type TestingHashable struct {
	data []byte
}

func (s TestingHashable) ToBeHashed() (protocol.HashID, []byte) {
	return "test", s.data
}

func randString() (b TestingHashable) {
	d := make([]byte, 100)
	_, err := rand.Read(d)
	if err != nil {
		panic(err)
	}

	fmt.Println("new random string is", d)
	return TestingHashable{d}
}

// initialize an entropy node
func StartEntropyNode(id int) {
	fmt.Printf("[Node%d] is running\n", id)

	// get specified curve
	config.ReadConfig()
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
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

	pk, sk := crypto.VrfKeygen()
	fmt.Println("public key is", pk)
	fmt.Println("secret key is", sk)

	var signal chan interface{}
	go WatchConfig(privateKey, sk, id, signal)

	s := <-signal
	fmt.Printf("===>[EXIT]Node[%d] exit because of:%s\n", id, s)
}

func WatchConfig(sk *ecdsa.PrivateKey, privateKey crypto.VrfPrivkey, id int, sig chan interface{}) {
	previousOutput := string(config.GetPreviousOutput())
	fmt.Println("init output", previousOutput)

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../output.yml")
	myViper.WatchConfig()
	myViper.OnConfigChange(func(e fsnotify.Event) {
		// lock file
		f, err := os.Open("../lock")
		if err != nil {
			panic(err)
		}
		// if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
			log.Println("add share lock in no block failed", err)
		}
		fmt.Println(time.Now())
		fmt.Println("Config Change")

		// config.ReadConfig()
		if err := myViper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		newOutput := string(config.GetPreviousOutput())
		if previousOutput != newOutput && newOutput != "" {
			fmt.Println("output change", newOutput)

			// calculate VRF result
			previousOutput = newOutput
			msg := randString()
			vrfResult, ok := privateKey.Prove(msg)
			if !ok {
				panic("Failed to construct VRF proof")
			}
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

				// generate commit
				groupLength := 2048
				c, h, rKSubOne, rK, a1, a2, a3, z := commitment.GenerateTimeCommitment(groupLength)
				fmt.Println("[TC]c is", c)
				fmt.Println("[TC]h is", h)
				fmt.Println("[TC]rKSubOne is", rKSubOne)
				fmt.Println("[TC]rK is", rK)
				fmt.Println("[TC]a1 is", a1)
				fmt.Println("[TC]a2 is", a2)
				fmt.Println("[TC]a3 is", a3)
				fmt.Println("[TC]z is", z)

				cMarshal, _ := c.MarshalJSON()
				hMarshal, _ := h.MarshalJSON()
				rKSubOneMarshal, _ := rKSubOne.MarshalJSON()
				rKMarshal, _ := rK.MarshalJSON()
				a1Marshal, _ := a1.MarshalJSON()
				a2Marshal, _ := a2.MarshalJSON()
				a3Marshal, _ := a3.MarshalJSON()
				zMarshal, _ := z.MarshalJSON()
				// timeCommitment := [4][]byte{cMarshal, hMarshal, rKSubOneMarshal, rKMarshal}
				sendVRFMsg(sk, privateKey, vrfResult, msg.data, int64(id), sig)
				time.Sleep(100 * time.Millisecond)
				sendTCMsg(sk, int64(id), cMarshal, hMarshal, rKSubOneMarshal, rKMarshal, a1Marshal, a2Marshal, a3Marshal, zMarshal, sig)
				fmt.Printf("Time commitment is:%v,%v,%v,%v\n", cMarshal, hMarshal, rKSubOneMarshal, rKMarshal)
			}
		}

		// unlock file
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			log.Println("unlock share lock failed", err)
		}
	})
}

func sendVRFMsg(sk *ecdsa.PrivateKey, privateKey crypto.VrfPrivkey, vrfResult crypto.VRFProof, msg []byte, id int64, sig chan interface{}) {
	// new time commitment message
	// send time commitment message to origin nodes
	vrfMsg := &message.EntropyVRFMessage{
		PublicKey: privateKey.Pubkey(),
		VRFResult: vrfResult,
		ClientID:  id,
		Msg:       msg,
	}

	fmt.Println(privateKey.Pubkey())
	fmt.Println(msg)
	fmt.Println(vrfResult)

	// get consensus nodes' information
	nodeConfig := config.GetConsensusNode()
	for i := 0; i < len(nodeConfig); i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: util.EntropyPortByID(i)})
		if err != nil {
			fmt.Println(time.Now())
			fmt.Printf("dial tcp err:%s\n", err)
			continue
		}

		cMsg := message.CreateConMsg(message.MTVRFVerify, vrfMsg, sk, id)
		bs, err := json.Marshal(cMsg)
		if err != nil {
			panic(err)
		}
		fmt.Println(len(bs))

		go WriteTCP(conn, bs)
	}
}

func sendTCMsg(sk *ecdsa.PrivateKey, id int64, cMarshal []byte, hMarshal []byte, rKSubOneMarshal []byte, rKMarshal []byte, a1Marshal []byte, a2Marshal []byte, a3Marshal []byte, zMarshal []byte, sig chan interface{}) {
	// new time commitment message
	// send time commitment message to origin nodes
	tcMsg := &message.EntropyTCMessage{
		ClientID:               id,
		TimeCommitmentC:        string(cMarshal),
		TimeCommitmentH:        string(hMarshal),
		TimeCommitmentrKSubOne: string(rKSubOneMarshal),
		TimeCommitmentrK:       string(rKMarshal),
		TimeCommitmentA1:       string(a1Marshal),
		TimeCommitmentA2:       string(a2Marshal),
		TimeCommitmentA3:       string(a3Marshal),
		TimeCommitmentZ:        string(zMarshal),
	}

	// get consensus nodes' information
	nodeConfig := config.GetConsensusNode()
	for i := 0; i < len(nodeConfig); i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: util.EntropyPortByID(i)})
		if err != nil {
			fmt.Println(time.Now())
			fmt.Printf("dial tcp err:%s\n", err)
			continue
		}

		cMsg := message.CreateConMsg(message.MTCollect, tcMsg, sk, id)
		bs, err := json.Marshal(cMsg)
		if err != nil {
			panic(err)
		}
		fmt.Println(len(bs))

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
