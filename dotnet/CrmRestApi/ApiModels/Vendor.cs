namespace CrmRestApi.ApiModels
{
    public class Vendor : NewVendor
    {
        public int Id { get; set; }


        public Vendor() : base()
        {
            Id = -1;
        }

        public Vendor(int id, string name) : base(name)
        {
            Id = id;
        }
    }
}
