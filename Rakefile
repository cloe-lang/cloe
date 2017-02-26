bins = %i(parse tisp)
examples = %i(parse_examples run_examples)


bins.each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end


task :build => bins


examples.zip(bins).each do |name, bin|
  task name => bin do
    Dir.glob('examples/*.r').each do |file|
      sh "bin/#{bin} #{file}"
    end
  end
end


task :examples => examples


task :unittest do
  sh 'go test ./...'
end


task :cmdtest => :tisp do
  Dir.glob('test/*.tisp') do |file|
    sh "bin/tisp #{file}"
  end
end


task :test => %i(unittest cmdtest)


task :lint do
  sh 'go vet ./...; golint ./...'
end


task :default => %i(test build)


task :clean do
  sh 'git clean -dfx'
end
