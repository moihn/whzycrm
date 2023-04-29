using CrmBlazorApp.DbModels;

namespace CrmBlazorApp.Data
{
    public class CrmDataService
    {
        private readonly ModelContext _context;

        public CrmDataService(ModelContext context)
        {
            _context = context;
        }

        public Task<Client[]> GetClientsAsync()
        {
            var clients = _context.Clients.Select(dbClient => new Client
            {
                Name = dbClient.Name,
                Id = dbClient.ClientId,
                CountryId = dbClient.CountryId
            }).ToArray();
            return Task.FromResult(clients);
        }
    }
}
