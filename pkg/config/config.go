package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConf(path string, genre string, conf interface{}) error {
	viper.SetConfigType(genre)
	//viper.AutomaticEnv()         // read in environment variables that match
	//viper.SetEnvPrefix("gorush") // will be uppercased automatically
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//viper.AddConfigPath(confPath)
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("read config file error", err)
		return err
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Println("unmarshal config file error", err)
		return err
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
	return err
}
