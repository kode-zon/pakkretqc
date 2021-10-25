#!/bin/bash
PAKKRETQC_ALM_ENDPOINT=https://qcweb.kasikornbank.com:8443
processFn() {
    if [ -z "${PAKKRETQC_ALM_ENDPOINT}"]; then
        echo "Error!! missing environment variable "
        echo "Example: "
        echo "  export PAKKRETQC_ALM_ENDPOINT=https://your.qcweb.server:9999"
        return
    else
        echo "PAKKRETQC_ALM_ENDPOINT=${PAKKRETQC_ALM_ENDPOINT}"
    fi
    version="lastest"
    if ! command -v jq &> /dev/null
    then
        echo "command jq doesn't exist"
    else
        jsvalue=($(jq -r '.version' ../package.json))
        echo "version = ${jsvalue[@]}"
        version=${jsvalue[@]}
    fi

    echo "output docker will tag as [pakkretqc:${version}]"

    docker build --tag "pakkretqc:${version}" \
        --network=host \
        --build-arg PAKKRETQC_ALM_ENDPOINT=${PAKKRETQC_ALM_ENDPOINT} \
        --build-arg http_proxy="http://127.0.0.1:3129" \
        --build-arg https_proxy="http://127.0.0.1:3129" \
        --build-arg no_proxy="localhost,127.0.0.1,.kasikornbank.com,::1" \
        -f private.Dockerfile .
}
processFn



