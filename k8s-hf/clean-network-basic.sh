#!/bin/bash

kubectl delete -n hfbn configmap $(kubectl get -n hfbn configmap | awk '{if(NR>1) print $1}' | xargs)
kubectl delete -n hfbn pvc $(kubectl get -n hfbn pvc | awk '{if(NR>1) print $1}' | xargs)
kubectl delete -n hfbn pv $(kubectl get -n hfbn pv | awk '{if(NR>1) print $1}' | xargs)
kubectl delete namespace hfbn
