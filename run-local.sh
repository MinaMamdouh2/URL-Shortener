#!/bin/bash

# Load env vars
export $(grep -v '^#' .env | xargs)

go run app/services/url-shortener-api/main.go