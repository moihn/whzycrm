using CrmRestApi.DbModels;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Options;


namespace CrmRestApi.Services;

public class ZyCrmService
{
    private readonly IDbContextFactory<DbModels.ModelContext> _dbContextFactory;

    public ZyCrmService(IDbContextFactory<DbModels.ModelContext> contextFactory)
    {
        _dbContextFactory = contextFactory;
    }

    public IEnumerable<ApiModels.Vendor> GetVendors()
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            // return _crm.GetVendors(context);
            foreach (var vendor in context.Vendors.ToList())
            {
                var id = Convert.ToInt32(vendor.Id);
                var name = Convert.ToString(vendor.Name);
                yield return new ApiModels.Vendor(id, name ?? "");
            }
        }
    }

    public ApiModels.Vendor? GetVendorById(int vendorId)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            // return _crm.GetVendors(context);
            var vendor = context.Vendors.Find(Convert.ToDecimal(vendorId));
            if (vendor == null)
            {
                return null;
            }
            return new ApiModels.Vendor(Convert.ToInt32(vendor.Id), vendor.Name ?? "");
        }
    }

    public IEnumerable<ApiModels.Vendor> AddVendor(ApiModels.NewVendor inputVendor)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            var newDbVendor = new DbModels.Vendor();
            newDbVendor.Name = inputVendor.Name;
            context.Vendors.Add(newDbVendor);
            context.SaveChanges();
        }
        return GetVendors();
    }

    public IEnumerable<ApiModels.Vendor> DeleteVendorById(int vendorId)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            // find and remove
            {
                var vendor = context.Vendors.Find(Convert.ToDecimal(vendorId));
                if (vendor != null)
                {
                    context.Vendors.Remove(vendor);
                    context.SaveChanges();
                }
            }
        }
        return GetVendors();
    }

    public void UpdateVendor(ApiModels.Vendor inputVendor)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            var vendor = context.Vendors.Find(Convert.ToDecimal(inputVendor.Id));
            if (vendor != null)
            {
                vendor.Name = inputVendor.Name;
                context.SaveChanges();
            }
        }
    }

    public IEnumerable<ApiModels.VendorProductSummary> GetProductSummariesByVendor(int vendorId)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            var dbProducts = context.VendorProducts.Where(vp => vp.VendorId == vendorId).ToList();
            foreach (var dbProduct in dbProducts)
            {
                var product = new ApiModels.VendorProductSummary();
                product.VendorProductId = product.VendorProductId;
                product.Reference = product.Reference;
                product.VendorId = product.VendorId;
                product.Description = product.Description;
                product.MaterialTypeId = product.MaterialTypeId;
                product.ProductTypeId = product.ProductTypeId;
                product.UnitTypeId = product.UnitTypeId;
                product.Length = product.Length;
                product.Width = product.Width;
                product.Height = product.Height;
                product.Weight = product.Weight;

                yield return product;
            }
        }
    }

    public ApiModels.VendorProduct? GetVendorProductDetail(int vendorProductId)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            var dbProduct = context.VendorProducts.Find(vendorProductId);
            if (dbProduct != null)
            {
                var product = new ApiModels.VendorProduct();
                product.VendorProductId = product.VendorProductId;
                product.Reference = product.Reference;
                product.VendorId = product.VendorId;
                product.Description = product.Description;
                product.MaterialTypeId = product.MaterialTypeId;
                product.ProductTypeId = product.ProductTypeId;
                product.UnitTypeId = product.UnitTypeId;
                product.Length = product.Length;
                product.Width = product.Width;
                product.Height = product.Height;
                product.Weight = product.Weight;
                
                foreach(var vendorPrice in dbProduct.VendorProductPrices.OrderByDescending(s => s.StartDate))
                {
                    var vpp = new ApiModels.VendorProductPrice();
                    vpp.StartDate = vendorPrice.StartDate;
                    vpp.Price = vendorPrice.Price;
                    vpp.CurrencyId = Convert.ToInt32(vendorPrice.CurrencyId);
                    vpp.PriceTypeId = Convert.ToInt32(vendorPrice.PriceTypeId);
                    product.PriceHistory.Add(vpp);
                }
                return product;
            }
            return null;
        }
    }

    public void UploadNewVendorProduct(ApiModels.NewVendorProduct newVendorProduct)
    {
        using (var context = _dbContextFactory.CreateDbContext())
        {
            var newVendorProductRow = new DbModels.VendorProduct();
            newVendorProductRow.Reference = newVendorProduct.Reference;
            if (newVendorProduct.TestPerformed.HasValue)
            {
                // newVendorProductRow.TestPerformed = newVendorProduct.TestPerformed.Value;
            }
        }
    }
}
