package config

import (
    "github.com/cihub/seelog"
    "gopkg.in/yaml.v2"
    "os"
    "sync"
    "time"
)

// Configuration 项目配置
type Configuration struct {
    // ChatGPT请求的最大并发数
    Listen              string              `json:"listen" yaml:"listen"`
    Concurrency         int                 `json:"concurrency" yaml:"concurrency"`
    Guests              GuestConfigs        `yaml:"guests"`
    Gpt                 *GptConfig          `json:"gpt" yaml:"gpt"`
    TencentCos          TencentCos          `json:"tencent_cos" yaml:"tencent_cos"`
    BaiduAi             BaiduAi             `json:"baidu_ai" yaml:"baidu_ai"`
    Db                  DbConfig            `json:"db" yaml:"db"`
    Mock                MockConfig          `yaml:"mock"`
    OpenAiAccountConfig OpenAiAccountConfig `yaml:"open_ai_account"`
    CensorEnabled       bool                `yaml:"censor_enabled"`
    ApiPlatform         ApiPlatform         `yaml:"api_platform"`
}

type ApiPlatform struct {
    // 请求地址，默认为 https://api.openai.com/v1/
    BaseUrl string `json:"base_url" yaml:"base_url"`
}

type OpenAiAccountConfig struct {
    QueryInterval time.Duration `yaml:"query_interval"`
    Concurrency   int           `yaml:"concurrency"`
}

type MockConfig struct {
    Enabled      bool   `yaml:"enabled"`
    Response     string `yaml:"response"`
    FreeOfCharge bool   `yaml:"free_of_charge"` // 是否免费
}

func (r *MockConfig) IsFree() bool {
    if !r.Enabled {
        return false
    }
    if !r.FreeOfCharge {
        return false
    }
    return true
}

type DbConfig struct {
    Host string `json:"host" yaml:"host"`
    Port string `json:"port" yaml:"port"`
    Name string `json:"name" yaml:"name"`
    User string `json:"user" yaml:"user"`
    Pass string `json:"pass" yaml:"pass"`
}

type Account struct {
    Email    string `json:"email" yaml:"email"`
    Password string `json:"password" yaml:"password"`
    ApiKey   string `json:"api_key" yaml:"api_key"`
    Token    string `json:"token" yaml:"token"`
}

// GptConfig 项目配置
type GptConfig struct {
    // ChatGPT请求的最大并发数
    Concurrency int `json:"concurrency" yaml:"concurrency"`
    // 自动通过好友
    AutoPass bool `json:"auto_pass" yaml:"auto_pass"`
    // 会话超时时间
    SessionTimeout time.Duration `json:"session_timeout" yaml:"session_timeout"`
    // GPT请求最大字符数
    MaxTokens uint `json:"max_tokens" yaml:"max_tokens"`
    // GPT模型
    Model string `json:"model" yaml:"model"`
    // 热度
    Temperature float64 `json:"temperature" yaml:"temperature"`
    // 回复前缀
    ReplyPrefix string `json:"reply_prefix" yaml:"reply_prefix"`
}

type TencentCos struct {
    SecretId   string `json:"secret_id" yaml:"secret_id"`
    SecretKey  string `json:"secret_key" yaml:"secret_key"`
    BucketUrl  string `json:"bucket_url" yaml:"bucket_url"`
    ServiceUrl string `json:"service_url" yaml:"service_url"`
    CiUrl      string `json:"ci_url" yaml:"ci_url"`
}
type BaiduAi struct {
    AppKey    string `json:"app_key" yaml:"app_key"`
    SecretKey string `json:"secret_key" yaml:"secret_key"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
    defer func() {
        seelog.Infof("当前配置：%+v", config)
    }()
    once.Do(func() {
        configFile := "配置/config.yml"
        // 给配置赋默认值
        config = &Configuration{
            Concurrency: 1,
            Gpt: &GptConfig{
                AutoPass:       false,
                SessionTimeout: 60,
                MaxTokens:      512,
                Model:          "text-davinci-003",
                Temperature:    0.9,
            },
        }

        // 判断配置文件是否存在，存在直接JSON读取
        _, err := os.Stat(configFile)
        if err == nil {
            f, err := os.Open(configFile)
            if err != nil {
                seelog.Criticalf("open config err: %v", err)
                return
            }
            defer f.Close()
            //encoder := json.NewDecoder(f)
            //err = encoder.Decode(config)
            //if err != nil {
            //    seelog.Criticalf("decode config err: %v", err)
            //    return
            //}
            encoder := yaml.NewDecoder(f)
            err = encoder.Decode(config)
            if err != nil {
                seelog.Criticalf("decode config err: %v", err)
                return
            }
        }
        // 有环境变量使用环境变量
        dbHost := os.Getenv("DB_HOST")
        if dbHost != "" {
            config.Db.Host = dbHost
        }
        dbUser := os.Getenv("DB_USER")
        if dbUser != "" {
            config.Db.User = dbUser
        }
        dbPass := os.Getenv("DB_PASS")
        if dbPass != "" {
            config.Db.Pass = dbPass
        }
        dbName := os.Getenv("DB_NAME")
        if dbName != "" {
            config.Db.Name = dbName
        }
        dbPort := os.Getenv("DB_PORT")
        if dbPort != "" {
            config.Db.Port = dbPort
        }
        baseUrl := os.Getenv("AP_BASE_URL")
        if baseUrl != "" {
            config.ApiPlatform.BaseUrl = baseUrl
        }
    })

    return config
}
