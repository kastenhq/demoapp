#! /usr/bin/bash

BASE=$(pwd)

pushd $(pwd)/go/src
GOPATH=${BASE}/go swagger generate model --spec=${BASE}/swagger.yaml

GOPATH=${BASE}/go swagger generate client --spec=${BASE}/swagger.yaml \
            --name=rest \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --template-dir=${BASE}/build/templates

GOPATH=${BASE}/go swagger generate server --spec=${BASE}/swagger.yaml \
            --skip-models \
            --skip-validation \
            --client-package=restclient \
            --name=store \
            --server-package=storeserver \
            --operation=healthz \
            --operation=deleteImageData \
            --operation=getImageData

GOPATH=${BASE}/go swagger generate server --spec=${BASE}/swagger.yaml \
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

GOPATH=${BASE}/go glide update -v

GOPATH=${BASE}/go go test -v  ./...

GOPATH=${BASE}/go go install -v  ./...

popd
