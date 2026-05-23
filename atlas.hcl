env "local" {
  url = getenv("DATABASE_URL")
  dev = "docker://postgres/16/dev?search_path=public"
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
