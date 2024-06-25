#!/bin/bash

parameters=(
  "310Mi" "1" "2Gi" "test1"
  "430Mi" "0.6" "3Gi" "test2"
  "370Mi" "0.9" "5Gi" "test3"
  "80Mi" "0.3" "1Gi" "test4"
  "700Mi" "1.7" "7Gi" "test5"
  "280Mi" "0.4" "5Gi" "test6"
  "520Mi" "0.2" "2Gi" "test7"
  "831Mi" "0.6" "8Gi" "test8"
  "463Mi" "0.9" "5Gi" "test9"
  "365Mi" "0.6" "2Gi" "test10"
  "421Mi" "1.1" "4Gi" "test11"
  "120Mi" "1" "1Gi" "test12"
  "157Mi" "0.7" "6Gi" "test13"
  "213Mi" "0.1" "4Gi" "test14"
  "586Mi" "0.3" "4Gi" "test15"
)

counter=0

for ((i=0; i<${#parameters[@]}; i+=4)); do
  memory=${parameters[i]}
  cpu=${parameters[i+1]}
  ephemeral_storage=${parameters[i+2]}
  deployment_name=${parameters[i+3]}

  cat > "./simple/${deployment_name}.yaml" <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: menet
  name: ${deployment_name}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: stress-job
  template:
    metadata:
      labels:
        app: stress-job
    spec:
      containers:
        - name: job
          image: polinux/stress
          resources:
            requests:
              memory: ${memory}
              cpu: ${cpu}
              ephemeral-storage: ${ephemeral_storage}
          command: ["/bin/sh", "-c"]
          args: ["stress --vm 1 --hdd 1 --hdd-bytes ${ephemeral_storage%?} --vm-bytes ${memory%?} --vm-hang 1"]
      restartPolicy: Always
EOF
done
