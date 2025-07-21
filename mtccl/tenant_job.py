# class TenantJob:
#     job_id: str
#     arrival_time: int # in nanoseconds
#     parallel_config: ParallelConfig
   
from dataclasses import dataclass
import os
import json
import argparse
from typing import Dict, List, Optional, Tuple
import rich
from neusight.Tracing.trace import get_model
from neusight.Prediction.predictor import dump_df
from workload_generator import WorkloadGenerator
from parse import show_df
import pandas as pd

@dataclass
class ParallelConfig:
    """并行配置类"""
    dp: int = 1  # Data parallel degree
    tp: int = 1  # Tensor parallel degree  
    pp: int = 1  # Pipeline parallel degree
    pp_num_microbatch: int = 1  # Pipeline parallel micro-batch size
    
    def __init__(self, dp: int, pp: int, tp: int, pp_num_microbatch: int):
        self.dp = dp
        self.tp = tp
        self.pp = pp
        self.pp_num_microbatch = pp_num_microbatch
    
    def total_gpus(self) -> int:
        """计算总GPU数量"""
        return self.dp * self.tp * self.pp
    
    def __str__(self) -> str:
        """字符串表示"""
        return f"DP={self.dp}_TP={self.tp}_PP={self.pp}_MB={self.pp_num_microbatch}"
    
    def to_neusight_options(self) -> str:
        """转换为neusight选项字符串"""
        options = []
        if self.dp > 1:
            options.append(f"dp{self.dp}")
        if self.tp > 1:
            options.append(f"tp{self.tp}")
        if self.pp > 1:
            options.append(f"pp{self.pp}_{self.pp_num_microbatch}")
        return "_".join(options) if options else ""
    
    def is_distributed(self) -> bool:
        """判断是否为分布式配置"""
        return self.dp > 1 or self.tp > 1 or self.pp > 1
    
    def get_micro_batch_size(self, global_batch_size: int) -> int:
        """计算micro-batch size"""
        if self.dp > 1:
            return global_batch_size // self.dp
        elif self.pp > 1:
            return global_batch_size // self.pp_num_microbatch
        else:
            return global_batch_size

@dataclass
class TrainingWorkload:
    batch_size: int
    sequence_length: int
    model_config_path: str
    device_config_path: str
    fusion: bool
    execution_type: str = "train"
    
    def __init__(self, batch_size: int,
                sequence_length: int, 
                model_config_path: str, 
                device_config_path: str, 
                execution_type: str = "train", 
                fusion: bool = True):
        self.batch_size = batch_size
        self.sequence_length = sequence_length
        self.model_config_path = model_config_path
        self.device_config_path = device_config_path
        self.execution_type = execution_type
        self.fusion = fusion

@dataclass
class Tenant:
    tenant_id: str
    gpu_ids: List[int]
    
    def __init__(self, tenant_id: str, gpu_ids: List[int]):
        self.tenant_id = tenant_id
        self.gpu_ids = gpu_ids
        
    def get_gpu_ids(self) -> List[int]:
        return self.gpu_ids
    
    def get_tenant_id(self) -> str:
        return self.tenant_id

class LayerComputationDelay:
    """每层计算延迟信息"""
    layer_id: int
    forward_latency: float  # ms
    backward_latency: float  # ms
    communication_latency: float  # ms
    accumulation_latency: float  # ms
    total_latency: float  # ms
    
    def to_dict(self) -> Dict:
        return {
            "layer_id": self.layer_id,
            "forward_latency": self.forward_latency,
            "backward_latency": self.backward_latency,
            "communication_latency": self.communication_latency,
            "accumulation_latency": self.accumulation_latency,
            "total_latency": self.total_latency
        }

class CommunicationVolume:
    """通信量信息"""
    comm_type: str  # "AllGather", "AllReduce", "P2P" 
    size_gb: float
    group_size: int
    frequency: str  # "per_step", "per_layer" 
    
    def to_dict(self) -> Dict:
        return {
            "comm_type": self.comm_type,
            "size_gb": self.size_gb,
            "group_size": self.group_size,
            "frequency": self.frequency
        }

class TenantJob:
    job_id: str
    arrival_time: int # in nanoseconds
    tenant_id: str
    gpu_ids: List[int]
    gpu_type: str
    parallel_config: Optional[ParallelConfig]
    workload: Optional[TrainingWorkload]
    
    def __init__(self, job_id: str, arrival_time: int,
                tenant_id: str, 
                gpu_ids: List[int], 
                gpu_type: str, 
                parallel_config: Optional[ParallelConfig] = None, 
                workload: Optional[TrainingWorkload] = None):
        self.job_id = job_id
        self.arrival_time = arrival_time
        self.tenant_id = tenant_id
        self.gpu_ids = gpu_ids
        self.gpu_type = gpu_type
        self.parallel_config = parallel_config
        self.workload = workload
        
        self._layer_delays: Optional[List[LayerComputationDelay]] = None
        self._comm_volumes: Optional[List[CommunicationVolume]] = None
        self._neusight_results: Optional[Dict] = None

    def get_job_id(self) -> str:
        return self.job_id
    
    def get_parallel_config(self) -> ParallelConfig:
        return self.parallel_config
    
    def get_workload(self) -> TrainingWorkload:
        return self.workload
    
    def get_tenant_id(self) -> str:
        return self.tenant_id
    
    def get_gpu_ids(self) -> List[int]:
        return self.gpu_ids
    
    def get_computation_graph(self) -> Tuple[List[LayerComputationDelay], object]:
        if not self.workload:
            raise ValueError("Workload is not set")
        
        from neusight.Tracing.trace import trace_graph
        
        df, _ = trace_graph(
                    model_config_path=self.workload.model_config_path, 
                    sequence_length=self.workload.sequence_length, 
                    batch_size=self.workload.batch_size, 
                    is_train=True, 
                    bench=False,
                    single_layer=True, 
                    fusion=self.workload.fusion,
                    distributed=self.parallel_config.is_distributed(),
                    dp_degree=self.parallel_config.dp,
                    pp_degree=self.parallel_config.pp,
                    pp_num_microbatch=self.parallel_config.pp_num_microbatch,
                    tp_degree=self.parallel_config.tp,
        )
        
        return df

def load_tenant_job(job_data: Dict) -> TenantJob:
    model_config_path = job_data["model_config_path"]
    device_config_path = job_data["device_config_path"]
    tenant_job = TenantJob(
        job_id=job_data["job_id"],
        arrival_time=job_data["arrival_time"],
        tenant_id=job_data["tenant_id"],
        gpu_ids=job_data["gpu_ids"],
        gpu_type=job_data["gpu_type"],
        parallel_config=ParallelConfig(
            dp=job_data["parallel_config"]["dp"],
            pp=job_data["parallel_config"]["pp"],
            tp=job_data["parallel_config"]["tp"],
            pp_num_microbatch=job_data["parallel_config"]["pp_num_microbatch"],
        ),
        workload=TrainingWorkload(
            batch_size=job_data["workload"]["batch_size"],
            sequence_length=job_data["workload"]["sequence_length"],
            model_config_path=model_config_path,
            device_config_path=device_config_path,
            fusion=False,
        ),
    )
    return tenant_job


        
if __name__ == "__main__":
    parse = argparse.ArgumentParser()
    parse.add_argument("--jobs_dir", type=str, default="inputs/workloads/clos1_jobs")
    args = parse.parse_args()
    jobs_dir = args.jobs_dir
    tenant_jobs = []
    for job_file in os.listdir(jobs_dir):
        if job_file.endswith(".json"):
            with open(os.path.join(jobs_dir, job_file), "r") as f:
                job_data = json.load(f)
                tenant_job = load_tenant_job(job_data)
                tenant_jobs.append(tenant_job)
    
    # todo: add for loop
    tenant_job = tenant_jobs[1]
    rich.inspect(tenant_job)
    csv_file = jobs_dir + f"/layer_delays_{tenant_job.job_id}.csv"
    if not os.path.exists(csv_file):
        layer_delays = tenant_job.get_computation_graph()
        csv_file = dump_df(layer_delays, jobs_dir + f"/layer_delays_{tenant_job.job_id}.csv")
    
    # model, n_layer = get_model(model_config_path = tenant_job.workload.model_config_path,
    #                 is_train=True, 
    #                 device="cuda",
    #                 fusion=tenant_job.workload.fusion,
    #                 )
    # rich.inspect(model)
    # rich.inspect(n_layer)
    
    wg = WorkloadGenerator(csv_file, parallel_config=tenant_job.parallel_config, 
                           device_config_path=tenant_job.workload.device_config_path,
                           fusion=tenant_job.workload.fusion,
    )
    
    show_df(wg.df)