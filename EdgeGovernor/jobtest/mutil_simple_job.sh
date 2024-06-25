#!/bin/bash

# 定义作业的基础路径
_job_base_path="/tmp/jobs"
mkdir -p "$_job_base_path"

_job_index=1

create_job() {
  local job_name="workload-$((_job_index + 1))"  # 使用局部变量和正确的算术增加
  local job_file_path="${_job_base_path}/${job_name}.yaml"

  # 定义资源请求参数，这里直接作为函数参数传入
  local request_memory=$1
  local request_cpu=$2
  local request_storage=$3

  # 创建 YAML 文件内容
  cat <<EOF > "$job_file_path"
apiVersion: v1
kind: Pod
metadata:
  namespace: run-job
  name: \${job_name}
spec:
  containers:
  - name: job
    image: polinux/stress
    resources:
      requests:
        memory: "$request_memory"
        cpu: "$request_cpu"
        ephemeral-storage: "$request_storage"
    command: ["stress"]
    args: ["--vm", "1", "--vm-bytes", "$request_memory", "--vm-hang", "0"]
  restartPolicy: Never
EOF

  echo "Created job $job_name as $job_file_path"
  ((_job_index++))  # 确保 _job_index 增加
}

# 删除旧的作业文件并创建新的作业
#rm -f "${_job_base_path}"/*
create_job "400Mi" "1" "3Gi"
create_job "164Mi" "0.6" "5Gi"
create_job "562Mi" "0.9" "10Gi"
create_job "80Mi" "0.3" "1Gi"
create_job "1000Mi" "1.7" "12Gi"
create_job "620Mi" "0.4" "5Gi"
create_job "70Mi" "0.2" "1Gi"
# ... （其他作业调用）

echo "[INFO] All jobs yaml files are available in $_job_base_path"
echo "[INFO] Deploying resources in the cluster..."

## 部署所有作业
#for job_file in "${_job_base_path}"/*.yaml; do
#  kubectl -n run-job apply -f "$job_file"
#done