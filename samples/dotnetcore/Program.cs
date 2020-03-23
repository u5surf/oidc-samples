namespace OidcSample
{
    using System;
    using Microsoft.AspNetCore.Hosting;
    using Microsoft.Extensions.Hosting;

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
                    string host = Environment.GetEnvironmentVariable("HOST") ?? "0.0.0.0";
                    string port = Environment.GetEnvironmentVariable("PORT") ?? "8000";
                    Uri uri = new Uri("http://" + host + ":" + port);
                    webBuilder.UseStartup<Startup>().UseUrls(uri.ToString());
                });
    }
}
