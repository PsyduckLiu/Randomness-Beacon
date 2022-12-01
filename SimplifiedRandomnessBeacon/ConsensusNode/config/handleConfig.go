package config

import (
	"consensusNode/util"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

type Output struct {
	PreviousOutput string `mapstructure:"previousOutput"`
}

type NodeConfig struct {
	Ip string
	Pk string
}

// write new previousputput
func WriteOutput(output string) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteOutput]Open lock failed:%s", err))
	}
	// share lock, concurrently read
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteOutput]Add share lock failed:%s", err))
	}

	// set config file
	outputViper := viper.New()
	outputViper.SetConfigFile("../Configuration/output.yml")

	// read config and keep origin settings
	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteOutput]Read config file failed:%s", err))
	}

	outputViper.Set("PreviousOutput", output)

	// write new settings
	if err := outputViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteOutput]Write config file failed:%s", err))
	}
	// outputViper.Debug()

	fmt.Println("\n===>[Write Output]Write output success")

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Unlock share lock failed:%s", err))
	}
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

// write new id-ip-pk into config
func NewConsensusNode(id int64, ip string, pk string) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewConsensusNode]Open lock failed:%s", err))
	}
	// share lock, concurrently read
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(fmt.Errorf("===>[ERROR from NewConsensusNode]Add share lock failed:%s", err))
	}

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")

	// read config and keep origin settings
	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from NewConsensusNode]Read config file failed:%s", err))
	}

	var ipString string = "Node" + strconv.FormatInt(id, 10) + "_Ip"
	var pkString string = "Node" + strconv.FormatInt(id, 10) + "_Pk"
	configViper.Set(ipString, ip)
	configViper.Set(pkString, pk)

	// write new settings
	if err := configViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from NewConsensusNode]Write config file failed:%s", err))
	}
	fmt.Printf("\n===>[NewConsensusNode]Add new consensus node[%d] success\n", id)

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		panic(fmt.Errorf("===>[ERROR from NewConsensusNode]Unlock share lock failed:%s", err))
	}
}

// remove consensus node from config file
func RemoveConsensusNode(id int64) {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from RemoveConsensusNode]Open lock failed:%s", err))
	}
	// share lock, concurrently read
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(fmt.Errorf("===>[ERROR from RemoveConsensusNode]Add share lock failed:%s", err))
	}

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")

	// read config and keep origin settings
	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from RemoveConsensusNode]Read config file failed:%s", err))
	}

	// handle new id-ip-pk
	var ipString string = "Node" + strconv.FormatInt(id, 10) + "_Ip"
	var pkString string = "Node" + strconv.FormatInt(id, 10) + "_Pk"
	if configViper.GetString(ipString) != "0" {
		configViper.Set(ipString, "0")
		configViper.Set(pkString, "0")
	}

	// write new settings
	if err := configViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from RemoveConsensusNode]Write config file failed:%s", err))
	}
	fmt.Println("===>[RemoveConsensusNode]remove consensus node success")

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		panic(fmt.Errorf("===>[ERROR from RemoveConsensusNode]Unlock share lock failed:%s", err))
	}
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
