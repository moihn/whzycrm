using System.ComponentModel.DataAnnotations;

namespace CrmBlazorApp.Data.EditClientOrder
{
    public class ClientOrderItemDTO
    {
        public int OrderItemId { get; set; }

        [Required]
        public virtual DbModels.ClientProduct ClientProduct { get; set; } = null!;

        [Range(1, int.MaxValue, ErrorMessage = "Quantity cannot be 0")]
        public int Quantity { get; set; }

        [Range(0.01, float.MaxValue, ErrorMessage = "Price cannot be 0 or negative")]
        [DataType(DataType.Currency)]
        public decimal Price { get; set; }

        public DbModels.Currency Currency { get; set; } = null!;

        public DateOnly? AddedDate { get; } = DateOnly.FromDateTime(DateTime.Now);

        public DateOnly? AlternativeShipDate { get; set; }
    }
}
