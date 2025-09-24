const http = require('http');
const port = 8080;
const routes = {
  '/v1/ocr/jobs': (req, res) => {
    if (req.method === 'POST') {
      let body=''; req.on('data', d => body += d);
      req.on('end', () => {
        const jobId = 'ocr_' + Math.random().toString(36).slice(2,10);
        res.writeHead(200, {'Content-Type':'application/json'});
        res.end(JSON.stringify({ jobId }));
      });
    } else {
      res.writeHead(405); res.end();
    }
  },
};
http.createServer((req, res) => {
  if (req.url.startsWith('/v1/ocr/jobs/') && req.method === 'GET') {
    const jobId = req.url.split('/').pop();
    res.writeHead(200, {'Content-Type':'application/json'});
    res.end(JSON.stringify({ jobId, status: "succeeded", result_uri: "s3://gopan/ocr/demo.json" }));
    return;
  }
  const h = routes[req.url];
  if (h) return h(req, res);
  res.writeHead(404); res.end('not found');
}).listen(port, () => console.log('MinerU mock on :' + port));
