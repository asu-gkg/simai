from dataclasses import dataclass
import dataclasses
import json
from pathlib import Path
from neusight.Prediction.predictor import OperatorPredictor
from neusight.Tracing.parse import parse_trace
import pandas as pd
import ast
from neusight.Tracing.parse import parse_ops, parse_dependency, fuse_parse
from neusight.Prediction.predictor import NeusightPredictor

class LayerInfo:
    def __init__(self, layer_df):
        columns = ["Name", "OpName", "FwOps", "BwOps", "AccOps", "Prev", "Next", "InputShapes", "OutputShape"]
        self.name = layer_df.iloc[columns.index("Name")]
        self.op_name = layer_df.iloc[columns.index("OpName")]
        self.fw_ops = layer_df.iloc[columns.index("FwOps")]
        self.bw_ops = layer_df.iloc[columns.index("BwOps")]
        self.acc_ops = layer_df.iloc[columns.index("AccOps")]
        self.prev = layer_df.iloc[columns.index("Prev")]
        self.next = layer_df.iloc[columns.index("Next")]
        self.input_shapes = layer_df.iloc[columns.index("InputShapes")]
        self.output_shape = layer_df.iloc[columns.index("OutputShape")]
    


@dataclass
class Work_Item:
    name: str = dataclasses.field(default="none")
    forward_compute_time: int = dataclasses.field(default=0)
    forward_comm_type: str = dataclasses.field(default="NONE")
    forward_comm_size: int = dataclasses.field(default=0)
    backward_compute_time: int = dataclasses.field(default=0)
    backward_comm_type: str = dataclasses.field(default="NONE")
    backward_comm_size: int = dataclasses.field(default=0)
    weight_grad_compute_time: int = dataclasses.field(default=0)
    weight_grad_comm_type: str = dataclasses.field(default="NONE")
    weight_grad_comm_size: int = dataclasses.field(default=0)
    process_time: int = dataclasses.field(default=0)


class WorkloadGenerator:
    def __init__(self, csv_file, parallel_config, device_config_path, 
                 fusion=True):
        self.csv_file = csv_file
        self.parallel_config = parallel_config
        self.dp = parallel_config.dp
        self.tp = parallel_config.tp
        self.pp = parallel_config.pp
        self.pp_num_microbatch = parallel_config.pp_num_microbatch
        self.is_distributed = parallel_config.is_distributed()
        self.fusion = fusion
        self.df = None
        device_config_path = Path(device_config_path)
        device_config_path = device_config_path.absolute()
        with open(device_config_path, "r") as f:
            self.device_config = json.load(f)
        self.predictor = OperatorPredictor(predictor_path="./mtccl/NeuSight/scripts/asplos/data/predictor/MLP_WAVE", 
                    tile_dataset_dir=Path("./mtccl/NeuSight/scripts/asplos/data/dataset/train"))
        self.format_csv()
        
    def format_csv(self, is_train=False):
        # 当前的dist模式dp tp pp是互斥的，即只能有一种并行模式
        df = parse_trace(
                    self.csv_file, 
                    is_train=is_train, 
                    bench=False, 
                    fusion=self.fusion,
                    distributed=self.is_distributed,
                    dp_degree=self.dp,
                    pp_degree=self.pp,
                    pp_num_microbatch=self.pp_num_microbatch,
                    tp_degree=self.tp,
        )
        
        # ms
        df[[f"fw_latency", f"bw_latency", f"acc_latency"]] = df.apply(lambda x: self.predictor.predict(self.device_config, x), axis=1)
        self.df = df
    
    def generate_workload(self):
        columns = ["Name", "OpName", "FwOps", "BwOps", "AccOps", "Prev", "Next", "InputShapes", "OutputShape"]
        df = df[columns]
        
        for index, row in df.iterrows():
            layer_info = LayerInfo(row)
        # todo add work_item to workloads

    def generate_gpu_intensity(self):
        pass