using CrmBlazorApp.DbModels;
using CrmBlazorApp.Pages;
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

        public Task<DbModels.Client[]> GetClientsAsync()
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var clients = dbContext.Clients.Include(row => row.Country).ToArray();
                return Task.FromResult(clients);
            }
        }

        public Task<DbModels.Vendor[]> GetVendorsAsync()
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var vendors = dbContext.Vendors.ToArray();
                return Task.FromResult(vendors);
            }
        }

        public Task<DbModels.VendorProduct[]> GetProductsOfVendorAsync(int vendorId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var products = dbContext.VendorProducts.Include(row => row.MaterialType).Where(row => row.VendorId == vendorId).ToArray();
                return Task.FromResult(products);
            }
        }

        public Task<DbModels.VendorProduct?> GetVendorProductDetailAsync(int vendorProductId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var product = dbContext.VendorProducts
                    .Include(row => row.MaterialType)
                    .Include(row => row.Vendor)
                    .Include(row => row.VendorProductPrices).ThenInclude(p => p.Currency)
                    .SingleOrDefault(row => row.VendorProductId == vendorProductId);
                return Task.FromResult(product);
            }
        }

        public Task<DbModels.VendorProduct?> GetVendorProductDetailByReferenceAsync(string reference)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var product = dbContext.VendorProducts
                    .Include(row => row.MaterialType)
                    .Include(row => row.Vendor)
                    .Include(row => row.VendorProductPrices).ThenInclude(p => p.Currency)
                    .SingleOrDefault(row => row.Reference == reference);
                return Task.FromResult(product);
            }
        }

    }
}
