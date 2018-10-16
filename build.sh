#! /usr/bin/bash

GOPATH=$(pwd)/go swagger generate model --spec=swagger.yaml -t go/src/demoapp

GOPATH=$(pwd)/go swagger generate client --spec=swagger.yaml -t go/src/demoapp \
            --name=rest \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --template-dir=build/templates

GOPATH=$(pwd)/go swagger generate server --spec=swagger.yaml -t go/src/demoapp \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --name=store \
            --server-package=storeserver \
            --operation=healthz \
            --operation=deleteImageData \
            --operation=getImageData

GOPATH=$(pwd)/go swagger generate server --spec=swagger.yaml -t go/src/demoapp \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --name=meta \
            --server-package=metaserver \
            --operation=healthz \
            --operation=addImage \
            --operation=listImages \
            --operation=getImage \
            --operation=deleteImage
