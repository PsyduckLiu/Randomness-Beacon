package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Running bool   `mapstructure:"running"`
	Version string `mapstructure:"version"`
	// PreviousOutput string `mapstructure:"previousOutput"`
	EllipticCurve string `mapstructure:"ellipticCurve"`
	VRFType       int    `mapstructure:"vrfType"`
	TCType        int    `mapstructure:"tcType"`
	FType         int    `mapstructure:"fType "`
	Difficulty    int    `mapstructure:"difficulty"`
	Node0_ip      string `mapstructure:"node0_ip"`
	Node0_pk      string `mapstructure:"node0_pk"`
	Node1_ip      string `mapstructure:"node1_ip"`
	Node1_pk      string `mapstructure:"node1_pk"`
	Node2_ip      string `mapstructure:"node2_ip"`
	Node2_pk      string `mapstructure:"node2_pk"`
	Node3_ip      string `mapstructure:"node3_ip"`
	Node3_pk      string `mapstructure:"node3_pk"`
	Node4_ip      string `mapstructure:"node4_ip"`
	Node4_pk      string `mapstructure:"node4_pk"`
	Node5_ip      string `mapstructure:"node5_ip"`
	Node5_pk      string `mapstructure:"node5_pk"`
	Node6_ip      string `mapstructure:"node6_ip"`
	Node6_pk      string `mapstructure:"node6_pk"`
}

type Output struct {
	PreviousOutput string `mapstructure:"previousOutput"`
}

type NodeConfig struct {
	Ip string
	Pk string
}

// get difficulty from config
func GetDifficulty() int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return configViper.GetInt("Difficulty")
}

// get curve from config
func GetCurve() string {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return configViper.GetString("EllipticCurve")
}

// get previous output from config
func GetPreviousOutput() string {
	// set config file
	outputViper := viper.New()
	outputViper.SetConfigFile("../output.yml")

	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return outputViper.GetString("PreviousOutput")
}

// get consensus nodes from config
func GetConsensusNode() []NodeConfig {
	var nodeConfig []NodeConfig
	var configuration = new(Configurations)

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := configViper.Unmarshal(configuration); err != nil {
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

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../config.yml")
	outputViper := viper.New()
	outputViper.SetConfigFile("../output.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := configViper.Unmarshal(configuration); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s", err))
	}

	fmt.Printf("Reading using model:\n")
	fmt.Printf("Running:%v\n", configuration.Running)
	fmt.Printf("Version:%s\n", configuration.Version)
	fmt.Printf("PreviousOutput:%s\n", outputViper.GetString("PreviousOutput"))
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
