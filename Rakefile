bins = %i(parse desugar risp)
examples = %i(parse_examples desugar_examples run_examples)

bins.each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end

examples.zip(bins).each do |name, bin|
  task name => bin do
    Dir.glob('examples/*.r').each do |file|
      sh "bin/#{bin} #{file}"
    end
  end
end

task :examples => examples

task :test do
  sh 'go test -test.v ./...'
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => [:test, *bins]

task :clean do
  sh 'git clean -dfx'
end
