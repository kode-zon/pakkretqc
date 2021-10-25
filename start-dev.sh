processFn() {
    if [ -z "${PAKKRETQC_ALM_ENDPOINT}"]; then
        echo "Error!! missing environment variable "
        echo "Example: "
        echo "  export PAKKRETQC_ALM_ENDPOINT=https://your.qcweb.server:9999"
        return
    else
        echo "PAKKRETQC_ALM_ENDPOINT=${PAKKRETQC_ALM_ENDPOINT}"
    fi
    go run devtools/cmd/appbundler/main.go -w
}
processFn
