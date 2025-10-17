const express = require("express");
const port = process.env.PORT || 8080;

const SERVER_API = 'https://api.nft.marketplace.200lab.io'

const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

app.use('/v1', createProxyMiddleware({ target: SERVER_API, changeOrigin: true }));
app.listen(port);
