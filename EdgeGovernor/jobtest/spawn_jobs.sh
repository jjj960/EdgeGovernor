#!/bin/bash

_jobs_file_path="/tmp/jobs.yaml"
_job_template_file_name="job-template.yaml"
_job_index=1

create_job() {
  _job_name="workload-$_job_index"
  _request_memory=$2
  _request_memory_in_mb="${_request_memory%Mi}M"
  _request_cpu=$3
  _request_storage=$4

  sed -e "s/JOB_NAME_PLACEHOLDER/$_job_name/g" \
      -e "s/MEMORY_PLACEHOLDER/$_request_memory/g" \
      -e "s/MEMORY_STRESS_PLACEHOLDER/$_request_memory_in_mb/g" \
      -e "s/CPU_PLACEHOLDER/$_request_cpu/g" \
      -e "s/EPHEMERAL_STORAGE_PLACEHOLDER/$_request_storage/g" \
      $_job_template_file_name >> $_jobs_file_path

  ((_job_index+=1))
}

#rm $_jobs_file_path
create_job "400Mi" "1" "3Gi"
create_job "164Mi" "0.6" "5Gi"
create_job "562Mi" "0.9" "10Gi"
create_job "80Mi" "0.3" "1Gi"
create_job "1000Mi" "1.7" "12Gi"
create_job "620Mi" "0.4" "5Gi"
create_job "70Mi" "0.2" "1Gi"
#create_job "60s" "1000Mi" "1" "8Gi"
#create_job "54s" "382Mi" "1.2" "10Gi"
#create_job "27s" "1347Mi" "0.7" "9Gi"
#create_job "78s" "1719Mi" "1.8" "10Gi"
#create_job "36s" "2185Mi" "1.7" "5Gi"
#create_job "92s" "952Mi" "0.9" "9Gi"
#create_job "15s" "1264Mi" "1.2" "9Gi"
#create_job "63s" "1891Mi" "1.3" "10Gi"
#create_job "42s" "2635Mi" "1.2" "6Gi"
#create_job "70s" "557Mi" "1.6" "8Gi"
#create_job "10s" "774Mi" "1.0" "6Gi"

echo "[INFO] The jobs yaml file is available at $_jobs_file_path"
echo "[INFO] Deploying resources in the cluster..."

#kubectl -n menet apply -f $_jobs_file_path