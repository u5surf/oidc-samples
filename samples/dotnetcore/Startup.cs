namespace OidcSample
{
    using System.Linq;
    using System.Text;
    using System.Text.Encodings.Web;
    using Microsoft.AspNetCore.Authentication;
    using Microsoft.AspNetCore.Authentication.OpenIdConnect;
    using Microsoft.AspNetCore.Builder;
    using Microsoft.AspNetCore.Hosting;
    using Microsoft.AspNetCore.Http;
    using Microsoft.Extensions.Configuration;
    using Microsoft.Extensions.DependencyInjection;
    using Microsoft.Extensions.Hosting;
    using Microsoft.IdentityModel.Protocols.OpenIdConnect;

    public class Startup
    {
        private const string MiraclIssuer = "https://api.mpin.io";

        public Startup(IConfiguration configuration)
        {
            this.Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = OpenIdConnectDefaults.AuthenticationScheme;
                options.DefaultSignInScheme = OpenIdConnectDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = OpenIdConnectDefaults.AuthenticationScheme;
                options.DefaultScheme = OpenIdConnectDefaults.AuthenticationScheme;
            })
            .AddCookie(OpenIdConnectDefaults.AuthenticationScheme, options =>
            {
                options.Cookie.SameSite = SameSiteMode.None;
            })
            .AddOpenIdConnect("oidc", options =>
            {
                options.Authority = this.Configuration["ISSUER"] ?? MiraclIssuer;
                options.ClientId = this.Configuration["CLIENT_ID"];
                options.ClientSecret = this.Configuration["CLIENT_SECRET"];
                options.ResponseType = OpenIdConnectResponseType.Code;
                options.AuthenticationMethod = OpenIdConnectRedirectBehavior.RedirectGet;
                options.Scope.Clear();
                options.Scope.Add("openid email profile");
                options.CallbackPath = this.GetCallbackPath();

                options.GetClaimsFromUserInfoEndpoint = true;
                options.SaveTokens = true;

                options.ClaimsIssuer = "oidc";
                options.ProtocolValidator.RequireNonce = false;

                options.ClaimActions.MapJsonKey("email_verified", "email_verified", "bool");
                options.ClaimActions.MapJsonKey("sub", "sub", "string");
                options.Events = new OpenIdConnectEvents()
                {
                    OnAuthenticationFailed = c =>
                    {
                        c.HandleResponse();
                        c.Response.StatusCode = StatusCodes.Status400BadRequest;
                        return c.Response.WriteAsync("Authentication Failed");
                    },
                };
            });

            // Add framework services.
            services.AddControllersWithViews();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseAuthentication();
            app.UseAuthorization();
            app.UseRouting();

            app.Run(async context =>
            {
                var userResult = await context.AuthenticateAsync("oidc");
                var user = userResult.Principal;
                var props = userResult.Properties;

                // Not authenticated.
                if (user == null || !user.Identities.Any(identity => identity.IsAuthenticated))
                {
                    // This is what [Authorize] calls.
                    await context.ChallengeAsync("oidc", new AuthenticationProperties() { RedirectUri = "/" });
                    return;
                }

                var response = context.Response;

                // Return json structured user info for our tests.
                response.ContentType = "text/json";
                StringBuilder rr = new StringBuilder("{");
                foreach (var claim in context.User.Claims)
                {
                    rr.Append(string.Format("\"{0}\":\"{1}\",", claim.Type, claim.Value));
                }

                rr.Remove(rr.Length - 1, 1);
                rr.Append("}");
                await response.WriteAsync(rr.ToString());
            });
        }

        private static string HtmlEncode(string content) =>
            string.IsNullOrEmpty(content) ? string.Empty : HtmlEncoder.Default.Encode(content);

        private string GetCallbackPath()
        {
            if (string.IsNullOrEmpty(this.Configuration["REDIRECT_URL"]))
            {
                return "/login";
            }

            string redirectUrl = this.Configuration["REDIRECT_URL"];
            string callbackPath = redirectUrl.Remove(0, redirectUrl.IndexOf("://") + 3);
            callbackPath = callbackPath.Remove(0, callbackPath.IndexOf("/"));

            return callbackPath;
        }
    }
}
