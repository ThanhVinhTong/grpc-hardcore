const grpc = require('@grpc/grpc-js');
const protoLoader = require("@grpc/proto-loader");

const packageDef = protoLoader.loadSync("./todo.proto", {});
const grpcObject = grpc.loadPackageDefinition(packageDef);
const todoPackage = grpcObject.todoPackage;

const server = new grpc.Server();
const todos = [];

// SIMULATE NETWORK LATENCY (ms)
const ARTIFICIAL_LATENCY = 20; 

server.addService(todoPackage.Todo.service, {
    "createTodo": (call, callback) => {
        setTimeout(() => {
            const todoItem = { ...call.request, id: todos.length + 1 };
            todos.push(todoItem);
            callback(null, todoItem);
        }, ARTIFICIAL_LATENCY);
    },
    "readTodos": (call, callback) => {
        setTimeout(() => {
            callback(null, { items: todos });
        }, ARTIFICIAL_LATENCY);
    },
    "streamTodos": (call) => {
        // gRPC Advantage: Sending many items over ONE stream without request overhead
        todos.forEach(todo => {
            call.write(todo);
        });
        call.end();
    }
});

server.bindAsync("0.0.0.0:40000", grpc.ServerCredentials.createInsecure(), (error, port) => {
    if (error) return console.error(error);
    console.log(`gRPC Server (with ${ARTIFICIAL_LATENCY}ms latency) at port ${port}`);
});