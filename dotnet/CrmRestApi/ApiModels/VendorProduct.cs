namespace CrmRestApi.ApiModels
{
    public class VendorProduct : VendorProductSummary
    {
        public class ClientProductItemOrderHistory
        {
            public class OrderedItem
            {
                public DateOnly OrderDate { get; set; }
                public int OrderedQuantity { get; set; }
                public decimal Price { get; set; }
                public int CurrencyId { get; set; }
            }
            public int ClientId { get; set; }
            public ICollection<OrderedItem> ItemOrderHistory = new List<OrderedItem>();
        }

        // This is a vendor product change history sorted by price start date in descending order
        public ICollection<VendorProductPrice> PriceHistory { get; set; } = new List<VendorProductPrice>();

        // This is a client product order history for this product
        public ICollection<ClientProductItemOrderHistory> OrderHistoryByClients = new List<ClientProductItemOrderHistory>();
    }
}
