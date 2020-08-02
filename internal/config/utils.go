package config

import (
	"fmt"

	"github.com/kathleenfrench/pls/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
)

const configPathRelativeToHome = ".config/pls"
const configFileName = "config"
const configFileType = "yaml"

func constructConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		utils.ExitWithError(err)
	}

	return fmt.Sprintf("%s/%s", home, configPathRelativeToHome)
}
