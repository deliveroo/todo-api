require File.expand_path('.rake/utils.rb', File.dirname(__FILE__))

# Define a local BIN for project binaries.
BIN          = File.expand_path(File.join ".", "bin")
MIGRATIONS   = File.expand_path("./migrations")

ENV["GOBIN"] = BIN

# Register binaries built from source.
bin :dotenv,      go: "github.com/joho/godotenv/cmd/godotenv"
bin :linter,      go: "github.com/golangci/golangci-lint/cmd/golangci-lint"
bin :migrate,     go: "github.com/johngibb/migrate/cmd/migrate"
bin :modd,        go: "github.com/cortesi/modd/cmd/modd"
bin :todo_api,    go: "./cmd/todo-api"
bin :waitforpg,   go: "./pkg/waitforpg"

task default: %w[build]

desc "Removes build artifacts"
task :clean do
  run "rm -fr #{BIN}"
end

desc "Builds the application"
task :build do
  run "go install ./..."
end

desc "Lints the source code"
task :lint => linter do
  run "#{linter} run"
end

namespace :db do
  desc "Creates the database"
  task :create do
    run "createdb $DATABASE_NAME", env: :local
  end

  desc "Drops the database"
  task :drop do
    run "dropdb --if-exists $DATABASE_NAME", env: :local
  end

  desc "Applies all pending migrations"
  task :migrate => migrate do
    run "#{migrate} up -src #{MIGRATIONS} -conn $DATABASE_URL", env: :local
    Rake::Task["db:schema"].execute
  end

  desc "Dumps the database schema"
  task :schema do
    run "pg_dump $DATABASE_URL --schema-only --no-owner > schema.sql", env: :local
  end

  namespace :migrate do
    desc "Creates a new empty migration"
    task :create, [:name] => migrate do |t, args|
      run "#{migrate} create -src #{MIGRATIONS} #{args.name || 'unnamed'}"
    end

    desc "Displays the current migration status."
    task :status => migrate do
      run "#{migrate} status -src #{MIGRATIONS} -conn $DATABASE_URL", env: :local
    end
  end
end

namespace :docker do
  desc "Starts the docker env for local development"
  task :up do
    run "docker-compose up --detach --remove-orphans", silent: true, env: :local
  end

  desc "Stops the docker env for local development"
  task :down do
    run "docker-compose down", env: :local
  end

  namespace :test do
    desc "Starts the docker env for automated tests"
    task :up do
      run "docker-compose --project-name todo-api-test up --detach --remove-orphans", env: :test
      run "#{waitforpg} $TODO_API_TEST_DATABASE_URL", env: :test
      run "dropdb --if-exists $DATABASE_NAME", env: :test, silent: true
      run "createdb $DATABASE_NAME", env: :test, silent: true
      run "psql -f schema.sql $TODO_API_TEST_DATABASE_URL -v ON_ERROR_STOP=1", env: :test, silent: true
    end

    desc "Stops the docker env for automated tests"
    task :down do
      run "docker-compose --project-name todo-api-test down", env: :test
    end
  end
end

desc "Runs all tests"
task :test => "docker:test:up" do
  run "go test ./...", env: :test
end

namespace :test do
  desc "Connects to the test database"
  task :psql do
    run "psql --dbname $DATABASE_NAME", env: :test
  end

  desc "Runs tests with dependencies"
  task :deps => "docker:test:up" do
    run "go test -count=1 ./selftest ./repo", env: :test
  end

  # Used by modd to run tests just for files that have changed.
  task :modd do |_, args|
    dirmods = ARGV.map { |a| a }[1..] + ["./selftest"]
    run "go test -count=1 #{dirmods.join(" ")}", env: :test
  end

  desc "Runs repo tests"
  task :repo => "docker:test:up" do
    run "go test -count=1 ./repo", env: :test
  end

  desc "Runs selftest"
  task :self => "docker:test:up" do
    run "go test -count=1 ./selftest", env: :test
  end

  desc "Runs short tests"
  task :short do
    run "go test -short ./..."
  end
end

desc "Runs selftest-api"
task :run => "docker:up" do
  run todo_api, env: :local, exec: true
end

desc "Connects to the local database"
task :psql do
  run "psql --dbname $DATABASE_NAME", env: :local
end

desc "Watches and builds on file changes"
task :watch do
  run modd, exec: true
end

# Define "rake help" for those with muscle memory of "make help".
task :help do; run "rake -T"; end
