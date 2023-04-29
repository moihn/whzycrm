using CrmRestApi.ApiModels;
using CrmRestApi.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace CrmRestApi.Controllers;

[ApiController]
[Route("api/[controller]")]
public class VendorController : ControllerBase
{
    private readonly ZyCrmService _crm;
    public VendorController(ZyCrmService crmService) 
    {
        // _dbContextFactory = contextFactory;
        _crm = crmService;
    }

    [HttpGet(Name = "GetVendors")]
    public IEnumerable<Vendor> Get()
    {
        return _crm.GetVendors();

    }

    [HttpGet("{vendorId}", Name = "GetVendorbyId")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public ActionResult<Vendor> Get(int vendorId)
    {
        var vendor = _crm.GetVendorById(vendorId);
        if (vendor == null)
        {
            return NotFound($"Vendor with ID {vendorId} is not found");
        }
        return vendor;
    }

    [HttpPost(Name = "AddVendor")]
    public IEnumerable<Vendor> Post(NewVendor newVendor)
    {
        return _crm.AddVendor(newVendor);
    }

    [HttpPut(Name = "UpdateVendor")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public ActionResult<Vendor> Update(Vendor vendor)
    {
        _crm.UpdateVendor(vendor);
        return Get(vendor.Id);
    }

    [HttpDelete("{vendorId}", Name = "DeleteVendor")]
    public IEnumerable<Vendor> Delete(int vendorId)
    {
        return _crm.DeleteVendorById(vendorId);
    }
}
