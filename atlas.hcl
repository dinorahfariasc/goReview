env "local" {
  url = getenv("DATABASE_URL")
  dev = "postgres://goreview:goreview@host.docker.internal:5433/dev?sslmode=disable"
  src = "file://db/schema"

  migration {
    dir = "file://db/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}