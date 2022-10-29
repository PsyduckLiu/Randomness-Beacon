package config

import (
	"consensusNode/signature"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

type Configurations struct {
	Running        bool
	Version        string
	PreviousOutput string
	EllipticCurve  string
	VRFType        int
	TCType         int
	FType          int
	Difficulty     int
	// ConsensusNodes NodesConfig
	Node0_ip string
	Node0_pk string
	Node1_ip string
	Node1_pk string
	Node2_ip string
	Node2_pk string
	Node3_ip string
	Node3_pk string
	Node4_ip string
	Node4_pk string
	Node5_ip string
	Node5_pk string
	Node6_ip string
	Node6_pk string
}

// type NodesConfig struct {
// 	Node0 NodeConfig
// 	Node1 NodeConfig
// 	Node2 NodeConfig
// 	Node3 NodeConfig
// 	Node4 NodeConfig
// 	Node5 NodeConfig
// 	Node6 NodeConfig
// }

type NodeConfig struct {
	Ip string
	Pk string
}

// config file Setup
// assign curve,previous output,VRF type,Time Commitment Type,F Function Type
func SetupConfig() {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		log.Println("add share lock in no block failed", err)
	}

	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	oldConfig := viper.AllSettings()
	fmt.Printf("All settings #1 %+v\n\n", oldConfig)

	// first time setup
	if !viper.GetBool("Running") {
		fmt.Println("First time Setup")

		// generate elliptic curve
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(fmt.Errorf("setup conf curve(private Key) failed, err:%s", err))
		}
		marshalledKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("setup conf curve(marshalled Key) failed, err:%s", err))
		}
		viper.Set("EllipticCurve", string(marshalledKey))

		// generate random init input
		message := []byte("asdkjhdk")
		randomNum := signature.Digest(message)
		viper.Set("PreviousOutput", string(randomNum))

		// generate random difficulty
		selectBigInt, _ := rand.Int(rand.Reader, big.NewInt(2))
		selectInt, err := strconv.Atoi(selectBigInt.String())
		if err != nil {
			panic(err)
		}
		viper.Set("Difficulty", selectInt)

		viper.Set("Running", true)

		// TODO: VRF type,Time Commitment Type,F Function Type

		oldConfig := viper.AllSettings()
		fmt.Printf("All settings #1 %+v\n\n", oldConfig)

		// write new settings
		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("setup conf failed, err:%s", err))
		}

		newnewConfig := viper.AllSettings()
		fmt.Printf("All settings #3 %+v\n\n", newnewConfig)

		fmt.Println("Finish time Setup")
	}

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}
}

// write new previousputput
func WriteOutput(output string) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		log.Println("add share lock in no block failed", err)
	}

	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// oldConfig := viper.AllSettings()
	// fmt.Printf("All settings #1 %+v\n\n", oldConfig)

	viper.Set("PreviousOutput", output)

	// newConfig := viper.AllSettings()
	// fmt.Printf("All settings #2 %+v\n\n", newConfig)
	// fmt.Println(viper.GetString("Node1_IP"))
	// fmt.Println(viper.GetString("Node1_PK"))
	// fmt.Println(viper.GetString("Node2_IP"))
	// fmt.Println(viper.GetString("Node2_PK"))
	// fmt.Println(viper.GetString("Node3_IP"))
	// fmt.Println(viper.GetString("Node3_PK"))
	// fmt.Println("start", time.Now())

	// write new settings
	if err := viper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
	ReadConfig()
	// if err := viper.WriteConfigAs("../output.yml"); err != nil {
	// 	panic(fmt.Errorf("setup conf failed, err:%s", err))
	// }

	// fmt.Println("end", time.Now())
	fmt.Println("Write output")
	// time.Sleep(1 * time.Second)
	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}

	// newnewConfig := viper.AllSettings()
	// fmt.Printf("All settings #3 %+v\n\n", newnewConfig)
	// fmt.Println(viper.GetString("ConsensusNodes.Node1.PK"))

}

// get difficulty from config
func GetDifficulty() int {
	var configuration = new(Configurations)

	// set config file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	return configuration.Difficulty
}

// get curve from config
func GetCurve() string {

	var configuration = new(Configurations)

	// set config file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	return configuration.EllipticCurve
}

// get previous input from config
func GetPreviousInput() string {
	var configuration = new(Configurations)

	// set config file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	return configuration.PreviousOutput
}

// write new id-ip-pk into config
func NewConsensusNode(id int64, ip string, pk string) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		log.Println("add share lock in no block failed", err)
	}

	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// handle new id-ip-pk
	// newNode := NodeConfig{
	// 	Ip: ip,
	// 	Pk: pk,
	// }
	switch id {
	// case 0:
	// 	viper.Set("ConsensusNodes.Node0", newNode)
	// case 1:
	// 	viper.Set("ConsensusNodes.Node1", newNode)
	// case 2:
	// 	viper.Set("ConsensusNodes.Node2", newNode)
	// case 3:
	// 	viper.Set("ConsensusNodes.Node3", newNode)
	// case 4:
	// 	viper.Set("ConsensusNodes.Node4", newNode)
	// case 5:
	// 	viper.Set("ConsensusNodes.Node5", newNode)
	// case 6:
	// 	viper.Set("ConsensusNodes.Node6", newNode)

	case 0:
		// viper.Set("ConsensusNodes.Node0.Ip", ip)
		// viper.Set("ConsensusNodes.Node0.Pk", pk)
		viper.Set("Node0_Ip", ip)
		viper.Set("Node0_Pk", pk)
	case 1:
		// viper.Set("ConsensusNodes.Node1.Ip", ip)
		// viper.Set("ConsensusNodes.Node1.Pk", pk)
		viper.Set("Node1_Ip", ip)
		viper.Set("Node1_Pk", pk)
	case 2:
		// viper.Set("ConsensusNodes.Node2.Ip", ip)
		// viper.Set("ConsensusNodes.Node2.Pk", pk)
		viper.Set("Node2_Ip", ip)
		viper.Set("Node2_Pk", pk)
	case 3:
		// viper.Set("ConsensusNodes.Node3.Ip", ip)
		// viper.Set("ConsensusNodes.Node3.Pk", pk)
		viper.Set("Node3_Ip", ip)
		viper.Set("Node3_Pk", pk)
	case 4:
		// viper.Set("ConsensusNodes.Node4.Ip", ip)
		// viper.Set("ConsensusNodes.Node4.Pk", pk)
		viper.Set("Node4_Ip", ip)
		viper.Set("Node4_Pk", pk)
	case 5:
		// viper.Set("ConsensusNodes.Node5.Ip", ip)
		// viper.Set("ConsensusNodes.Node5.Pk", pk)
		viper.Set("Node5_Ip", ip)
		viper.Set("Node5_Pk", pk)
	case 6:
		// viper.Set("ConsensusNodes.Node6.Ip", ip)
		// viper.Set("ConsensusNodes.Node6.Pk", pk)
		viper.Set("Node6_Ip", ip)
		viper.Set("Node6_Pk", pk)
	}

	// write new settings
	if err := viper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
	fmt.Println("new consensus node")

	oldConfig := viper.AllSettings()
	fmt.Printf("All settings #1 %+v\n\n", oldConfig)
	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}
}

func RemoveConsensusNode(id int64) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		log.Println("add share lock in no block failed", err)
	}

	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// handle new id-ip-pk
	switch id {
	case 0:
		if viper.GetString("consensusnodes.node0.pk") != "0" {
			viper.Set("consensusnodes.node0.pk", "0")
			viper.Set("consensusnodes.node0.ip", "0")
		}
	case 1:
		if viper.GetString("consensusnodes.node1.pk") != "0" {
			viper.Set("consensusnodes.node1.pk", "0")
			viper.Set("consensusnodes.node1.ip", "0")
		}
	case 2:
		if viper.GetString("consensusnodes.node2.pk") != "0" {
			viper.Set("consensusnodes.node2.pk", "0")
			viper.Set("consensusnodes.node2.ip", "0")
		}
	case 3:
		if viper.GetString("consensusnodes.node3.pk") != "0" {
			viper.Set("consensusnodes.node3.pk", "0")
			viper.Set("consensusnodes.node3.ip", "0")
		}
	case 4:
		if viper.GetString("consensusnodes.node4.pk") != "0" {
			viper.Set("consensusnodes.node4.pk", "0")
			viper.Set("consensusnodes.node4.ip", "0")
		}
	case 5:
		if viper.GetString("consensusnodes.node5.pk") != "0" {
			viper.Set("consensusnodes.node5.pk", "0")
			viper.Set("consensusnodes.node5.ip", "0")
		}
	case 6:
		if viper.GetString("consensusnodes.node6.pk") != "0" {
			viper.Set("consensusnodes.node6.pk", "0")
			viper.Set("consensusnodes.node6.ip", "0")
		}
	}

	// write new settings
	if err := viper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
	fmt.Println("remove consensus node")

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}
}

// get consensus nodes from config
func GetConsensusNode() []NodeConfig {
	var nodeConfig []NodeConfig
	var configuration = new(Configurations)

	// set config file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	for i := 0; i < 7; i++ {
		var node NodeConfig
		switch i {
		case 0:
			// node.Ip = configuration.ConsensusNodes.Node0.Ip
			// node.Pk = configuration.ConsensusNodes.Node0.Pk
			node.Ip = configuration.Node0_ip
			node.Pk = configuration.Node0_pk
			nodeConfig = append(nodeConfig, node)
		case 1:
			// node.Ip = configuration.ConsensusNodes.Node1.Ip
			// node.Pk = configuration.ConsensusNodes.Node1.Pk
			node.Ip = configuration.Node1_ip
			node.Pk = configuration.Node1_pk
			nodeConfig = append(nodeConfig, node)
		case 2:
			// node.Ip = configuration.ConsensusNodes.Node2.Ip
			// node.Pk = configuration.ConsensusNodes.Node2.Pk
			node.Ip = configuration.Node2_ip
			node.Pk = configuration.Node2_pk
			nodeConfig = append(nodeConfig, node)
		case 3:
			// node.Ip = configuration.ConsensusNodes.Node3.Ip
			// node.Pk = configuration.ConsensusNodes.Node3.Pk
			node.Ip = configuration.Node3_ip
			node.Pk = configuration.Node3_pk
			nodeConfig = append(nodeConfig, node)
		case 4:
			// node.Ip = configuration.ConsensusNodes.Node4.Ip
			// node.Pk = configuration.ConsensusNodes.Node4.Pk
			node.Ip = configuration.Node4_ip
			node.Pk = configuration.Node4_pk
			nodeConfig = append(nodeConfig, node)
		case 5:
			// node.Ip = configuration.ConsensusNodes.Node5.Ip
			// node.Pk = configuration.ConsensusNodes.Node5.Pk
			node.Ip = configuration.Node5_ip
			node.Pk = configuration.Node5_pk
			nodeConfig = append(nodeConfig, node)
		case 6:
			// node.Ip = configuration.ConsensusNodes.Node6.Ip
			// node.Pk = configuration.ConsensusNodes.Node6.Pk
			node.Ip = configuration.Node6_ip
			node.Pk = configuration.Node6_pk
			nodeConfig = append(nodeConfig, node)
		}
	}

	return nodeConfig
}

// get messages from config file
func ReadConfig() {
	var configuration = new(Configurations)

	// set fonfig file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	fmt.Printf("Reading using model:\n")
	fmt.Printf("Running:%v\n", configuration.Running)
	fmt.Printf("Version:%s\n", configuration.Version)
	fmt.Printf("PreviousOutput:%s\n", configuration.PreviousOutput)
	fmt.Printf("EllipticCurve:%v\n", configuration.EllipticCurve)
	fmt.Printf("Consensusnodes:\n")
	fmt.Printf("Node[0]'s ip is %s\n", configuration.Node0_ip)
	fmt.Printf("Node[0]'s pk is %s\n", configuration.Node0_pk)
	fmt.Printf("Node[1]'s ip is %s\n", configuration.Node1_ip)
	fmt.Printf("Node[1]'s pk is %s\n", configuration.Node1_pk)
	fmt.Printf("Node[2]'s ip is %s\n", configuration.Node2_ip)
	fmt.Printf("Node[2]'s pk is %s\n", configuration.Node2_pk)
	fmt.Printf("Node[3]'s ip is %s\n", configuration.Node3_ip)
	fmt.Printf("Node[3]'s pk is %s\n", configuration.Node3_pk)
	fmt.Printf("Node[4]'s ip is %s\n", configuration.Node4_ip)
	fmt.Printf("Node[4]'s pk is %s\n", configuration.Node4_pk)
	fmt.Printf("Node[5]'s ip is %s\n", configuration.Node5_ip)
	fmt.Printf("Node[5]'s pk is %s\n", configuration.Node5_pk)
	fmt.Printf("Node[6]'s ip is %s\n", configuration.Node6_ip)
	fmt.Printf("Node[6]'s pk is %s\n", configuration.Node6_pk)
	// fmt.Printf("Node[0]'s ip is %s\n", configuration.ConsensusNodes.Node0.Ip)
	// fmt.Printf("Node[0]'s pk is %s\n", configuration.ConsensusNodes.Node0.Pk)
	// fmt.Printf("Node[1]'s ip is %s\n", configuration.ConsensusNodes.Node1.Ip)
	// fmt.Printf("Node[1]'s pk is %s\n", configuration.ConsensusNodes.Node1.Pk)
	// fmt.Printf("Node[2]'s ip is %s\n", configuration.ConsensusNodes.Node2.Ip)
	// fmt.Printf("Node[2]'s pk is %s\n", configuration.ConsensusNodes.Node2.Pk)
	// fmt.Printf("Node[3]'s ip is %s\n", configuration.ConsensusNodes.Node3.Ip)
	// fmt.Printf("Node[3]'s pk is %s\n", configuration.ConsensusNodes.Node3.Pk)
	// fmt.Printf("Node[4]'s ip is %s\n", configuration.ConsensusNodes.Node4.Ip)
	// fmt.Printf("Node[4]'s pk is %s\n", configuration.ConsensusNodes.Node4.Pk)
	// fmt.Printf("Node[5]'s ip is %s\n", configuration.ConsensusNodes.Node5.Ip)
	// fmt.Printf("Node[5]'s pk is %s\n", configuration.ConsensusNodes.Node5.Pk)
	// fmt.Printf("Node[6]'s ip is %s\n", configuration.ConsensusNodes.Node6.Ip)
	// fmt.Printf("Node[6]'s pk is %s\n", configuration.ConsensusNodes.Node6.Pk)
	fmt.Printf("Difficulty:%d\n", configuration.Difficulty)
	fmt.Printf("VRFType:%d\n", configuration.VRFType)
	fmt.Printf("TCType:%d\n", configuration.TCType)
	fmt.Printf("FType:%d\n", configuration.FType)
}
