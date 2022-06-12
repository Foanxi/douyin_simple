package initalize

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error resources file: %w \n", err))
	}

	if err := viper.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct %w \n", err))
	}
}
