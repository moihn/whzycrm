using CrmRestApi.Services;
using CrmRestApi.Configuration;
using Microsoft.EntityFrameworkCore;
using System.Data.SqlClient;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.

builder.Services.AddControllers();
var conStrBuilder = new SqlConnectionStringBuilder(
    // builder.Configuration.GetConnectionString("Oracle")
    builder.Configuration.GetConnectionString("Postgres")
);
builder.Services.AddDbContextFactory<CrmRestApi.DbModels.ModelContext>(
    // options => options.UseOracle(conStrBuilder.ConnectionString)
    options => options.UseNpgsql(conStrBuilder.ConnectionString)
);
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
// builder.Services.Configure<OracleSettings>(builder.Configuration.GetSection("OracleSettings"));
builder.Services.AddSingleton<ZyCrmService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseAuthorization();

app.MapControllers();

app.Run();
