package config

import (
	"fmt"

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
	Node0_ip       string
	Node0_pk       string
	Node1_ip       string
	Node1_pk       string
	Node2_ip       string
	Node2_pk       string
	Node3_ip       string
	Node3_pk       string
	Node4_ip       string
	Node4_pk       string
	Node5_ip       string
	Node5_pk       string
	Node6_ip       string
	Node6_pk       string
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
	Ip string
	Pk string
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
