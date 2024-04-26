#!/bin/bash
go mod init url-shortener
go mod tidy
go get github.com/ilyakaznacheev/cleanenv
go get github.com/joho/godotenv