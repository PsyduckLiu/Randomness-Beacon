package config

import (
	"crypto/elliptic"
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Version        string
	PreviousOutput string
	EllipticCurve  elliptic.Curve
}

var configuration = new(Configurations)

func ReadConfig() {
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
	fmt.Println(configuration)

	fmt.Printf("reading using model:\nversion=%s,previousOutput=%s,ellipticCurve=%v\n",
		configuration.Version,
		configuration.PreviousOutput,
		configuration.EllipticCurve,
	)

	fmt.Printf("reading without model:\nversion=%s,previousOutput=%s\n",
		viper.GetString("version"),
		viper.GetString("previousOutput"),
	)
}
