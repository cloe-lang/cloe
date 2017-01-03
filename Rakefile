task :risp do
  sh 'cd src/cmd/risp && go build main.go'
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => :risp

task :clean do
  sh 'git clean -dfx'
end
