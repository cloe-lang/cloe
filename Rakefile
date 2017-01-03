task :build do
  sh 'cd src && go build main.go'
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => :build
