# Central oauth2 sso for services

[![Build Status](https://travis-ci.org/vedhavyas/oauth2_central.svg?branch=master)](https://travis-ci.org/vedhavyas/oauth2_central)

[![GoDoc](https://godoc.org/github.com/vedhavyas/oauth2_central?status.png)](https://godoc.org/github.com/vedhavyas/oauth2_central)

[![Go Report Card](https://goreportcard.com/badge/github.com/vedhavyas/oauth2_central)](https://goreportcard.com/report/github.com/vedhavyas/oauth2_central)

[![Code Climate](https://codeclimate.com/github/vedhavyas/oauth2_central/badges/gpa.svg)](https://codeclimate.com/github/vedhavyas/oauth2_central)

[![Test Coverage](https://codeclimate.com/github/vedhavyas/oauth2_central/badges/coverage.svg)](https://codeclimate.com/github/vedhavyas/oauth2_central/coverage)

## Config File
Rename config.json.example to config.json and fill in all details

Pass the path of the config file as command line argument like this - ./oauth2_central -config-file=path/to/file
if none is passed, program will look for config.json in the project root.

## Test, Install, and Run
`make all` to test and build the project
`./oauth2_central` to run the project

