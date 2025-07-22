from dataclasses import dataclass
import dataclasses
import json
from pathlib import Path
from neusight.Prediction.predictor import OperatorPredictor
from neusight.Tracing.parse import parse_trace
import pandas as pd
import ast
import rich
from neusight.Tracing.parse import parse_ops, parse_dependency, fuse_parse
from neusight.Prediction.predictor import NeusightPredictor

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
    
    def has_overhead(self):
        return self.forward_compute_time > 0 or self.backward_compute_time > 0 or self.weight_grad_compute_time > 0 or self.forward_comm_size > 0 or self.backward_comm_size > 0 or self.weight_grad_comm_size > 0

class LayerInfo:
    def __init__(self, layer_df):

        columns = ["Name", "OpName", "FwOps", "BwOps", "AccOps", "Prev", "Next", "InputShapes", "OutputShape", "fw_latency", "bw_latency", "acc_latency"]
        self.name = layer_df.iloc[columns.index("Name")]
        self.op_name = layer_df.iloc[columns.index("OpName")]
        self.fw_ops = layer_df.iloc[columns.index("FwOps")]
        self.bw_ops = layer_df.iloc[columns.index("BwOps")]
        self.acc_ops = layer_df.iloc[columns.index("AccOps")]
        self.prev = layer_df.iloc[columns.index("Prev")]
        self.next = layer_df.iloc[columns.index("Next")]
        self.input_shapes = layer_df.iloc[columns.index("InputShapes")]
        self.output_shape = layer_df.iloc[columns.index("OutputShape")]
        self.fw_latency = layer_df.iloc[columns.index("fw_latency")]
        self.bw_latency = layer_df.iloc[columns.index("bw_latency")]
        self.acc_latency = layer_df.iloc[columns.index("acc_latency")]
        self.set_comm()
        
    def set_comm(self):           
        # bw_ops example: [['ALLREDUCE', (B * S * H,)]]
        # fw_ops example: [['SENDRECV', (B * S * H,)]]
        if self.op_name in ['allreduce', 'sendrecv']:
            if self.fw_ops and len(self.fw_ops) > 0 and len(self.fw_ops[0]) > 1:
                self.fw_comm_type = self.fw_ops[0][0]
                self.fw_comm_size = self.fw_ops[0][1][0]
            else:
                self.fw_comm_size = 0
                self.fw_comm_type = "None"
            if self.bw_ops and len(self.bw_ops) > 0 and len(self.bw_ops[0]) > 1:
                self.bw_comm_type = self.bw_ops[0][0]
                self.bw_comm_size = self.bw_ops[0][1][0]
            else:
                self.bw_comm_type = "None"
                self.bw_comm_size = 0
            self.fw_latency = 0
            self.bw_latency = 0
            self.acc_latency = 0
        else:
            self.fw_comm_type = "None"
            self.fw_comm_size = 0
            self.bw_comm_type = "None"
            self.bw_comm_size = 0
    
    def to_work_item(self) -> Work_Item:
        return Work_Item(
            name=self.name, 
            forward_compute_time=self.fw_latency,
            forward_comm_size=self.fw_comm_size,
            forward_comm_type=self.fw_comm_type,
            backward_compute_time=self.bw_latency,
            backward_comm_size=self.bw_comm_size,
            backward_comm_type=self.bw_comm_type,
            weight_grad_compute_time=self.acc_latency,
            weight_grad_comm_size=0,
            weight_grad_comm_type="None"
        )

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
        self.workloads = []
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
        # columns = ["Name", "OpName", "FwOps", "BwOps", "AccOps", "Prev", "Next", "InputShapes", "OutputShape"]
        # df = df[columns]
        
        for _, row in self.df.iterrows():
            layer_info = LayerInfo(row)
            item = layer_info.to_work_item()
            if item.has_overhead():
                self.workloads.append(item)
    
    def dump_file(self, filename):
        # 输出 txt 文件
        txt_filename = filename + ".txt"
        with open(txt_filename, "w") as f:
            columns = list(Work_Item.__dataclass_fields__.keys())
            f.write("\t".join(columns) + "\n")
            for item in self.workloads:
                f.write(
                    "\t".join([str(getattr(item, k)) for k in columns])
                    + "\n"
                )
        
        # 输出 json 文件
        json_filename = filename + ".json"
        with open(json_filename, "w") as f:
            json_data = []
            for item in self.workloads:
                item_dict = dataclasses.asdict(item)
                json_data.append(item_dict)
            json.dump(json_data, f, indent=2, ensure_ascii=False)
        
        rich.print(f"fwd_latency_e2e: {sum(w.forward_compute_time for w in self.workloads)}")
        rich.print(f"已输出文件: {txt_filename}, {json_filename}")

    def generate_gpu_intensity(self):
        pass