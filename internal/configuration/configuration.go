package configuration

import (
	"fmt"
	"strings"

	"github.com/FloRichardPro/minimal-api/internal/controllers"
	"github.com/FloRichardPro/minimal-api/internal/routers"
	"github.com/FloRichardPro/minimal-api/internal/services"
	"github.com/spf13/viper"
)

var (
	// Config holds all config data for the all API.
	Config *Conf
)

type Conf struct {
	Routers     *routers.Conf     `mapstructure:"routers"`
	Controllers *controllers.Conf `mapstructure:"controllers"`
	Services    *services.Conf    `mapstructure:"services"`
}

func (c *Conf) LoadConf(path string) error {
	Config = new(Conf)
	Config.Routers = &routers.Config
	Config.Controllers = &controllers.Config
	Config.Services = &services.Config

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("can't load API configuration : %w", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("can't unmarshall API configuration : %w", err)
	}

	// routers.Config = *c.Routers

	return nil
}
