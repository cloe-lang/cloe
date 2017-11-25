Feature: Others
  Background:
    Given an executable named "leak_memory.sh" with:
    """
    #!/bin/sh

    set -e

    tisp $1 > /dev/null &
    pid=$!

    sleep 1 # Wait for memory usage to be stable.

    ok=false
    last_mem=0

    for _ in $(seq 10)
    do
      mem=$(ps ho rss $pid)

      if [ $last_mem -ge $mem  ]
      then
        ok=true
        break
      fi

      last_mem=$mem
      sleep 1
    done

    kill $pid
    $ok
    """

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

  Scenario: Ensure no memory leak with infinite effects
    This test succeeds only with Go 1.8 onward because of argument liveness.
    Given a file named "main.tisp" with:
    """
    (let many42 (prepend (write 42) many42))
    ..many42
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0

  Scenario: Ensure no memory leak with deep recursion
    Given a file named "main.tisp" with:
    """
    (def (f n)
      (if (= n 0)
          "OK!"
          (f (- n 1))))

    (write (f 100000000))
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0

  Scenario: Ensure no memory leak with map function
    Given a file named "main.tisp" with:
    """
    (let l (prepend 42 l))

    ..(map write l)
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0
