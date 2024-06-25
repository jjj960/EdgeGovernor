# EdgeGovernor
EdgeGovernor is an IoT framework used for cloud edge collaboration and supporting edge autonomy. Compared with traditional cloud edge collaborative systems, EdgeGovernor enhances fault tolerance and adaptability.

The demo video can be found at https://youtu.be/vXz17xd9SW0.

(1) EdgeGovernor endows IoT systems with the capability of edge cluster autonomy, ensuring cluster management even in the missing of cloud node failure. 

(2) Some built-in task scheduling and resource management algorithms are designed in a hot-swappable manner, allowing users to replace algorithms according to actual requirements. 

(3) EdgeGovernor exhibits strong compatibility and can manage numerous edge computing platforms developed based on Kubernetes. 

## Framework architecture

![arch](F:\Desktop\ASE 2024\code\EdgeGovernor\figures\arch.png)

In the cloud-edge collaborative and edge-autonomous environments, to achieve efficient management of cluster workflows, the EdgeGovernor framework utilizes a strategy based on the leader-follower pattern and the designation of candidate node roles. The framework's architecture is composed of multiple critical layers and services, aimed at fulfilling the diverse needs of various task scenarios while ensuring the stability of the framework.

The edge computing infrastructure layer serves as the cornerstone of the EdgeGovernor architecture, jointly composed of edge servers and edge devices. EdgeGovernor is responsible for collecting various types of information from both edge servers and edge devices. These edge servers are all equipped with Docker and Kubernetes environments to support the flexible deployment of tasks in diverse scenarios. To ensure the security and reliability of data in the event of a cluster node failure, EdgeGovernor requires that the databases of all nodes form a cluster, enabling real-time data backup and ensuring the high availability of the framework.

The data privacy layer uses AES 256-CFB symmetric encryption to secure data from the host and node-to-node communication. This approach eschews common asymmetric encryption due to the centralization risks it may bring, which could undermine system security during cloud server attacks. The CFB mode's benefit is encryption without enlarging data packets, and as a stream cipher, it can detect tampering. 

The data storage layer is responsible for storing all data generated by EdgeGovernor, where key service data and task data are backed up on multiple nodes to ensure data integrity and maintain cluster state consistency in the event of node failures.


The service layer, a crucial part of EdgeGovernor, is responsible for managing cluster operations. The node management service, cluster state management service, and data persistence service work together closely to ensure that the edge cluster is still able to achieve autonomy even in the event of cloud node failures. The task management service is responsible for managing the lifecycle of all tasks within the cluster, including task creation, assignment, and monitoring. Additionally, the task management service works in conjunction with the scheduler, leveraging the scheduler's efficient resource management algorithms to further ensure the optimized allocation of cluster resources and the high efficiency of task execution.

## Framework Evaluation

To evaluate the performance of EdgeGovernor in the autonomous management of edge clusters, we constructed a full-physical simulation cluster based on a cloud-edge collaborative architecture, composed of Raspberry Pi and Nvidia TX2 devices, with specific configurations as shown in Figure 3. The cluster architecture consists of one cloud node and three edge nodes. We simulated node failure scenarios by randomly starting and stopping devices. The experimental results show that in the case of disconnection between the cloud node and the system, EdgeGovernor can still manage and control the edge cluster stably, maintain data consistency, and ensure the normal execution of tasks. Additionally, we conducted a performance evaluation of EdgeGovernor's task scheduling algorithm using the cluster's comprehensive resource standard deviation metric. A lower value of this metric indicates a more uniform resource distribution. The results showed that compared to Kubernetes (K8s) LeastRequestedPriority and BalancedResourceAllocation algorithms, EdgeGovernor reduced the resource standard deviation by 2.62% and 5.95%, respectively.

![test1](F:\Desktop\ASE 2024\code\EdgeGovernor\figures\test1.jpg)

| Device            | Role  | CPU (core) | Memory (GB) | Bandwidth (Gb/s) | Disk (GB) |      |
| ----------------- | ----- | ---------- | ----------- | ---------------- | --------- | ---- |
| ThinkPad P52s     | Cloud | 4          | 16          | 0.1              | 512       |      |
| Raspberry Pi 4B   | Edge1 |            | 8           | 0.1              | 64        |      |
| Raspberry Pi 3B   | Edge2 |            |             | 0.1              | 16        |      |
| Nvidia Jetson TX2 | Edge3 |            |             | 0.1              | 32        |      |

