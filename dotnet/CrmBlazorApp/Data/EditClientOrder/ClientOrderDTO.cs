using System.ComponentModel.DataAnnotations;

namespace CrmBlazorApp.Data.EditClientOrder
{
    public class ClientOrderDTO
    {
        public int OrderId { get; set; }
        public ClientDTO Client { get; } = null!;
        [Required]
        public string ClientOrderReference { get; set; } = null!;
        [Required]
        public string OrderReference { get; set; } = null!;
        public DateOnly CreationDate { get; set; } = DateOnly.FromDateTime(DateTime.Now);

        [Required]
        public DateOnly? OrderDate { get; set; }

        public bool IsDraft { get; }
        public DateOnly? ShipmentDate { get; set; }
        [Required]
        public ICollection<ClientOrderItemDTO> ClientOrderItems { get; } = new List<ClientOrderItemDTO>();

        public DbModels.ClientOrderStatus Status { get; set; } = null!;

    }
}
