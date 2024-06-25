#!/bin/bash

# 切换到 simple 目录
cd ./simple

# 应用所有 YAML 文件
for file in *.yaml; do
  kubectl apply -f "$file"
done