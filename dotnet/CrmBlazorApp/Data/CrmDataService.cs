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

        public Task<DbModels.Client?> GetClientAsync(int clientId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var client = dbContext.Clients
                    .SingleOrDefault(row => row.ClientId == clientId);
                return Task.FromResult(client);
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

        public Task<DbModels.ClientProduct?> GetClientProductByIdAsync(int clientProductId)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var product = dbContext.ClientProducts
                    .Include(row => row.ClientProductItems)
                        .ThenInclude(item => item.VendorProduct)
                            .ThenInclude(vp => vp.VendorProductPrices)
                    .SingleOrDefault(row => row.ClientProductId == clientProductId);
                return Task.FromResult(product);
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

        public Task<DbModels.ClientProduct> SaveNewClientProductAsync(Data.NewClientProductDTO productDto, List<DbModels.VendorProduct> vps)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                var product = new DbModels.ClientProduct()
                {
                    ClientId = productDto.ClientId,
                    Reference = productDto.Reference,
                    Description = productDto.Description,
                    Barcode = productDto.Barcode
                };
                foreach (var vp in vps)
                {
                    var cpi = new DbModels.ClientProductItem()
                    {
                        VendorProduct = vp,
                        ClientProduct = product
                    };
                    dbContext.VendorProducts.Attach(vp);
                    dbContext.ClientProductItems.Add(cpi);
                    product.ClientProductItems.Add(cpi);
                }
                
                dbContext.ClientProducts.Add(product);
                dbContext.SaveChanges();
                return Task.FromResult(product);
            }

        }

        public Task<DbModels.ClientProduct> SaveUpdatedClientProductAsync(DbModels.ClientProduct product, List<int> removedVendorProductIds, List<DbModels.VendorProduct> newVendorProducts)
        {
            using (var dbContext = _dbContextFactory.CreateDbContext())
            {
                dbContext.Attach(product);
                foreach (var removedVendorProductId in removedVendorProductIds)
                {
                    var item = product.ClientProductItems.FirstOrDefault(item => item.VendorProductId == removedVendorProductId);
                    if (item != null)
                    {
                        dbContext.ClientProductItems.Remove(item);
                    }
                }
                foreach (var vp in newVendorProducts)
                {
                    var vp2 = dbContext.VendorProducts.SingleOrDefault(item => item.VendorProductId == vp.VendorProductId);
                    if (vp2 != null && !product.ClientProductItems.Select(item => item.VendorProduct.Reference).Contains(vp2.Reference))
                    {
                        var cpi = new DbModels.ClientProductItem()
                        {
                            VendorProduct = vp2,
                            ClientProduct = product
                        };
                        product.ClientProductItems.Add(cpi);
                    }
                }

                dbContext.ClientProducts.Update(product);
                dbContext.SaveChanges();
                return Task.FromResult(product);
            }

        }
    }
}
