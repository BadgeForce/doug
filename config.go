package doug

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config struct that holds configs
type Config struct {
	S3Conf    s3                       `toml:"s3"`
	Projects  []map[string]interface{} `toml:"projects"`
	Artifacts map[string][]interface{}
	Github    `toml:"github"`
}

//Github . . .
type Github struct {
	Secret string `toml:"secret"`
}

//AWS . . .
type s3 struct {
	Regions []string `toml:"regions"`
	Bucket  string   `toml:"bucket"`
}

//Configs is the variable that holds config information for microservice
var Configs Config

func postDecode(config *Config) {
	config.Artifacts = make(map[string][]interface{})
	for _, project := range config.Projects {
		name := project["name"].(string)
		artifacts := project["artifacts"].([]interface{})
		config.Artifacts[name] = append(config.Artifacts[name], artifacts...)
	}
}

//Initializes configurations
func InitializeConfig(path string) {
	if _, err := toml.DecodeFile(path, &Configs); err != nil {
		log.Println("Check your configs.")
		log.Fatalln(err.Error())
		return
	}
	postDecode(&Configs)
	if len(Configs.S3Conf.Regions) == 0 {
		log.Fatalln("S3 configs do not have any regions specified")
	}
}
