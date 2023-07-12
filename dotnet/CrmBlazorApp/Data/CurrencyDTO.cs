using System.ComponentModel.DataAnnotations;

namespace CrmBlazorApp.Data
{
    public class CurrencyDTO
    {
        public int CurrencyId { get; set; }

        [Required]
        public string? IsoSymbol { get; set; }

        public string? Description { get; set; }
    }
}
