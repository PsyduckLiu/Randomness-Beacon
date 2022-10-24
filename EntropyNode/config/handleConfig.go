package config

import (
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
