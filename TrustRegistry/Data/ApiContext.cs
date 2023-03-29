using System;
using Microsoft.EntityFrameworkCore;
using TrustRegistry.Models;
namespace TrustRegistry.Data
{
	public class ApiContext : DbContext
	{
		public DbSet<PackageRegistry> Registries { get; set; }

		public ApiContext(DbContextOptions<ApiContext> options)
			:base(options)
		{
		}
	}
}

