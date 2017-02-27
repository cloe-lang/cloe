bins = %i(parse tisp)


bins.each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end


task :build => bins


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


task :format do
  Dir.glob('src/**/*.go').each do |file|
    sh "go fmt #{file}"
  end
end


task :docker do
  tag = 'raviqqe/tisp-build'
  sh "sudo docker build -t #{tag} etc/docker"
  sh "sudo docker push #{tag}"
end


task :default => %i(test build)


task :clean do
  sh 'git clean -dfx'
end
