using System.ComponentModel.DataAnnotations;

namespace CrmBlazorApp.Data
{
    public class NewClientProductDTO
    {
        [Required]
        public int ClientId { get; set; }

        [Required]
        public string Reference { get; set; } = "";

        [Required]
        public string? Description { get; set; }

        public string? Barcode { get; set; }
        public int InnerQty { get; set; } = 1;
        public int? CartonLength { get; set; }
        public int? CartonHeight { get; set; }
        public int? CartonWidth { get; set; }
        public decimal? GrossWeight { get; set; }
        public decimal? NetWeight { get; set; }
    }
}
