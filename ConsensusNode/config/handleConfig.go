package config

import (
	"consensusNode/signature"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Running        bool
	Version        string
	PreviousOutput []byte
	EllipticCurve  []byte
	VRFType        int
	TCType         int
	FType          int
	ConsensusNodes NodesConfig
}

type NodesConfig struct {
	Node0 NodeConfig
	Node1 NodeConfig
	Node2 NodeConfig
	Node3 NodeConfig
	Node4 NodeConfig
	Node5 NodeConfig
	Node6 NodeConfig
}

type NodeConfig struct {
	Pk []byte
	Ip string
}

// config file Setup
// assign curve,previous output,VRF type,Time Commitment Type,F Function Type
func SetupConfig() {
	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

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
		viper.Set("EllipticCurve", marshalledKey)

		// generate random init input
		message := []byte("asdkjhdk")
		randomNum := signature.Digest(message)
		viper.Set("PreviousOutput", randomNum)
		// TODO: VRF type,Time Commitment Type,F Function Type

		viper.Set("Running", true)

		// write new settings
		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("setup conf failed, err:%s", err))
		}
	}
}

// get curve from config
func GetCurve() []byte {
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
func GetPreviousInput() []byte {
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
func NewConsensusNode(id int64, ip string, pk []byte) {
	// set config file
	viper.SetConfigFile("../config.yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// handle new id-ip-pk
	switch id {
	case 0:
		viper.Set("consensusnodes.node0.pk", pk)
		viper.Set("consensusnodes.node0.ip", ip)
	case 1:
		viper.Set("consensusnodes.node1.pk", pk)
		viper.Set("consensusnodes.node1.ip", ip)
	case 2:
		viper.Set("consensusnodes.node2.pk", pk)
		viper.Set("consensusnodes.node2.ip", ip)
	case 3:
		viper.Set("consensusnodes.node3.pk", pk)
		viper.Set("consensusnodes.node3.ip", ip)
	case 4:
		viper.Set("consensusnodes.node4.pk", pk)
		viper.Set("consensusnodes.node4.ip", ip)
	case 5:
		viper.Set("consensusnodes.node5.pk", pk)
		viper.Set("consensusnodes.node5.ip", ip)
	case 6:
		viper.Set("consensusnodes.node6.pk", pk)
		viper.Set("consensusnodes.node6.ip", ip)
	}

	// write new settings
	if err := viper.WriteConfig(); err != nil {
		panic(fmt.Errorf("setup conf failed, err:%s", err))
	}
}

func RemoveConsensusNode(id int64) {
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
}

// get consensus nodes from config
func GetConsensusNode() [7]NodeConfig {
	var nodeConfig [7]NodeConfig
	var configuration = new(Configurations)

	// set config file
	viper.SetConfigFile("../config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	nodeConfig[0].Ip = configuration.ConsensusNodes.Node0.Ip
	nodeConfig[0].Pk = configuration.ConsensusNodes.Node0.Pk
	nodeConfig[1].Ip = configuration.ConsensusNodes.Node1.Ip
	nodeConfig[1].Pk = configuration.ConsensusNodes.Node1.Pk
	nodeConfig[2].Ip = configuration.ConsensusNodes.Node2.Ip
	nodeConfig[2].Pk = configuration.ConsensusNodes.Node2.Pk
	nodeConfig[3].Ip = configuration.ConsensusNodes.Node3.Ip
	nodeConfig[3].Pk = configuration.ConsensusNodes.Node3.Pk
	nodeConfig[4].Ip = configuration.ConsensusNodes.Node4.Ip
	nodeConfig[4].Pk = configuration.ConsensusNodes.Node4.Pk
	nodeConfig[5].Ip = configuration.ConsensusNodes.Node5.Ip
	nodeConfig[5].Pk = configuration.ConsensusNodes.Node5.Pk
	nodeConfig[6].Ip = configuration.ConsensusNodes.Node6.Ip
	nodeConfig[6].Pk = configuration.ConsensusNodes.Node6.Pk

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
	fmt.Printf("Node[0]'s ip is %s\n", configuration.ConsensusNodes.Node0.Ip)
	fmt.Printf("Node[0]'s pk is %s\n", configuration.ConsensusNodes.Node0.Pk)
	fmt.Printf("Node[1]'s ip is %s\n", configuration.ConsensusNodes.Node1.Ip)
	fmt.Printf("Node[1]'s pk is %s\n", configuration.ConsensusNodes.Node1.Pk)
	fmt.Printf("Node[2]'s ip is %s\n", configuration.ConsensusNodes.Node2.Ip)
	fmt.Printf("Node[2]'s pk is %s\n", configuration.ConsensusNodes.Node2.Pk)
	fmt.Printf("Node[3]'s ip is %s\n", configuration.ConsensusNodes.Node3.Ip)
	fmt.Printf("Node[3]'s pk is %s\n", configuration.ConsensusNodes.Node3.Pk)
	fmt.Printf("Node[4]'s ip is %s\n", configuration.ConsensusNodes.Node4.Ip)
	fmt.Printf("Node[4]'s pk is %s\n", configuration.ConsensusNodes.Node4.Pk)
	fmt.Printf("Node[5]'s ip is %s\n", configuration.ConsensusNodes.Node5.Ip)
	fmt.Printf("Node[5]'s pk is %s\n", configuration.ConsensusNodes.Node5.Pk)
	fmt.Printf("Node[6]'s ip is %s\n", configuration.ConsensusNodes.Node6.Ip)
	fmt.Printf("Node[6]'s pk is %s\n", configuration.ConsensusNodes.Node6.Pk)
	fmt.Printf("VRFType:%d\n", configuration.VRFType)
	fmt.Printf("TCType:%d\n", configuration.TCType)
	fmt.Printf("FType:%d\n", configuration.FType)
}