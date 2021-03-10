#/bin/bash

export GO_ENV=production
LOG_FILE="log-"$GO_ENV"-"`date '+%Y%m%d'`".txt"

go run app.go  #2>&1 | tee $LOG_FILE