using Microsoft.AspNetCore.Mvc;
using System;
using System.IO;
using Newtonsoft.Json;
using System.Net;
using System.Text;

namespace dotnet.Controllers
{
    public class InputName {
        public string name { get; set;}
    }

    public class OutputResponse {
        public string greeting { get; set;}
    }

    [Route("/")]
    [ApiController]
    public class GreetingController: ControllerBase
    {
        public void log(String aid,  String l)
        {
            if (aid != "development" && aid != "Development") {
                try{
                    string url = "http://localhost:8889/log?aid="+aid;
                    var data = Encoding.ASCII.GetBytes(l);
                    var request = (HttpWebRequest)WebRequest.Create(url);
                    request.Method = "POST";
                    request.ContentType = "plain/text";
                    using(var stream = request.GetRequestStream())
                    {
                        stream.Write(data, 0, data.Length);
                    }
                    request.GetResponse();
                } catch(Exception e){
                    // ignore
                    Console.WriteLine(e);
                }
            } else {
                Console.WriteLine(l);
            }
        }

        [HttpPost]
        public ActionResult<String> Post()
        {
            using(var reader = new StreamReader(Request.Body))
            {
                string aid = Request.Headers["Direktiv-ActionID"];

                log(aid, "Reading Input");
                var body = reader.ReadToEnd();
                
                InputName input = JsonConvert.DeserializeObject<InputName>(body);
                log(aid, "Decoding Input");
                OutputResponse or = new OutputResponse();
                
                or.greeting = "Welcome to Direktiv, "+ input.name +"!";
                String output = JsonConvert.SerializeObject(or);
                log(aid, "Writing Output");
                return Content(output, "application/json");
            }
        }
    }
}