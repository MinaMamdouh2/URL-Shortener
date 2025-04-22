#!/bin/bash

# Load env vars
export $(grep -v '^#' .env | xargs)

cd app/services/url-shortener-api && air
