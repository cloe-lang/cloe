%i(risp parse).each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end

task :test do
  sh 'go test -test.v ./...'
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => %i(risp parse)

task :clean do
  sh 'git clean -dfx'
end
