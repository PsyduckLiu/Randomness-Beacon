package config

import (
	"entropyNode/util"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Output struct {
	PreviousOutput string `mapstructure:"previousOutput"`
}

type NodeConfig struct {
	Ip string
	Pk string
}

// get difficulty from config file
func GetDifficulty() int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetDifficulty]Read config file failed:%s", err))
	}

	return configViper.GetInt("Difficulty")
}

// get curve from config file
func GetCurve() string {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetCurve]Read config file failed:%s", err))
	}

	return configViper.GetString("EllipticCurve")
}

// get previous output from config file
func GetPreviousOutput() string {
	// set config file
	outputViper := viper.New()
	outputViper.SetConfigFile("../Configuration/output.yml")

	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetPreviousOutput]Read config file failed:%s", err))
	}

	return outputViper.GetString("PreviousOutput")
}

// get consensus nodes from config file
func GetConsensusNode() []NodeConfig {
	var nodeConfig []NodeConfig
	// var configuration = new(Configurations)

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetConsensusNode]Read config file failed:%s", err))
	}
	// if err := configViper.Unmarshal(configuration); err != nil {
	// 	panic(fmt.Errorf("===>[ERROR from GetConsensusNode]Unmarshal conf failed:%s", err))
	// }

	for i := 0; i < util.TotalNodeNum; i++ {
		var node NodeConfig
		var ipString string = "Node" + strconv.FormatInt(int64(i), 10) + "_Ip"
		var pkString string = "Node" + strconv.FormatInt(int64(i), 10) + "_Pk"
		node.Ip = configViper.GetString(ipString)
		node.Pk = configViper.GetString(pkString)

		nodeConfig = append(nodeConfig, node)
	}

	return nodeConfig
}

// read configurations from config file
func ReadConfig() {
	// var configuration = new(Configurations)

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")
	outputViper := viper.New()
	outputViper.SetConfigFile("../Configuration/output.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from ReadConfig]Read config file failed:%s", err))
	}
	// if err := configViper.Unmarshal(configuration); err != nil {
	// 	panic(fmt.Errorf("===>[ERROR from ReadConfig]Unmarshal conf failed:%s", err))
	// }

	fmt.Printf("\nReading Configuration:\n")
	fmt.Printf("Running:%v\n", configViper.GetString("Running"))
	fmt.Printf("Version:%s\n", configViper.GetString("Version"))
	fmt.Printf("PreviousOutput:%s\n", outputViper.GetString("PreviousOutput"))
	fmt.Printf("EllipticCurve:%v\n", configViper.GetString("EllipticCurve"))
	fmt.Printf("Consensusnodes:\n")
}
