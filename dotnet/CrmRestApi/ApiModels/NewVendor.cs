using System.ComponentModel.DataAnnotations;

namespace CrmRestApi.ApiModels
{
    public class NewVendor
    {
        [Required]
        public string Name { get; set; }
    
        public NewVendor()
        {
            Name = "";
        }

        public NewVendor(string name)
        {
            Name = name;
        }
    }
}
