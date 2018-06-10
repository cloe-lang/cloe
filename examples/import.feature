Feature: Import statement
  Scenario: Import a module
    Given a file named "main.cloe" with:
    """
    (import "./mod")
    (seq! ..(mod.map write [1 2 3 4 5]))
    """
    And a file named "mod.cloe" with:
    """
    (def (map func list)
      (match list
        [] []
        [first ..rest] [(func first) ..(map func rest)]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    1
    2
    3
    4
    5
    """

  Scenario: Import a module in a directory
    Given a file named "main.cloe" with:
    """
    (import "./modules/mod")
    (mod.Hello "world")
    """
    And a file named "modules/mod.cloe" with:
    """
    (def (Hello name) (write (merge "Hello, " name "!")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    Hello, world!
    """

  Scenario: Import a module with invalid path
    Given a file named "main.cloe" with:
    """
    (import "mod")
    """
    And a file named "mod.cloe" with:
    """
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0

  Scenario: Print values in cached modules
    Given a file named "main.cloe" with:
    """
    (import "./mod1")
    (import "./mod2")

    (seq!
      (write mod1.stdin . end "")
      (write mod2.stdin . end ""))
    """
    And a file named "mod1.cloe" with:
    """
    (import "./mod2")

    (let stdin mod2.stdin)
    """
    And a file named "mod2.cloe" with:
    """
    (let stdin (read))
    """
    When I successfully run `sh -c 'echo Hello | cloe main.cloe'`
    Then the stdout should contain exactly:
    """
    Hello
    Hello
    """

  Scenario: Import a module via language path
    Given a file named "main.cloe" with:
    """
    (import "foo")

    (foo.hello)
    """
    And a directory named "cloe-modules"
    And a file named "cloe-modules/foo.cloe" with:
    """
    (def (hello)
      (write "Hello, world!"))
    """
    And a file named "main.sh" with:
    """
    CLOE_PATH=$PWD/cloe-modules cloe main.cloe
    """
    When I successfully run `sh main.sh`
    Then the stdout should contain exactly "Hello, world!"
