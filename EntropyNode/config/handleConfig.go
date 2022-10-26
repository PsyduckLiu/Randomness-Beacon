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
	Difficulty     int
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
			node.Ip = configuration.ConsensusNodes.Node0.Ip
			node.Pk = configuration.ConsensusNodes.Node0.Pk
			nodeConfig = append(nodeConfig, node)
		case 1:
			node.Ip = configuration.ConsensusNodes.Node1.Ip
			node.Pk = configuration.ConsensusNodes.Node1.Pk
			nodeConfig = append(nodeConfig, node)
		case 2:
			node.Ip = configuration.ConsensusNodes.Node2.Ip
			node.Pk = configuration.ConsensusNodes.Node2.Pk
			nodeConfig = append(nodeConfig, node)
		case 3:
			node.Ip = configuration.ConsensusNodes.Node3.Ip
			node.Pk = configuration.ConsensusNodes.Node3.Pk
			nodeConfig = append(nodeConfig, node)
		case 4:
			node.Ip = configuration.ConsensusNodes.Node4.Ip
			node.Pk = configuration.ConsensusNodes.Node4.Pk
			nodeConfig = append(nodeConfig, node)
		case 5:
			node.Ip = configuration.ConsensusNodes.Node5.Ip
			node.Pk = configuration.ConsensusNodes.Node5.Pk
			nodeConfig = append(nodeConfig, node)
		case 6:
			node.Ip = configuration.ConsensusNodes.Node6.Ip
			node.Pk = configuration.ConsensusNodes.Node6.Pk
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
	fmt.Printf("Difficulty:%d\n", configuration.Difficulty)
	fmt.Printf("VRFType:%d\n", configuration.VRFType)
	fmt.Printf("TCType:%d\n", configuration.TCType)
	fmt.Printf("FType:%d\n", configuration.FType)
}
