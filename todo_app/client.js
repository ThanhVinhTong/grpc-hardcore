const grpc = require('@grpc/grpc-js');
const protoLoader = require("@grpc/proto-loader");
const axios = require('axios');
const { performance } = require('perf_hooks');

const packageDef = protoLoader.loadSync("./todo.proto", {});
const grpcObject = grpc.loadPackageDefinition(packageDef);
const todoPackage = grpcObject.todoPackage;
const grpcClient = new todoPackage.Todo("localhost:40000", grpc.credentials.createInsecure());

const REST_URL = "http://localhost:3000/todos";
const ITERATIONS = 10000;
const CONCURRENCY = 1000; // We send 50 requests at a time

async function benchmark() {
    console.log(`\n🚀 ADVANCED BENCHMARK: ${ITERATIONS} requests @ ${CONCURRENCY} concurrency`);
    console.log(`📡 Simulating 20ms network latency per request\n`);

    // --- REST THROUGHPUT ---
    const restStart = performance.now();
    for (let i = 0; i < ITERATIONS; i += CONCURRENCY) {
        const batch = [];
        for (let j = 0; j < CONCURRENCY && (i + j) < ITERATIONS; j++) {
            batch.push(axios.post(REST_URL, { title: "bench" }));
        }
        await Promise.all(batch);
    }
    const restDuration = (performance.now() - restStart) / 1000;
    const restRps = (ITERATIONS / restDuration).toFixed(2);

    // --- gRPC THROUGHPUT ---
    const grpcStart = performance.now();
    for (let i = 0; i < ITERATIONS; i += CONCURRENCY) {
        const batch = [];
        for (let j = 0; j < CONCURRENCY && (i + j) < ITERATIONS; j++) {
            batch.push(new Promise(r => grpcClient.createTodo({ title: "bench" }, r)));
        }
        await Promise.all(batch);
    }
    const grpcDuration = (performance.now() - grpcStart) / 1000;
    const grpcRps = (ITERATIONS / grpcDuration).toFixed(2);

    // --- gRPC STREAMING PERFORMANCE ---
    console.log(`\n⏳ Measuring Streaming (gRPC Benefit)...`);
    const streamStart = performance.now();
    let streamCount = 0;
    await new Promise((resolve) => {
        const stream = grpcClient.streamTodos({});
        stream.on('data', () => streamCount++);
        stream.on('end', resolve);
    });
    const streamDuration = performance.now() - streamStart;

    console.log("\n" + "=".repeat(40));
    console.log("FINAL PERFORMANCE RESULTS (THROUGHPUT)");
    console.log("=".repeat(40));
    console.table({
        "REST API": { "RPS (Req/Sec)": restRps, "Total Time": restDuration.toFixed(2) + "s" },
        "gRPC Unary": { "RPS (Req/Sec)": grpcRps, "Total Time": grpcDuration.toFixed(2) + "s" }
    });

    console.log(`\n📈 STREAMING ADVANTAGE:`);
    console.log(`Read ${streamCount} items via gRPC Stream in ${streamDuration.toFixed(2)}ms`);
    console.log(`(This avoided the 20ms per-request latency entirely!)`);

    process.exit(0);
}

benchmark().catch(err => console.error(err));