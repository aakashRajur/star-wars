package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func PrintEnv() {
	for _, each := range os.Environ() {
		pair := strings.Split(each, "=")
		fmt.Println(pair[0], "=", pair[1])
	}
}

func LoadEnv() {
	env := os.Getenv("env")
	_ = godotenv.Load("wiki.env." + env + ".local")
	if env != "test" {
		_ = godotenv.Load("wiki.env.local")
	}
	_ = godotenv.Load("wiki.env." + env)
	_ = godotenv.Load()
}

func GetInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return -1
	}
	return value
}

func GetBool(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return value
}

func GetString(key string) string {
	return os.Getenv(key)
}
