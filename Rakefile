%i(risp parse).each do |bin|
  task bin do
    sh "cd src/cmd/#{bin} && go build main.go"
  end
end

task :lint do
  sh 'go vet ./...; golint ./...'
end

task :default => :risp

task :clean do
  sh 'git clean -dfx'
end
