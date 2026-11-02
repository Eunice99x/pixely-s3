data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "./cmd/atlas-loader"]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev?search_path=public"
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://db/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
