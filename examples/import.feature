Feature: Import statement
  Scenario: Import a module
    Given a file named "main.coel" with:
    """
    (import "./mod")
    (eseq ..(mod.map write [1 2 3 4 5]))
    """
    And a file named "mod.coel" with:
    """
    (def (map func list)
      (match list
        [] []
        [first ..rest] (prepend (func first) (map func rest))))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    1
    2
    3
    4
    5
    """

  Scenario: Import a module in a directory
    Given a file named "main.coel" with:
    """
    (import "./modules/mod")
    (mod.Hello "world")
    """
    And a file named "modules/mod.coel" with:
    """
    (def (Hello name) (write (merge "Hello, " name "!")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Hello, world!
    """

  Scenario: Import a module with invalid path
    Given a file named "main.coel" with:
    """
    (import "mod")
    """
    And a file named "mod.coel" with:
    """
    """
    When I run `coel main.coel`
    Then the exit status should not be 0

  Scenario: Print values in cached modules
    Given a file named "main.coel" with:
    """
    (import "./mod1")
    (import "./mod2")

    (eseq
      (write mod1.stdin . end "")
      (write mod2.stdin . end ""))
    """
    And a file named "mod1.coel" with:
    """
    (import "./mod2")

    (let stdin mod2.stdin)
    """
    And a file named "mod2.coel" with:
    """
    (let stdin (read))
    """
    When I successfully run `sh -c 'echo Hello | coel main.coel'`
    Then the stdout should contain exactly:
    """
    Hello
    Hello
    """
