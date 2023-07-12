using System.ComponentModel.DataAnnotations;

namespace CrmBlazorApp.Data.AddClientOrder
{
    public class ClientOrderDTO
    {
        public ClientDTO Client { get; set; } = null!;
        [Required]
        public string ClientOrderReference { get; set; } = null!;
        [Required]
        public string OrderReference { get; set; } = null!;
        public DateOnly OrderDate { get; set; }
        public DateOnly CreationDate { get; set; } = DateOnly.FromDateTime(DateTime.Now);
        public DateOnly? ShipmentDate { get; set; }
        [Required]
        public ICollection<ClientOrderItemDTO> ClientOrderItems { get; } = new List<ClientOrderItemDTO>();

        public ClientOrderStatusDTO Status { get; } = null!;
    }
}
