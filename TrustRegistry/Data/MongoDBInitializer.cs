using System.IO;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using MongoDB.Driver;

public static class MongoDBInitializer
{
    public static void Initialize(IServiceCollection services, IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("mongodb+srv://steveswag210:pass123@cluster0.jnqprc8.mongodb.net/test");

        services.AddSingleton<IMongoClient>(new MongoClient(connectionString));

        // Other MongoDB related services can be registered here
        // For example: services.AddScoped<IMongoDatabase>(provider => provider.GetService<IMongoClient>().GetDatabase("my-db-name"));
    }
}