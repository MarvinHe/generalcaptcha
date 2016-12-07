package generalcaptcha

import (
	"math/rand"
	"time"
	"strconv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"github.com/pkg/errors"
)

// every mobile should get captcha less than this number
var maxDailyCount int = 10
var timeout time.Duration = 5 * 60 * time.Second
var redisStore RedisStore

type CaptchaService struct {
	maxDailyCount int
	timeout time.Duration
	redisStore RedisStore
}

func init() {
	config := loanConfigFromYaml("config.yaml")
	timeout = time.Duration(config["Timeout"].(int)) * 60 * time.Second
	maxDailyCount = config["MaxDailyCount"].(int)
	redisStore = RedisStore{
		RedisKeyPrefix: config["Prefix"].(string),
		Options: EnsureMapString(config["Options"].(map[interface{}]interface{}))}
	redisStore.InitialClient()
	rand.Seed(time.Now().UnixNano())
}

func loanConfigFromYaml(filename string) (config map[string]interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(data, &config)
	return config
}

func GenerateCaptcha(mobile string) (captcha string, err error) {
	if redisStore.GetCountInDay(mobile) >= maxDailyCount {
		captcha, err = "", errors.New("exceed the maxDailyCount")
	} else {
		captcha = calculateCaptcha(mobile)
		redisStore.Store(mobile, captcha, timeout)
		redisStore.IncreCountInDay(mobile)
	}
	return
}

func CheckCaptcha(mobile, captcha string) bool {
	return captcha != "" && captcha == getCaptcha(mobile)
}

func getCaptcha(mobile string) string {
	return redisStore.Get(mobile)
}

func calculateCaptcha(mobile string) string {
	return strconv.Itoa(rand.Intn(1000000))
}
