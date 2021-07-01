## Python

## Code

```py
from http.server import BaseHTTPRequestHandler, HTTPServer
import requests
import json
import signal
import sys

PORT = 8080

# Headers
DirektivActionIDHeader     = "Direktiv-ActionID"
DirektivErrorCodeHeader    = "Direktiv-ErrorCode"
DirektivErrorMessageHeader = "Direktiv-ErrorMessage"

InputNameField = "name"

class DirektivHandler(BaseHTTPRequestHandler):
    def _log(self, actionID, msg):
        if actionID != "development" and actionID != "Development":
            try:
                r = requests.post("http://localhost:8889/log?aid=%s" % actionID, headers={"Content-type": "plain/text"}, data = msg)
                if r.status_code != 200:
                    self._send_error("com.greeting-bad-log.error", "log request failed to direktiv")
            except:
                self._send_error("com.greeting-bad-log.error", "failed to log to direktiv")
        else: 
            print(msg)

    def _send_error(self, errorCode, errorMsg):
        self.send_response(400)
        self.send_header('Content-type', 'application/json')
        self.send_header(DirektivErrorCodeHeader, 'application/json')
        self.send_header(DirektivErrorMessageHeader, errorMsg)
        self.end_headers()
        self.wfile.write(json.dumps({"error": errorMsg}).encode())
        return 

    def do_POST(self):
        actionID = ""
        if DirektivActionIDHeader in self.headers:
            actionID = self.headers[DirektivActionIDHeader]
        else:
            return self._send_error("com.greeting-bad-header.error", "Header '%s' must be set" % DirektivActionIDHeader)

        self._log(actionID, "Decoding Input")
        self.data_string = self.rfile.read(int(self.headers['Content-Length']))
        reqData = json.loads(self.data_string)
        
        if InputNameField in reqData:
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()

            # Respond Data
            self._log(actionID, "Writing Output")
            self.wfile.write(json.dumps({"greeting": "Welcome to Direktiv, %s" % reqData[InputNameField]}).encode())
            return
        else:
            return self._send_error("com.greeting-input.error","json field '%s' must be set" % InputNameField)


httpd = HTTPServer(('', PORT), DirektivHandler)
print('Starting greeter server on ":%s"' % PORT)

def shutdown(*args):
    print('Shutting down Server')
    httpd.server_close()
    sys.exit(0)

signal.signal(signal.SIGINT, shutdown)
httpd.serve_forever()
```

## Dockerfile

```dockerfile
FROM python:3

WORKDIR /usr/src/app
COPY main.py ./
COPY requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

CMD [ "python", "./main.py" ]
```