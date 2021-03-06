# Demoapp

[![Go Report Card](https://goreportcard.com/badge/github.com/kastenhq/demoapp)](https://goreportcard.com/badge/github.com/kastenhq/demoapp)

## Overview
This is a Demo app. can be used to demo or testing purposes.
It's a kubernetes-native picture gallery.

### Design Goals

This application was created as an example for Kubernets-native application

Why this app is kubernetes-native:

1. The application has 2 backend microservices metadata and store.
2. These services are written in golang.
3. This application follows API-first pattern.
4. The application has swagger generated servers and clients.

## Quick Start

### Pre-requisites

- docker
- golang
- helm

### How To Build

Just run build.sh

```
bash build.sh
```
This script will run swagger, glide, go install, go test and build docker images.

```
~/demoapp$ docker images
REPOSITORY                                               TAG                   IMAGE ID            CREATED              SIZE
frontend-server                                          latest                7ffa0c3358a6        About a minute ago   109MB
metadata-server                                          latest                a0b50306f58e        About a minute ago   15.8MB
store-server                                             latest                42ed98766f26        About a minute ago   14.1MB
```

### How To Just Deploy Without Build

```
helm install helm/demoapp/ --name demoapp --namespace demoapp
```

This chart will deploy `demoapp`. Alongside with it, MongoDB will be deployed and configured.
Also `all-in-one jaeger` will be deployed.

#### How To Access

Get access by default `demoapp` will be exposed via cluster ingress controller.

To get access to `jaeger`

```
kubectl port-forward svc/demoapp-jaeger-query -n kube-system 8001:80
```
then open http://127.0.0.1:8001

You also can deploy app with internal ingress controller

```
helm install helm/demoapp/ --name demoapp --namespace demoapp --set nginx.enabled=true
```
Then you should be able to `kubectl port-forward` to nginx service


### How to add photos

!For now, picture gallery works with `png` only!

```
for i in $(ls ./my_pictures); do (echo -n '{ "base64": "'"$(base64 -w 0 ./my_pictures/$i)"'"}')|curl -H "Content-Type: application/json" -d @- https://mytestcluster/demoapp/metadata/v0/images ; done
```
