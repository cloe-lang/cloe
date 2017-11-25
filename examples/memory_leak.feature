Feature: Memory leak
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

  Scenario: Run infinite effects
    This test succeeds only with Go 1.8 onward because of argument liveness.
    Given a file named "main.tisp" with:
    """
    (let many42 (prepend (write 42) many42))
    ..many42
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0

  Scenario: Evaluate deep recursion
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

  Scenario: Apply a map function to an infinite list
    Given a file named "main.tisp" with:
    """
    (let l (prepend 42 l))

    ..(map write l)
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0

  Scenario: Apply a map function to an infinite list of map functions
    Given a file named "main.tisp" with:
    """
    (let l (prepend map l))

    ..(map (\ (x) (write (typeOf x))) l)
    """
    When I run `sh leak_memory.sh main.tisp`
    Then the exit status should be 0
