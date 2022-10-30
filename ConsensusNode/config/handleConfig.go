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
	Running        bool   `mapstructure:"running"`
	Version        string `mapstructure:"version"`
	PreviousOutput string `mapstructure:"previousOutput"`
	EllipticCurve  string `mapstructure:"ellipticCurve"`
	VRFType        int    `mapstructure:"vrfType"`
	TCType         int    `mapstructure:"tcType"`
	FType          int    `mapstructure:"fType "`
	Difficulty     int    `mapstructure:"difficulty"`
	Node0_ip       string `mapstructure:"node0_ip"`
	Node0_pk       string `mapstructure:"node0_pk"`
	Node1_ip       string `mapstructure:"node1_ip"`
	Node1_pk       string `mapstructure:"node1_pk"`
	Node2_ip       string `mapstructure:"node2_ip"`
	Node2_pk       string `mapstructure:"node2_pk"`
	Node3_ip       string `mapstructure:"node3_ip"`
	Node3_pk       string `mapstructure:"node3_pk"`
	Node4_ip       string `mapstructure:"node4_ip"`
	Node4_pk       string `mapstructure:"node4_pk"`
	Node5_ip       string `mapstructure:"node5_ip"`
	Node5_pk       string `mapstructure:"node5_pk"`
	Node6_ip       string `mapstructure:"node6_ip"`
	Node6_pk       string `mapstructure:"node6_pk"`
}

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

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// oldConfig := myViper.AllSettings()
	// fmt.Printf("All settings #1 %+v\n\n", oldConfig)

	// first time setup
	if !myViper.GetBool("Running") {
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
		myViper.Set("EllipticCurve", string(marshalledKey))

		// generate random init input
		message := []byte("asdkjhdk")
		randomNum := signature.Digest(message)
		myViper.Set("PreviousOutput", string(randomNum))

		// generate random difficulty
		selectBigInt, _ := rand.Int(rand.Reader, big.NewInt(2))
		selectInt, err := strconv.Atoi(selectBigInt.String())
		if err != nil {
			panic(err)
		}
		myViper.Set("Difficulty", selectInt)

		myViper.Set("Running", true)

		// TODO: VRF type,Time Commitment Type,F Function Type

		// write new settings
		if err := myViper.WriteConfig(); err != nil {
			panic(fmt.Errorf("setup conf failed, err:%s", err))
		}

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

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// oldConfig := myViper.AllSettings()
	// fmt.Printf("All settings #1 %+v\n\n", oldConfig)

	myViper.Set("PreviousOutput", output)

	// oldConfig := myViper.AllSettings()
	// fmt.Printf("All settings #1 %+v\n\n", oldConfig)

	// fmt.Println(time.Now())
	// write new settings
	myViper.Debug()
	if err := myViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
	// if err := myViper.ReadInConfig(); err != nil {
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }
	myViper.Debug()

	// if err := myViper.WriteConfigAs("../config.yml"); err != nil {
	// 	panic(fmt.Errorf("setup conf failed, err:%s", err))
	// }

	// fmt.Println(time.Now())
	// // time.Sleep(100 * time.Millisecond)
	// newConfig := myViper.AllSettings()
	// fmt.Printf("All settings #2 %+v\n\n", newConfig)
	ReadConfig()

	fmt.Println("Write output")

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("unlock share lock failed", err)
	}
}

// get difficulty from config
func GetDifficulty() int {
	var configuration = new(Configurations)

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := myViper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	return configuration.Difficulty
}

// get curve from config
func GetCurve() string {
	var configuration = new(Configurations)

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := myViper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	return configuration.EllipticCurve
}

// get previous input from config
func GetPreviousInput() string {
	var configuration = new(Configurations)

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := myViper.Unmarshal(configuration); err != nil {
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

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	switch id {
	case 0:
		myViper.Set("Node0_Ip", ip)
		myViper.Set("Node0_Pk", pk)
	case 1:
		myViper.Set("Node1_Ip", ip)
		myViper.Set("Node1_Pk", pk)
	case 2:
		myViper.Set("Node2_Ip", ip)
		myViper.Set("Node2_Pk", pk)
	case 3:
		myViper.Set("Node3_Ip", ip)
		myViper.Set("Node3_Pk", pk)
	case 4:
		myViper.Set("Node4_Ip", ip)
		myViper.Set("Node4_Pk", pk)
	case 5:
		myViper.Set("Node5_Ip", ip)
		myViper.Set("Node5_Pk", pk)
	case 6:
		myViper.Set("Node6_Ip", ip)
		myViper.Set("Node6_Pk", pk)
	}

	// write new settings
	if err := myViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
	fmt.Println("new consensus node")

	// oldConfig := myViper.AllSettings()
	// fmt.Printf("All settings #1 %+v\n\n", oldConfig)

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

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// handle new id-ip-pk
	switch id {
	case 0:
		if myViper.GetString("Node0_Ip") != "0" {
			myViper.Set("Node0_Ip", "0")
			myViper.Set("Node0_Pk", "0")
		}
	case 1:
		if myViper.GetString("Node1_Ip") != "0" {
			myViper.Set("Node1_Ip", "0")
			myViper.Set("Node1_Pk", "0")
		}
	case 2:
		if myViper.GetString("Node2_Ip.pk") != "0" {
			myViper.Set("Node2_Ip", "0")
			myViper.Set("Node2_Pk", "0")
		}
	case 3:
		if myViper.GetString("Node3_Ip") != "0" {
			myViper.Set("Node3_Ip", "0")
			myViper.Set("Node3_Pk", "0")
		}
	case 4:
		if myViper.GetString("Node4_Ip") != "0" {
			myViper.Set("Node4_Ip", "0")
			myViper.Set("Node4_Pk", "0")
		}
	case 5:
		if myViper.GetString("Node5_Ip") != "0" {
			myViper.Set("Node5_Ip", "0")
			myViper.Set("Node5_Pk", "0")
		}
	case 6:
		if myViper.GetString("Node6_Ip") != "0" {
			myViper.Set("Node6_Ip", "0")
			myViper.Set("Node6_Pk", "0")
		}
	}

	// write new settings
	if err := myViper.WriteConfig(); err != nil {
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

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../config.yml")

	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := myViper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	for i := 0; i < 7; i++ {
		var node NodeConfig
		switch i {
		case 0:
			node.Ip = configuration.Node0_ip
			node.Pk = configuration.Node0_pk
			nodeConfig = append(nodeConfig, node)
		case 1:
			node.Ip = configuration.Node1_ip
			node.Pk = configuration.Node1_pk
			nodeConfig = append(nodeConfig, node)
		case 2:
			node.Ip = configuration.Node2_ip
			node.Pk = configuration.Node2_pk
			nodeConfig = append(nodeConfig, node)
		case 3:
			node.Ip = configuration.Node3_ip
			node.Pk = configuration.Node3_pk
			nodeConfig = append(nodeConfig, node)
		case 4:
			node.Ip = configuration.Node4_ip
			node.Pk = configuration.Node4_pk
			nodeConfig = append(nodeConfig, node)
		case 5:
			node.Ip = configuration.Node5_ip
			node.Pk = configuration.Node5_pk
			nodeConfig = append(nodeConfig, node)
		case 6:
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

	myViper := viper.New()
	// set fonfig file
	myViper.SetConfigFile("../config.yml")

	if err := myViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := myViper.Unmarshal(configuration); err != nil {
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
	fmt.Printf("Difficulty:%d\n", configuration.Difficulty)
	fmt.Printf("VRFType:%d\n", configuration.VRFType)
	fmt.Printf("TCType:%d\n", configuration.TCType)
	fmt.Printf("FType:%d\n", configuration.FType)
}
