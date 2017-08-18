Feature: Data types
  Scenario: Use number literals
    Given a file named "main.tisp" with:
    """
    (let x 123)
    (let x -456)
    (let x 123.456)
    (let x -456.789)
    (let x 0xDEADBEEF)
    (let x 01234567)
    """
    Then I successfully run `tisp main.tisp`

  Scenario: Use string literals
    Given a file named "main.tisp" with:
    """
    (let x "foo")
    (let x "Hello, world!")
    (let x "My name is Bob.\\nYour name is not Bob.")
    (let x "Job:\\tProgrammer?")
    """
    Then I successfully run `tisp main.tisp`

  Scenario: Use string literals
    Given a file named "main.tisp" with:
    """
    (write ({"foo" 123 ..{"bar" 456} ..{42 2049} ..{nil true true false}} 42))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    2049
    """
