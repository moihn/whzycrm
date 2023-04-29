namespace CrmRestApi.ApiModels
{
    public class VendorProductSummary
    {
        public int VendorProductId { get; set; }

        public string Reference { get; set; } = string.Empty;

        public int VendorId { get; set; }
        public string Description { get; set; } = string.Empty;

        public int? MaterialTypeId { get; set; }
        public int ProductTypeId { get; set; }

        public int UnitTypeId { get; set; }
        public int? Length { get; set; }

        public int? Width { get; set; }
        public int? Height { get; set; }
        public int? Weight { get; set; }
    }
}
