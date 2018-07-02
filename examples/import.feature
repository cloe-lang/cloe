Feature: Import statement
  Scenario: Import a module
    Given a file named "main.cloe" with:
    """
    (import "./mod")
    (seq! ..(mod.map print [1 2 3 4 5]))
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
    (mod.hello "world")
    """
    And a file named "modules/mod.cloe" with:
    """
    (def (hello name) (print (merge "Hello, " name "!")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Import nested modules
    Given a file named "main.cloe" with:
    """
    (import "./foo/bar")

    (bar.hello)
    """
    And a directory named "foo"
    And a file named "foo/bar.cloe" with:
    """
    (import "./baz")

    (let hello baz.hello)
    """
    And a file named "foo/baz.cloe" with:
    """
    (def (hello) (print "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Import a directory as a module
    Given a file named "main.cloe" with:
    """
    (import "./foo")

    (foo.hello)
    """
    And a file named "foo/module.cloe" with:
    """
    (def (hello) (print "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Import a module with an alternative prefix
    Given a file named "main.cloe" with:
    """
    (import bar "./foo")

    (bar.hello)
    """
    And a file named "foo.cloe" with:
    """
    (def (hello) (print "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Import a module and expand members inside
    Given a file named "main.cloe" with:
    """
    (import . "./foo")

    (hello)
    """
    And a file named "foo.cloe" with:
    """
    (def (hello) (print "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

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
      (print mod1.stdin . end "")
      (print mod2.stdin . end ""))
    """
    And a file named "mod1.cloe" with:
    """
    (import "./mod2")

    (let stdin mod2.stdin)
    """
    And a file named "mod2.cloe" with:
    """
    (import "os")

    (let stdin (os.readStdin))
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
    And a file named ".cloe/src/foo.cloe" with:
    """
    (def (hello)
      (print "Hello, world!"))
    """
    And a file named "main.sh" with:
    """
    CLOE_PATH=$PWD/.cloe cloe main.cloe
    """
    When I successfully run `sh main.sh`
    Then the stdout should contain exactly "Hello, world!"
