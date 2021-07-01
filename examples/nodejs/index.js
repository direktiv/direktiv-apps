const http = require("http")
const fetch = require("node-fetch")

const host = "localhost"
const port = "8080"

const DirektivActionIDHeader = "direktiv-actionid"
const DirektivErrorCodeHeader = "Direktiv-ErrorCode"
const DirektivErrorMessageHeader = "Direktiv-ErrorMessage"

const greetingHandler = function(req, res) {
    const aid = req.headers[DirektivActionIDHeader]
    var body = ""

    log(aid, "Reading Input")
    
    req.on('data', function(chunk) {
        body += chunk
    })

    req.on('end', function(){
        log(aid, "Decoding Input")
        let json = JSON.parse(body)
        if(!json.name){
            respondWithError(res, `com.greeting-decode.error`, "Unable to find 'name' object in JSON")
            return
        }

        log(aid, "Writing Output")
        respond(res, JSON.stringify({
            "greeting": `Welcome to Direktiv, ${json.name}!`
        }))
    })
}

const log = async (aid, l) => {
    if (aid != "development" && aid != "Development") {
        try {
            fetch(`http://localhost:8889/log?aid=${aid}`, {
                method: "POST",
                body: l,
                headers: {
                    "Content-Type": "plain/text"
                }
            })
        } catch(e) {
            console.log('unable to post log message', e)
        }
    } else {
        console.log(l)
    }
}

const respondWithError = (res, code, err) => {
    res.setHeader(DirektivErrorCodeHeader, code)
    res.setHeader(DirektivErrorMessageHeader, err)
    res.end(JSON.stringify({
        "errorCode": code,
        "errorMsg": err, 
    }))
}

const respond = (res, data) => {
    res.setHeader("Content-Type", "application/json")
    res.end(data)
}

const server = http.createServer(greetingHandler)

server.listen(port, host, () => {
    console.log(`Server is running on http://${host}:${port}`)
})

process.once('SIGTERM', function (code) {
    server.close()
})