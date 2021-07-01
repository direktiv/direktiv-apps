## Java

## Code

```java
package greeting;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.InetSocketAddress;
import java.util.List;
import java.io.OutputStream;

import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import java.net.URL;
import java.net.URLConnection;

import org.json.*;

public class Greeting {

    public static void main(String[] args)  throws Exception {
        InetSocketAddress isa = new InetSocketAddress(8080);
        HttpServer server = HttpServer.create(isa, 0);
        server.createContext("/", new GreetingHandler());
        server.setExecutor(null);
        server.start();
        System.out.println("Server started...");
        Runtime.getRuntime().addShutdownHook(new Thread(){
            @Override
                public void run() {
                    server.stop(0);
                } 
        });
    }


    static class GreetingHandler implements HttpHandler {
        
          public void log(String aid, String l)  {
            if (aid != "development" && aid != "Development") {
                try {
                    URL url = new URL("http://localhost:8889/log?aid="+ aid);
                    URLConnection con = url.openConnection();
                    HttpURLConnection http = (HttpURLConnection)con;
                    http.setRequestMethod("POST");
                    http.setRequestProperty("Content-Type", "plain/text");
                    http.getOutputStream().write(l.getBytes("UTF-8"));
                    http.connect();
                } catch(Exception e) {
                    // left empty as it should always work
                    return;
                }
            } else {
                System.out.println(l);
            }
        }

        public void respond(HttpExchange t, String output) throws IOException {
            Headers headers = t.getResponseHeaders();
            headers.set("Content-Type", "application/json");
            t.sendResponseHeaders(200, output.getBytes().length);

            OutputStream os = t.getResponseBody();
            os.write(output.getBytes());
            os.close();
        }

        public void respondWithError(HttpExchange t, String code, String l) throws IOException {
            Headers headers = t.getResponseHeaders();
            headers.set("Direktiv-ErrorCode", code);
            headers.set("Direktiv-ErrorMessage", l);
            t.sendResponseHeaders(200, l.getBytes().length);
            OutputStream os = t.getResponseBody();
            os.write(l.getBytes());
            os.close();
        }

        public void handle(HttpExchange t) throws IOException {
            System.out.println("Handling Request");
            Headers headers = t.getRequestHeaders();
            List<String> aidlist = headers.get("Direktiv-ActionID");
            String aid = aidlist.get(0);

            log(aid, "Reading Input");
            
            InputStreamReader isr = new InputStreamReader(t.getRequestBody(), "utf-8");
            BufferedReader br = new BufferedReader(isr);

            int b;
            StringBuilder buf = new StringBuilder(512);
            while((b = br.read()) != -1) {
                buf.append((char)b);
            }

            br.close();
            isr.close();

            log(aid, "Decoding Input");

            String name = "";

            try{
                System.out.println("Before object convert");
                JSONObject json = new JSONObject(buf.toString());
                System.out.println("After Convert");
                name = json.get("name").toString();
            } catch(Throwable e) {
                System.out.println(e);
                respondWithError(t, "com.greeting-decodejson.error", "Unable to find 'name' object in JSON");
                return;
            }

            log(aid, "Writing Output");

            respond(t, "{ \"greeting\": \"Welcome to Direktiv, "+ name +"!\" }");
        }

    }


}
```

## Dockerfile

```dockerfile
FROM openjdk:17-jdk-alpine3.13 as build

WORKDIR /app

COPY ./greeting/Greeting.java ./greeting/Greeting.java
COPY ./json-java.jar ./
COPY ./manifest.txt ./


RUN javac -classpath "./json-java.jar" ./greeting/Greeting.java
RUN jar cfm Greeting.jar ./manifest.txt ./greeting/*.class
RUN jar tf Greeting.jar


CMD ["java", "-cp", "./Greeting.jar:./json-java.jar", "greeting/Greeting"]
```