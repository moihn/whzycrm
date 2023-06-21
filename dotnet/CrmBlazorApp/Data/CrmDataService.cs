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

        public Task<DbModels.ClientProduct[]> GetProductsOfClientAsync(int clientId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var products = dbContext.ClientProducts
                    .Include(row => row.ClientProductItems)
                        .ThenInclude(item => item.VendorProduct)
                            .ThenInclude(vp => vp.VendorProductPrices)
                    .Where(row => row.ClientId == clientId).ToArray();
                return Task.FromResult(products);
            }
        }

        public Task<DbModels.ClientOrder[]> GetOrdersOfClientAsync(int orderId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var orders = dbContext.ClientOrders
                    .Include(row => row.ClientOrderItems)
                        .ThenInclude(item => item.ClientProduct)
                            .ThenInclude(cp => cp.ClientProductItems)
                                .ThenInclude(cpi => cpi.VendorProduct)
                                    .ThenInclude(vp => vp.VendorProductPrices)
                    .Include(row => row.ClientOrderItems)
                        .ThenInclude(item => item.Currency)
                    .Include(row => row.Status)
                    .Where(row => row.ClientId == orderId).ToArray();
                return Task.FromResult(orders);
            }
        }

        public Task<DbModels.ClientOrder?> GetOrderOfClientByOrderIdAsync(int orderId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var order = dbContext.ClientOrders
                    .Include(row => row.ClientOrderItems)
                        .ThenInclude(item => item.ClientProduct)
                            .ThenInclude(cp => cp.ClientProductItems)
                                .ThenInclude(cpi => cpi.VendorProduct)
                                    .ThenInclude(vp => vp.VendorProductPrices)
                    .Include(row => row.ClientOrderItems)
                        .ThenInclude(item => item.Currency)
                    .Include(row => row.Status)
                    .SingleOrDefault(row => row.ClientId == orderId);
                return Task.FromResult(order);
            }
        }
    }
}
