package config

import (
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
	PreviousOutput string
	EllipticCurve  []byte
	VRFType        int
	TCType         int
	FType          int
	ConsensusNodes NodesConfig
}

type NodesConfig struct {
	node0 NodeConfig
	node1 NodeConfig
	node2 NodeConfig
	node3 NodeConfig
	node4 NodeConfig
	node5 NodeConfig
	node6 NodeConfig
}

type NodeConfig struct {
	pk   string
	port string
}

// config file Setup
// assign curve,previous output,VRF type,Time Commitment Type,F Function Type
func SetupConfig() {
	// set fonfig file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	// read config and keep origin settings
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// first time setup
	if viper.GetBool("Running") == false {
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
		// TODO: VRF type,Time Commitment Type,F Function Type

		viper.Set("Running", true)

		// write new settings
		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("setup conf failed, err:%s", err))
		}
	}

}

// get messages from config file
func ReadConfig() {
	var configuration = new(Configurations)

	// set fonfig file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

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
	fmt.Printf("VRFType:%d\n", configuration.VRFType)
	fmt.Printf("TCType:%d\n", configuration.TCType)
	fmt.Printf("FType:%d\n", configuration.FType)
}
