from fastapi import FastAPI, Request, Response, Header
from typing import Optional
from pydantic import BaseModel
from fastapi.responses import JSONResponse
import requests

app = FastAPI()

class Greeter(BaseModel):
    name: str

class GreeterResponse(BaseModel):
    msg: str

# In knative mode functions can return errors via headers
def report_error(code, msg):
  errHdrs = {
    "Direktiv-ErrorCode": code,
    "Direktiv-ErrorMessage": msg,
  }
  return Response(msg, status_code=200, headers=errHdrs, media_type=None)

# In knative mode functions log via a localhost post with their action id
# In this example it is a direct post, an alternative would be to tail a file
# and send the new lines with this post
def log(aid, msg):
    url = 'http://localhost:8889/log?aid={}'.format(aid)
    requests.post(url = url, data = msg)

@app.post("/", response_model=GreeterResponse)
async def root(r: Request, response: Response,
    direktiv_actionid: Optional[str] = Header(None)):

    b = await r.json()

    log(direktiv_actionid, 'I got the name {}'.format(b["name"]))

    try:
        if len(b["name"]) < 3:
            return report_error("greeter.shortname", "name is too short")
        return GreeterResponse(msg='Hello {}'.format(b["name"]))
    except KeyError:
        return report_error("greeter.missing", "name is missing")
