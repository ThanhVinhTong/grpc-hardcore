# Engineering Report: High-Concurrency gRPC vs REST Benchmark

This report documents a high-stress performance comparison between gRPC and REST (Express) under simulated real-world conditions.

## 1. Test Methodology (The "Real-World" Setup)
To move beyond "localhost" trivialities, we implemented three critical constraints:
1.  **Network Latency Simulation:** Every request is forced to wait **20ms** before being processed, simulating a real-world network hop.
2.  **Concurrency Stress:** We tested scaling from 50 to **1,000 parallel requests**.
3.  **Streaming vs Unary:** We compared traditional Request/Response (REST) against gRPC Server Streaming.

## 2. Unary Throughput Results (RPS)
As we increased parallel load, gRPC's HTTP/2 multiplexing allowed it to handle significantly more traffic than REST.

| Concurrency Level | REST (RPS) | gRPC (RPS) | gRPC Multiplier |
| :--- | :--- | :--- | :--- |
| **50 Parallel** | 1,003 | 1,136 | 1.1x |
| **100 Parallel** | 2,071 | 2,196 | 1.1x |
| **1,000 Parallel** | 2,423 | **5,358** | **2.2x Faster** |

### **The "Concurrency Wall"**
Notice how REST starts to plateau around **2,400 Requests Per Second**. This is because HTTP/1.1 (used by REST) struggles to manage the overhead of thousands of simultaneous socket connections. gRPC, using a single multiplexed HTTP/2 connection, scaled to over **5,300 RPS** effortlessly.

---

## 3. The "Streaming Destruction"
The most significant finding was the difference between fetching large datasets via individual requests vs. a gRPC stream.

*   **Scenario:** Reading ~41,000 items.
*   **Result:** gRPC Stream finished in **385ms**.
*   **Theoretical REST Time:** If these items were fetched individually with 20ms latency, it would take **820 seconds** (13 minutes). 

**gRPC Streaming avoided the per-request network penalty entirely.**

---

## 4. Final Verdict

### Why gRPC Won:
1.  **Multiplexing (HTTP/2):** Many requests over one TCP connection reduced the overhead of opening 1,000 separate sockets.
2.  **Binary Serialization (Protobuf):** Small binary packets outperformed verbose JSON strings as payload volume increased.
3.  **Persistent Streams:** The ability to "pipe" data rather than "request/respond" made the network latency irrelevant for bulk data transfers.

### When to choose which?
*   **Use REST** for public APIs, browser compatibility, and simple CRUD.
*   **Use gRPC** for high-performance internal microservices, real-time data streaming, and cross-language communication in high-traffic environments.
