#!/usr/bin/env bash

kubectl apply -f examples/networkconnectivity/networkconnectivity-crd.yaml
kubectl apply -f examples/networkmonitor/networkmonitor-crd.yaml
kubectl apply -f examples/networknotification/networknotification-crd.yaml
kubectl apply -f examples/networktrafficshaper/networktrafficshaper-crd.yaml