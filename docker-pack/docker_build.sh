#!/bin/bash
processFn() {
    if [ -z "${PAKKRETQC_ALM_ENDPOINT}"]; then
        echo "Error!! missing environment variable "
        echo "Example: "
        echo "  export PAKKRETQC_ALM_ENDPOINT=https://your.qcweb.server:9999"
        return
    else
        echo "PAKKRETQC_ALM_ENDPOINT=${PAKKRETQC_ALM_ENDPOINT}"
    fi
    jsvalue=($(jq -r '.version' ../package.json))
    echo "version = ${jsvalue[@]}"
    docker build --tag "pakkretqc:v$version" --build-arg PAKKRETQC_ALM_ENDPOINT=${PAKKRETQC_ALM_ENDPOINT} --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy .
}
processFn



