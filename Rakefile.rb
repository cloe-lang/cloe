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
  tmp_dir = 'tmp'
  mkdir_p tmp_dir

  Dir.glob('test/*.tisp') do |file|
    in_file = file.ext '.in'
    expected_out_file = file.ext '.out'
    actual_out_file = File.join(tmp_dir, File.basename(expected_out_file))

    sh %W(bin/tisp #{file}
          #{File.exist?(in_file) ? "< #{in_file}" : ''}
          #{File.exist?(expected_out_file) ? "> #{actual_out_file}" : ''}
    ).join ' '

    sh "diff #{expected_out_file} #{actual_out_file}" \
        if File.exist? expected_out_file
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
