from typing import Dict, Tuple
from typing import List

class FatTree:
    num_gpus: int
    gpu_type: str
    num_servers: int 
    num_tors: int
    num_spines: int
    num_aggs: int
    num_switches: int
    gpus_per_server: int
    servers_per_tor: int
    tors_per_spine: int
    spines_per_agg: int
    link_bandwidth: Dict[str, float] # in Gb/s
    link_latency: Dict[str, float] # in us
    nodes: List[str] # gpu, server, spine, tor, agg
    edges: List[Tuple[str, str, float, float]] # (src, dst, bandwidth, latency)
    
    
    def __init__(self, num_gpus: int, 
                num_servers: int, 
                num_spines: int, 
                gpus_per_server: int, 
                servers_per_tor: int, 
                tors_per_spine: int, 
                spines_per_agg: int, 
                link_bandwidth: Dict[str, float] = None,
                link_latency: Dict[str, float] = None):
        self.num_gpus = num_gpus
        self.num_servers = num_servers
        self.num_spines = num_spines
        self.gpus_per_server = gpus_per_server
        self.servers_per_tor = servers_per_tor
        self.tors_per_spine = tors_per_spine
        self.spines_per_agg = spines_per_agg
        self.link_bandwidth = link_bandwidth
        self.link_latency = link_latency
        self.nodes = []
        self.edges = []
        self.capacity = {}
        self.latency = {}
        
        if self.link_bandwidth is None:
            self.link_bandwidth = {
                "l0": 100, # server to tor
                "l1": 100, # tor to spine
                "intra-gpu": 512, # intra-gpu
            }
        if self.link_latency is None:
            self.link_latency = {
                "l0": 2.0, # server to tor
                "l1": 2.0, # tor to spine
                "intra-gpu": 0.3, # intra-gpu
            }
        self._build()
            
    def get_link_bandwidth(self, link_type: str) -> float:
        return self.link_bandwidth[link_type]
    
    def get_link_latency(self, link_type: str) -> float:
        return self.link_latency[link_type]
    
    def get_num_gpus(self) -> int:
        return self.num_gpus
    
    def get_num_servers(self) -> int:
        return self.num_servers
    
    def get_num_spines(self) -> int:
        return self.num_spines
    
    def _server_id(self, gpu_id: int) -> int:
        return gpu_id % self.num_servers
    
    def _add_link(self, src, dst, bw, lat):
        self.edges.append((src, dst))
        self.capacity[(src, dst)] = bw
        self.latency[(src, dst)] = lat
        
    def _build(self):
        for g in range(self.num_gpus):
            self.nodes.append(f"gpu{g}")
        for s in range(self.num_servers):
            self.nodes.append(f"server{s}")
        for t in range(self.num_tors):
            self.nodes.append(f"tor{t}")
        for a in range(self.num_aggs):
            self.nodes.append(f"agg{a}")
        for s in range(self.num_spines):
            self.nodes.append(f"spine{s}")
        
        # GPU <-> GPU (intra-server)
        for s in range(self.num_servers):
            gpu_ids = [s * self.gpus_per_server + i for i in range(self.gpus_per_server)]
            for i in gpu_ids:
                for j in gpu_ids:
                    if i == j:
                        continue
                    self._add_link(f"gpu{i}", f"gpu{j}", self.link_bandwidth["intra-gpu"], self.link_latency["intra-gpu"])
        
        # server <-> ToR
        for s in range(self.num_servers):
            tor_id = s // self.servers_per_tor
            self._add_link(f"server{s}", f"tor{tor_id}", self.link_bandwidth["l0"], self.link_latency["l0"])
            self._add_link(f"tor{tor_id}", f"server{s}", self.link_bandwidth["l0"], self.link_latency["l0"])
        
        # ToR <-> Spine
        for t in range(self.num_tors):
            for i in range(self.tors_per_spine):
                spine_id = (t * self.tors_per_spine + i) % self.num_spines
                self._add_link(f"tor{t}", f"spine{spine_id}", self.get_link_bandwidth("l1"), self.get_link_latency("l1"))
                self._add_link(f"spine{spine_id}", f"tor{t}", self.get_link_bandwidth("l1"), self.get_link_latency("l1"))

        # Spine <-> Agg 
        for s in range(self.num_spines):
            agg_id = s // self.spines_per_agg
            self._add_link(f"spine{s}", f"agg{agg_id}", self.get_link_bandwidth("l1"), self.get_link_latency("l1"))
            self._add_link(f"agg{agg_id}", f"spine{s}", self.get_link_bandwidth("l1"), self.get_link_latency("l1"))