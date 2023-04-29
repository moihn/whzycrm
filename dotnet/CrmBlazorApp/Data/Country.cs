namespace CrmBlazorApp.Data
{
    public class Country
    {
        public int Id { get; set; }
        public string Name { get; set; }

        public Country(int id, string name)
        {
            Id = id;
            Name = name;
        }
    }
}
