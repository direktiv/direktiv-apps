## .Net 

## Code

Program.cs
```cs
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Hosting;

namespace dotnet
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureWebHostDefaults(webBuilder =>
                {

                    webBuilder.UseStartup<Startup>();
                    webBuilder.UseUrls("http://0.0.0.0:8080/");
                });
    }
}

```

Startup.cs
```cs
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Server.Kestrel.Core;


namespace dotnet
{
    public class Startup
    {
        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services)
        {
            // If using Kestrel:
            services.Configure<KestrelServerOptions>(options =>
            {
                options.AllowSynchronousIO = true;
            });

            services.AddMvc(option => option.EnableEndpointRouting = false).SetCompatibilityVersion(CompatibilityVersion.Version_3_0);
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseMvc();
        }
    }
}

```

GreetingController.cs
```cs
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
```

## Dockerfile

```dockerfile
FROM mcr.microsoft.com/dotnet/sdk:5.0 AS build-env
WORKDIR /app

COPY *.csproj /
RUN dotnet restore /

COPY . ./
RUN dotnet publish -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:5.0

WORKDIR /app
COPY --from=build-env /app/out .


ENTRYPOINT ["dotnet", "dotnet.dll"]
```