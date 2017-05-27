total_coverage_file = 'coverage.txt' # This path is specified by codecov.

task :build do
  sh 'go build -o bin/tisp src/cmd/tisp/main.go'
end

task :fast_unit_test do
  sh 'go test ./...'
end

task :unit_test do
  coverage_file = "/tmp/tisp-unit-test-#{Process.pid}.coverage"

  `go list ./...`.split.each do |package|
    sh %W[go test
          -covermode atomic
          -coverprofile #{coverage_file}
          -race
          #{package}].join(' ')

    verbose false do
      if File.exist? coverage_file
        sh "cat #{coverage_file} >> #{total_coverage_file}"
        rm coverage_file
      end
    end
  end
end

task :test_build do
  sh 'go test -c -cover '\
     "-coverpkg $(go list ./... | perl -0777 -pe 's/\\n(.)/,\\1/g') "\
     './src/cmd/tisp'
  mkdir_p 'bin'
  mv 'tisp.test', 'bin'

  File.write 'bin/tisp', [
    '#!/bin/sh',
    'file=/tmp/tisp-test-$$.out',
    'coverage_file=/tmp/tisp-test-$$.coverage',
    %(ARGS="$@" #{File.absolute_path 'bin/tisp.test'} )\
    + '-test.coverprofile $coverage_file > $file &&',
    "cat $file | perl -0777 -pe 's/(.*)PASS.*/\\1/s' &&",
    'rm $file &&',
    "cat $coverage_file >> #{total_coverage_file} &&",
    'rm $coverage_file'
  ].join("\n") + "\n"

  chmod 0o755, 'bin/tisp'
end

task command_test: :test_build do
  tmp_dir = 'tmp'
  mkdir_p tmp_dir

  Dir.glob('test/*.tisp').each do |file|
    shell_script = file.ext '.sh'

    if File.exist? shell_script
      sh "sh #{shell_script}"
      next
    end

    in_file = file.ext '.in'
    expected_out_file = file.ext '.out'
    actual_out_file = File.join(tmp_dir, File.basename(expected_out_file))

    sh %W[
      bin/tisp #{file}
      #{File.exist?(in_file) ? "< #{in_file}" : ''}
      #{File.exist?(expected_out_file) ? "> #{actual_out_file}" : ''}
    ].join ' '

    sh "diff #{expected_out_file} #{actual_out_file}" \
        if File.exist? expected_out_file
  end

  Dir.glob 'test/xfail/*.tisp' do |file|
    sh "! bin/tisp #{file} > /dev/null 2>&1"
  end
end

task test: %i[unit_test command_test]

task :lint do
  verbose false do
    [
      'go vet',
      'golint',
      'gosimple',
      'unused',
      'staticcheck',
      'interfacer'
    ].each do |command|
      puts "# #{command}"
      sh "#{command} ./..."
    end
  end
end

task :format do
  sh 'go fix ./...'
  sh 'go fmt ./...'

  Dir.glob '**/*.go' do |file|
    sh "goimports -w #{file}"
  end

  sh 'rubocop -a'
end

task :docker do
  tag = 'tisplang/tisp-build'
  sh "sudo docker build --no-cache -t #{tag} etc/docker"
  sh "sudo docker push #{tag}"
end

task :install_deps do
  sh %w[
    go get -u
    github.com/golang/lint/golint
    github.com/kr/pretty
    github.com/mvdan/interfacer/...
    golang.org/x/tools/cmd/goimports
    honnef.co/go/tools/...
  ].join ' '

  sh 'go get -d -t ./...'
end

task install: %i[install_deps test build] do
  sh 'go get ./...'
end

task :images do
  Dir.glob 'img/*.svg' do |file|
    sh "inkscape --export-area-drawing --export-png #{file.ext 'png'} #{file}"
  end
end

task doc: :images do
  sh 'convert -resize 16x16 img/icon.png doc/theme/img/favicon.ico'
  cd 'doc'
  sh 'mkdocs gh-deploy -m "[skip ci] on Wercker"'
end

task default: %i[test build]

task :clean do
  sh 'git clean -dfx'
end
