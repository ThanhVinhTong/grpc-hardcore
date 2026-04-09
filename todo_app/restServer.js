const express = require('express');
const app = express();
app.use(express.json());

const todos = [];
const ARTIFICIAL_LATENCY = 20;

app.post('/todos', (req, res) => {
    setTimeout(() => {
        const todoItem = { ...req.body, id: todos.length + 1 };
        todos.push(todoItem);
        res.json(todoItem);
    }, ARTIFICIAL_LATENCY);
});

app.get('/todos', (req, res) => {
    setTimeout(() => {
        res.json({ items: todos });
    }, ARTIFICIAL_LATENCY);
});

app.listen(3000, () => {
    console.log(`REST Server (with ${ARTIFICIAL_LATENCY}ms latency) at port 3000`);
});
