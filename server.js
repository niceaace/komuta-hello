const http = require('http');
const port = process.env.PORT || 8080;
const server = http.createServer((req, res) => {
  if (req.url === '/health') { res.writeHead(200, {'Content-Type':'text/plain'}); res.end('ok'); return; }
  res.writeHead(200, {'Content-Type':'text/plain'});
  res.end('hello from komuta\n');
});
server.listen(port, '0.0.0.0', () => console.log('listening on ' + port));








