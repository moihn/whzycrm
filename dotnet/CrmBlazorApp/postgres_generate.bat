dotnet ef dbcontext scaffold Name=ConnectionStrings:Postgres Npgsql.EntityFrameworkCore.PostgreSQL --output-dir DbModels --force --no-build --context ModelContext