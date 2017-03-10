bins = %w(tisp-parse tisp)


bins.each do |bin|
  task bin do
    sh "go build -o bin/#{bin} src/cmd/#{bin}/main.go"
  end
end


task :build => bins


task :unit_test do
  sh 'go test -cover -race ./...'
end


test_files = Dir.glob('test/*.tisp')


task :parser_test => 'tisp-parse' do |t|
  test_files.each do |file|
    sh "bin/#{t.source} #{file} > /dev/null"
  end
end


task :interpreter_test => 'tisp' do
  tmp_dir = 'tmp'
  mkdir_p tmp_dir

  test_files.each do |file|
    shell_script = file.ext '.sh'

    if File.exist? shell_script
      sh "sh #{shell_script}"
      next
    end

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

  Dir.glob('test/xfail/*.tisp') do |file|
    sh "! bin/tisp #{file} > /dev/null 2>&1"
  end
end


task :command_test => %i(parser_test interpreter_test)


task :test => %i(unit_test command_test)


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
