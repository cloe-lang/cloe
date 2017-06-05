TOTAL_COVERAGE_FILE = 'coverage.txt'.freeze # This path is specified by codecov.
BIN_PATH = File.absolute_path 'bin'

def go_test(*args)
  sh %W[go test
        -covermode atomic
        #{`uname -m` =~ /x86_64/ ? '-race' : ''}
        #{args.join ' '}].join ' '
end

task :build do
  sh 'go build -o bin/tisp src/cmd/tisp/main.go'
end

task :fast_unit_test do
  sh 'go test ./...'
end

task :unit_test do
  coverage_file = "/tmp/tisp-unit-test-#{Process.pid}.coverage"

  `go list ./src/lib/...`.split.each do |package|
    go_test '-coverprofile', coverage_file, package

    verbose false do
      if File.exist? coverage_file
        sh "cat #{coverage_file} >> #{TOTAL_COVERAGE_FILE}"
        rm coverage_file
      end
    end
  end
end

task command_test: :build do
  cd 'test' do
    sh 'bundler install'
    sh "bundler exec cucumber PATH=#{BIN_PATH}:$PATH"
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

task :doc do
  cd 'doc'
  sh 'mkdocs gh-deploy -m "[skip ci] on Wercker"'
end

task default: %i[test build]

task :clean do
  sh 'git clean -dfx'
end
