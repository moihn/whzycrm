dotnet ef dbcontext scaffold name=ConnectionStrings:Postgres Npgsql.EntityFrameworkCore.PostgreSQL --output-dir DbModels --force --no-build --context ModelContext