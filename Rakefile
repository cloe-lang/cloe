%i(risp parse).each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end

task :parse_examples => :parse do
  Dir.glob('examples/*.r').each do |file|
    sh "bin/parse #{file}"
  end
end

task :run_examples => :run do
  Dir.glob('examples/*.r').each do |file|
    sh "bin/risp #{file}"
  end
end

task :test do
  sh 'go test -test.v ./...'
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => %i(test risp parse)

task :clean do
  sh 'git clean -dfx'
end
