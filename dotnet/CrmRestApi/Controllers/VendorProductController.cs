using CrmRestApi.ApiModels;
using CrmRestApi.Services;
using Microsoft.AspNetCore.Mvc;

// For more information on enabling Web API for empty projects, visit https://go.microsoft.com/fwlink/?LinkID=397860

namespace CrmRestApi.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class VendorProductController : ControllerBase
    {
        private readonly ZyCrmService _crm;
        public VendorProductController(ZyCrmService crmService)
        {
            // _dbContextFactory = contextFactory;
            _crm = crmService;
        }

        // GET: api/<VendorProductController>/
        [HttpGet("Vendor/{vendorId}", Name = "GetByVendorId")]
        public ActionResult<IEnumerable<VendorProductSummary>> GetProductsOfVendor(int vendorId)
        {
            var vendor = _crm.GetVendorById(vendorId);
            if (vendor != null)
            {
                return Ok(_crm.GetProductSummariesByVendor(vendorId));
            }
            return NotFound($"Vendor with ID {vendorId} is not found");
        }

        // GET api/<VendorProductController>/5
        [HttpGet("Product/{vendorProductId}", Name = "GetByVendorProductId")]
        public ActionResult<VendorProduct> GetProductDetail(int vendorProductId)
        {
            var product = _crm.GetVendorProductDetail(vendorProductId);
            if (product != null)
            {
                return Ok(product);
            }
            return NotFound($"Vendor product with ID {vendorProductId} is not found");
        }

        // POST api/<VendorProductController>
        [HttpPost(Name = "UploadNewVendorProduct")]
        public void Post(NewVendorProduct value)
        {
            // var newVendorProd = new
        }

        // PUT api/<VendorProductController>/5
        [HttpPut("{id}")]
        public void Put(int id, [FromBody] string value)
        {
        }

        // DELETE api/<VendorProductController>/5
        [HttpDelete("{id}")]
        public void Delete(int id)
        {
        }
    }
}
