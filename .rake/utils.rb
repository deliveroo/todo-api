# bin creates a function that compiles the given Go package to the bin folder on
# demand, and returns its full path.
#
# Example:
#    bin :mytool, "github.com/my/tool"
#    task :usemytool do
#      sh mytool "-s"
#    end
def bin(name, opts)
  if opts[:go]
    import = opts[:go]
    path = File.join(BIN, import.split('/').last)
    file(path) { run "go install #{import}" }
    define_method(name) do
      Rake::Task[path].invoke
      path
    end
  elsif opts[:brew]
    formula = opts[:brew]
    define_method(name) do
      run "hash #{name} || brew install #{formula}"
      name
    end
  end
end

# run executes the specified command using the shell, similar to sh. It adds
# support for providing env vars via a .env file, running the command silently,
# or running it via exec (to replace the current process).
def run(cmd, opts={})
  cmd = "sh -c '#{cmd.gsub("'"){%q{'"'"'}}}'"
  if opts[:env]
    sh "touch .env.local .env.test.local", verbose: false
    env = '.env'
    env += ".#{opts[:env]}" unless opts[:env] == :local
    cmd = "#{dotenv} -f #{env}.local,#{env} #{cmd}"
  end
  puts cmd if verbose == true
  if opts[:silent]
    out = `#{cmd} 2>&1`
    abort out unless $?.success?
  elsif opts[:exec]
    exec cmd
  else
    system(cmd)
  end
end
