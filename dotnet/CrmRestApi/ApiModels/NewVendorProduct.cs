

using System.ComponentModel.DataAnnotations;

namespace CrmRestApi.ApiModels
{
    public class NewVendorProduct
    {
        [Required]
        public string Reference { get; set; }

        [Required]
        public int    VendorId { get; set; }

        public bool? TestPerformed { get; set; }

        public NewVendorProduct()
        {
            Reference = "";
        }



    }
}
