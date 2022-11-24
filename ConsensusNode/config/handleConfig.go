package config

import (
	"consensusNode/util"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

// type Configurations struct {
// 	Running       bool   `mapstructure:"running"`
// 	Version       string `mapstructure:"version"`
// 	EllipticCurve string `mapstructure:"ellipticCurve"`
// 	VRFType       int    `mapstructure:"vrfType"`
// 	TCType        int    `mapstructure:"tcType"`
// 	FType         int    `mapstructure:"fType "`
// 	Difficulty    int    `mapstructure:"difficulty"`
// Node0_ip      string `mapstructure:"node0_ip"`
// Node0_pk      string `mapstructure:"node0_pk"`
// Node1_ip      string `mapstructure:"node1_ip"`
// Node1_pk      string `mapstructure:"node1_pk"`
// Node2_ip      string `mapstructure:"node2_ip"`
// Node2_pk      string `mapstructure:"node2_pk"`
// Node3_ip      string `mapstructure:"node3_ip"`
// Node3_pk      string `mapstructure:"node3_pk"`
// Node4_ip      string `mapstructure:"node4_ip"`
// Node4_pk      string `mapstructure:"node4_pk"`
// Node5_ip      string `mapstructure:"node5_ip"`
// Node5_pk      string `mapstructure:"node5_pk"`
// Node6_ip      string `mapstructure:"node6_ip"`
// Node6_pk      string `mapstructure:"node6_pk"`
// Node7_ip      string `mapstructure:"node7_ip"`
// Node7_pk      string `mapstructure:"node7_pk"`
// Node8_ip      string `mapstructure:"node8_ip"`
// Node8_pk      string `mapstructure:"node8_pk"`
// Node9_ip      string `mapstructure:"node9_ip"`
// Node9_pk      string `mapstructure:"node9_pk"`
// }

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

	// switch id {
	// case 0:
	// 	configViper.Set("Node0_Ip", ip)
	// 	configViper.Set("Node0_Pk", pk)
	// case 1:
	// 	configViper.Set("Node1_Ip", ip)
	// 	configViper.Set("Node1_Pk", pk)
	// case 2:
	// 	configViper.Set("Node2_Ip", ip)
	// 	configViper.Set("Node2_Pk", pk)
	// case 3:
	// 	configViper.Set("Node3_Ip", ip)
	// 	configViper.Set("Node3_Pk", pk)
	// case 4:
	// 	configViper.Set("Node4_Ip", ip)
	// 	configViper.Set("Node4_Pk", pk)
	// case 5:
	// 	configViper.Set("Node5_Ip", ip)
	// 	configViper.Set("Node5_Pk", pk)
	// case 6:
	// 	configViper.Set("Node6_Ip", ip)
	// 	configViper.Set("Node6_Pk", pk)
	// case 7:
	// 	configViper.Set("Node7_Ip", ip)
	// 	configViper.Set("Node7_Pk", pk)
	// case 8:
	// 	configViper.Set("Node8_Ip", ip)
	// 	configViper.Set("Node8_Pk", pk)
	// case 9:
	// 	configViper.Set("Node9_Ip", ip)
	// 	configViper.Set("Node9_Pk", pk)
	// }

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

	// switch id {
	// case 0:
	// 	if configViper.GetString("Node0_Ip") != "0" {
	// 		configViper.Set("Node0_Ip", "0")
	// 		configViper.Set("Node0_Pk", "0")
	// 	}
	// case 1:
	// 	if configViper.GetString("Node1_Ip") != "0" {
	// 		configViper.Set("Node1_Ip", "0")
	// 		configViper.Set("Node1_Pk", "0")
	// 	}
	// case 2:
	// 	if configViper.GetString("Node2_Ip") != "0" {
	// 		configViper.Set("Node2_Ip", "0")
	// 		configViper.Set("Node2_Pk", "0")
	// 	}
	// case 3:
	// 	if configViper.GetString("Node3_Ip") != "0" {
	// 		configViper.Set("Node3_Ip", "0")
	// 		configViper.Set("Node3_Pk", "0")
	// 	}
	// case 4:
	// 	if configViper.GetString("Node4_Ip") != "0" {
	// 		configViper.Set("Node4_Ip", "0")
	// 		configViper.Set("Node4_Pk", "0")
	// 	}
	// case 5:
	// 	if configViper.GetString("Node5_Ip") != "0" {
	// 		configViper.Set("Node5_Ip", "0")
	// 		configViper.Set("Node5_Pk", "0")
	// 	}
	// case 6:
	// 	if configViper.GetString("Node6_Ip") != "0" {
	// 		configViper.Set("Node6_Ip", "0")
	// 		configViper.Set("Node6_Pk", "0")
	// 	}
	// case 7:
	// 	if configViper.GetString("Node7_Ip") != "0" {
	// 		configViper.Set("Node7_Ip", "0")
	// 		configViper.Set("Node7_Pk", "0")
	// 	}
	// case 8:
	// 	if configViper.GetString("Node8_Ip") != "0" {
	// 		configViper.Set("Node8_Ip", "0")
	// 		configViper.Set("Node8_Pk", "0")
	// 	}
	// case 9:
	// 	if configViper.GetString("Node9_Ip") != "0" {
	// 		configViper.Set("Node9_Ip", "0")
	// 		configViper.Set("Node9_Pk", "0")
	// 	}
	// }

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
	// for i := 0; i < 7; i++ {
	// 	var node NodeConfig
	// 	switch i {
	// 	case 0:
	// 		node.Ip = configuration.Node0_ip
	// 		node.Pk = configuration.Node0_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 1:
	// 		node.Ip = configuration.Node1_ip
	// 		node.Pk = configuration.Node1_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 2:
	// 		node.Ip = configuration.Node2_ip
	// 		node.Pk = configuration.Node2_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 3:
	// 		node.Ip = configuration.Node3_ip
	// 		node.Pk = configuration.Node3_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 4:
	// 		node.Ip = configuration.Node4_ip
	// 		node.Pk = configuration.Node4_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 5:
	// 		node.Ip = configuration.Node5_ip
	// 		node.Pk = configuration.Node5_pk
	// 		nodeConfig = append(nodeConfig, node)
	// 	case 6:
	// 		node.Ip = configuration.Node6_ip
	// 		node.Pk = configuration.Node6_pk
	// 		nodeConfig = append(nodeConfig, node)
	// case 7:
	// 	node.Ip = configuration.Node7_ip
	// 	node.Pk = configuration.Node7_pk
	// 	nodeConfig = append(nodeConfig, node)
	// case 8:
	// 	node.Ip = configuration.Node8_ip
	// 	node.Pk = configuration.Node8_pk
	// 	nodeConfig = append(nodeConfig, node)
	// case 9:
	// 	node.Ip = configuration.Node9_ip
	// 	node.Pk = configuration.Node9_pk
	// 	nodeConfig = append(nodeConfig, node)
	// 	}
	// }

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
	// fmt.Printf("Node[0]'s ip is %s\n", configuration.Node0_ip)
	// fmt.Printf("Node[0]'s pk is %s\n", configuration.Node0_pk)
	// fmt.Printf("Node[1]'s ip is %s\n", configuration.Node1_ip)
	// fmt.Printf("Node[1]'s pk is %s\n", configuration.Node1_pk)
	// fmt.Printf("Node[2]'s ip is %s\n", configuration.Node2_ip)
	// fmt.Printf("Node[2]'s pk is %s\n", configuration.Node2_pk)
	// fmt.Printf("Node[3]'s ip is %s\n", configuration.Node3_ip)
	// fmt.Printf("Node[3]'s pk is %s\n", configuration.Node3_pk)
	// fmt.Printf("Node[4]'s ip is %s\n", configuration.Node4_ip)
	// fmt.Printf("Node[4]'s pk is %s\n", configuration.Node4_pk)
	// fmt.Printf("Node[5]'s ip is %s\n", configuration.Node5_ip)
	// fmt.Printf("Node[5]'s pk is %s\n", configuration.Node5_pk)
	// fmt.Printf("Node[6]'s ip is %s\n", configuration.Node6_ip)
	// fmt.Printf("Node[6]'s pk is %s\n", configuration.Node6_pk)
	// fmt.Printf("Node[7]'s ip is %s\n", configuration.Node7_ip)
	// fmt.Printf("Node[7]'s pk is %s\n", configuration.Node7_pk)
	// fmt.Printf("Node[8]'s ip is %s\n", configuration.Node8_ip)
	// fmt.Printf("Node[8]'s pk is %s\n", configuration.Node8_pk)
	// fmt.Printf("Node[9]'s ip is %s\n", configuration.Node9_ip)
	// fmt.Printf("Node[9]'s pk is %s\n", configuration.Node9_pk)
	// fmt.Printf("Difficulty:%d\n", configuration.Difficulty)
	// fmt.Printf("VRFType:%d\n", configuration.VRFType)
	// fmt.Printf("TCType:%d\n", configuration.TCType)
	// fmt.Printf("FType:%d\n", configuration.FType)
}
