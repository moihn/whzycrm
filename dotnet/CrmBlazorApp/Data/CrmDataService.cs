using CrmBlazorApp.DbModels;
using Microsoft.EntityFrameworkCore;

namespace CrmBlazorApp.Data
{
    public class CrmDataService
    {
        private readonly IDbContextFactory<ModelContext> _dbContextFactory;

        public CrmDataService(IDbContextFactory<ModelContext> factory)
        {
            _dbContextFactory = factory;
        }

        public Task<Client[]> GetClientsAsync()
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var clients = dbContext.Clients.Select(dbClient => new Client
                {
                    Name = dbClient.Name,
                    Id = dbClient.ClientId,
                    CountryId = dbClient.CountryId
                }).ToArray();
                return Task.FromResult(clients);
            }
        }
    }
}
