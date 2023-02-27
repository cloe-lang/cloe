# frozen_string_literal: true

TOTAL_COVERAGE_FILE = 'coverage.txt' # This path is specified by codecov.
BIN_PATH = File.absolute_path 'bin'

task :build do
  %w[cloe clutil].each do |command|
    sh "go build -o bin/#{command} ./src/cmd/#{command}"
  end
end

task :fast_unit_test do
  sh 'go test ./...'
end

task :unit_test do
  coverage_file = "/tmp/cloe-unit-test-#{Process.pid}.coverage"

  sh "echo mode: atomic > #{TOTAL_COVERAGE_FILE}"

  `go list ./src/lib/...`.split.each do |package|
    sh %W[go test
          -covermode atomic
          -coverprofile #{coverage_file}
          #{package}].join ' '

    verbose false do
      if File.exist? coverage_file
        sh "tail -n +2 #{coverage_file} >> #{TOTAL_COVERAGE_FILE}"
        rm coverage_file
      end
    end
  end
end

task command_test: :build do
  sh 'bundler install'
  sh %W[bundler exec cucumber
        -r examples/aruba.rb
        PATH=#{BIN_PATH}:$PATH
        examples].join ' '
end

task :performance_test do
  sh 'go test -v -tags performance -run "^TestPerformance" ./...'
end

task :data_race_test do
  raise 'Architecture is not amd64' unless `uname -m` =~ /x86_64/

  sh 'go test -race ./...'
end

task test: %i[unit_test command_test]

task :bench do
  sh "go test -bench . -run '^$' -benchmem ./..."
end

task :format do
  sh 'go fix ./...'
  sh 'go fmt ./...'
  sh 'rubocop -a'
end

task :lint do
  sh 'gem install rubocop'
  sh 'rubocop'
end

task install: %i[test build] do
  sh 'go get ./...'
end

task default: %i[test build]

task :clean do
  sh 'git clean -dfx'
end
