package config

import "github.com/fortnoxab/fnxlogrus"

type Config struct {
	Log               fnxlogrus.Config
	Token             string
	BitbucketURL      string
	BitbucketUser     string
	BitbucketPassword string
	Port              string `default:"8080"`
}
