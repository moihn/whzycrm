using AutoMapper;
using CrmBlazorApp.Data;

namespace CrmBlazorApp
{
    public class AutoMapperProfile : Profile
    {
        public AutoMapperProfile()
        {
            CreateMap<DbModels.ClientOrder, Data.EditClientOrder.ClientOrderDTO>();
            CreateMap<DbModels.Client, ClientDTO>();
            CreateMap<DbModels.ClientOrderItem, Data.EditClientOrder.ClientOrderItemDTO>();

        }
    }
}
