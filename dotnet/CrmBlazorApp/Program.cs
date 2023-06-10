using CrmBlazorApp.Data;
using System.Data.SqlClient;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddRazorPages();
// builder.Services.Configure<string>(builder.Configuration.GetSection("ServerURL"));
builder.Services.AddServerSideBlazor();
builder.Services.AddSingleton<WeatherForecastService>();

builder.Services.AddDbContextFactory<CrmBlazorApp.DbModels.ModelContext>(options =>
    options.UseNpgsql("Name=ConnectionStrings:CrmDb")
);

builder.Services.AddSingleton<CrmDataService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Error");
}


app.UseStaticFiles();

app.UseRouting();

app.MapBlazorHub();
app.MapFallbackToPage("/_Host");

app.Run();
