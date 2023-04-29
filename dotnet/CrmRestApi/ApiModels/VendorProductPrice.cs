namespace CrmRestApi.ApiModels
{
    public class VendorProductPrice
    {
        public DateOnly StartDate { get; set; }
        public decimal Price { get; set; }
        public int CurrencyId { get; set; }
        public int PriceTypeId { get; set; }
    }
}
