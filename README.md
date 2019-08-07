# NUMACC
**NUMA Configuration Checker for Container**  
NUMACC is a golang-based tool to check CPU affinity and NUMA configuration for containers and pods.  
NUMACC will indicates following information:   
1. Which CPU core that process runs on.  
2. Whether this process is pin to certain CPU core or not.      
3. NUMA node of net devices for given container/pod.  

Allowing user to know if the container/pod is under proper NUMA configuration.

## Requirement
go v1.10  

## Installation  
```shell script
git clone https://github.com/mjace/numacc
cd numacc
go get -d ./...
./numacc
```

## Usage 
To check NUMA configuration of container  
```shell script
./numacc cid <1234ABCFD>
```

## Example
```shell script
./numacc cid 6131f8f8cc
Process and CPU for Container  6131f8f8cc
PID	CurrentCpu	CpuAffinity
47431 	  43 	 0-71
68018 	  32 	 0-71
The NIC NUMA for container  6131f8f8cc
sriov-a	0
eth0	N/A
lo	N/A
```
