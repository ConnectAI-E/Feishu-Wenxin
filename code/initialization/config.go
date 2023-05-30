package initialization

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	FeishuAppId                string
	FeishuAppSecret            string
	FeishuAppEncryptKey        string
	FeishuAppVerificationToken string
	FeishuBotName              string
	OpenaiApiKeys              []string
	HttpPort                   int
	HttpsPort                  int
	UseHttps                   bool
	CertFile                   string
	KeyFile                    string
	OpenaiApiUrl               string
	HttpProxy                  string
	AzureOn                    bool
	AzureApiVersion            string
	AzureDeploymentName        string
	AzureResourceName          string
	AzureOpenaiToken           string
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	//content, err := ioutil.ReadFile("config.yaml")
	//if err != nil {
	//	fmt.Println("Error reading file:", err)
	//}
	//fmt.Println(string(content))

	config := &Config{
		FeishuAppId:                getViperStringValue("APP_ID", ""),
		FeishuAppSecret:            getViperStringValue("APP_SECRET", ""),
		FeishuAppEncryptKey:        getViperStringValue("APP_ENCRYPT_KEY", ""),
		FeishuAppVerificationToken: getViperStringValue("APP_VERIFICATION_TOKEN", ""),
		FeishuBotName:              getViperStringValue("BOT_NAME", ""),
		OpenaiApiKeys:              getViperStringArray("OPENAI_KEY", nil),
		HttpPort:                   getViperIntValue("HTTP_PORT", 9000),
		HttpsPort:                  getViperIntValue("HTTPS_PORT", 9001),
		UseHttps:                   getViperBoolValue("USE_HTTPS", false),
		CertFile:                   getViperStringValue("CERT_FILE", "cert.pem"),
		KeyFile:                    getViperStringValue("KEY_FILE", "key.pem"),
		OpenaiApiUrl:               getViperStringValue("API_URL", "https://api.openai.com"),
		HttpProxy:                  getViperStringValue("HTTP_PROXY", ""),
		AzureOn:                    getViperBoolValue("AZURE_ON", false),
		AzureApiVersion:            getViperStringValue("AZURE_API_VERSION", "2023-03-15-preview"),
		AzureDeploymentName:        getViperStringValue("AZURE_DEPLOYMENT_NAME", ""),
		AzureResourceName:          getViperStringValue("AZURE_RESOURCE_NAME", ""),
		AzureOpenaiToken:           getViperStringValue("AZURE_OPENAI_TOKEN", ""),
	}

	return config
}

func getViperStringValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

//OPENAI_KEY: sk-xxx,sk-xxx,sk-xxx
//result:[sk-xxx sk-xxx sk-xxx]
func getViperStringArray(key string, defaultValue []string) []string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	raw := strings.Split(value, ",")
	return filterFormatKey(raw)
}

func getViperIntValue(key string, defaultValue int) int {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}

func getViperBoolValue(key string, defaultValue bool) bool {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %v\n", key, defaultValue)
		return defaultValue
	}
	return boolValue
}

func (config *Config) GetCertFile() string {
	if config.CertFile == "" {
		return "cert.pem"
	}
	if _, err := os.Stat(config.CertFile); err != nil {
		fmt.Printf("Certificate file %s does not exist, using default file cert.pem\n", config.CertFile)
		return "cert.pem"
	}
	return config.CertFile
}

func (config *Config) GetKeyFile() string {
	if config.KeyFile == "" {
		return "key.pem"
	}
	if _, err := os.Stat(config.KeyFile); err != nil {
		fmt.Printf("Key file %s does not exist, using default file key.pem\n", config.KeyFile)
		return "key.pem"
	}
	return config.KeyFile
}

// 过滤出 "sk-" 开头的 key
func filterFormatKey(keys []string) []string {
	var result []string
	for _, key := range keys {
		if strings.HasPrefix(key, "sk-") || strings.HasPrefix(key,
			"fk") {
			result = append(result, key)
		}
	}
	return result

}
