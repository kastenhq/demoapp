#! /bin/bash

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
            --operation=getImageData \
            --operation=storeImageData \

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

GOPATH=${BASE}/go glide install --strip-vendor

GOPATH=${BASE}/go CGO_ENABLED=0 GO_EXTLINK_ENABLED=0 go install -v  ./...
popd

BIN=store-server envsubst < build/templates/Dockerfile | docker build -t store-server go/bin -f -
BIN=meta-server envsubst < build/templates/Dockerfile | docker build -t metadata-server go/bin -f -
docker build -t frontend-server frontend -f frontend/Dockerfile

GOPATH=${BASE}/go go test -v ./... -check.vv
