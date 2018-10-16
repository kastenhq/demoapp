#! /usr/bin/bash

GOPATH=$(pwd)/go swagger generate model --spec=swagger.yaml -t go/src/demoapp

GOPATH=$(pwd)/go swagger generate client --spec=swagger.yaml -t go/src/demoapp \
            --name=rest \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --template-dir=build/templates