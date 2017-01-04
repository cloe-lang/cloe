%i(risp parse).each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end

[%i(run_examples risp), %i(parse_examples parse)].each do |name, bin|
  task name => bin do
    Dir.glob('examples/*.r').each do |file|
      sh "bin/#{bin} #{file}"
    end
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
