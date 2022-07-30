package main

import (
	"./network"
	"./system"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

func init() {
	fmt.Println("Initializing Home Assistant PC control server...")
}

func main() {
	// Generate our config based on the config supplied by the user in the flags
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	// Run the server
	cfg.Run()
}

func getIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "HA PC control server",
		"ip":       network.GetLocalIP(),
		"hostname": network.GetHostname(),
	})
}

// Config struct for config
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Auth struct {
		User  string `yaml:"user"`
		Token string `yaml:"token"`
	} `yaml:"auth"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file, that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

// Run will run the server
func (config Config) Run() {
	router := gin.Default()
	router.GET("/", getIndex)

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		config.Auth.User: config.Auth.Token,
	}))

	authorized.POST("/system/:action", system.PostSystem)

	err := router.Run(config.Server.Host + ":" + config.Server.Port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
