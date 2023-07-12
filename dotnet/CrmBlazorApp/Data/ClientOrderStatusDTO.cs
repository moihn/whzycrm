using CrmBlazorApp.Data.EditClientOrder;

namespace CrmBlazorApp.Data
{
    public class ClientOrderStatusDTO
    {
        public int StatusId { get; set; }

        public string Description { get; set; } = null!;
    }
}
