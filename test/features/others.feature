Feature: Others
  Scenario: Run Tisp with an empty source
    Given a file named "main.tisp" with:
    """
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly ""

  Scenario: Read a source from stdin
    Given a file named "main.tisp" with:
    """
    (write "Hello, world!")
    """
    When I run the following script:
    """
    tisp < main.tisp
    """
    Then the stdout should contain exactly:
    """
    Hello, world!
    """

  Scenario: Ensure no memory leak
    This test succeeds only with Go 1.8 onward because of argument liveness.
    Given a file named "main.tisp" with:
    """
    (let (many42) (prepend (write 42) (many42)))
    ..(many42)
    """
    When I run the following script:
    """
    tisp main.tisp > /dev/null &
    pid=$!

    sleep 1 # Wait for memory usage to be stable.

    ok=false
    last_mem=0

    for _ in $(seq 10)
    do
      mem=$(ps ho vsz $pid)

      if [ $last_mem -ge $mem  ]
      then
        ok=true
        break
      fi

      last_mem=$mem
      sleep 1
    done &&

    kill $pid &&
    $ok
    """
    Then the exit status should be 0
