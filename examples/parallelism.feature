Feature: Parallelism
  Scenario: Evaluate effects in parallel
    Given a file named "main.cloe" with:
    """
    (write (par [1 2 3] [4 5 6] [7 8 9]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "[7 8 9]"

  Scenario: Evaluate effects sequentially
    Given a file named "main.cloe" with:
    """
    (seq!
      (write 0)
      (write 1)
      (write 2)
      (write 3)
      (write 4)
      (write 5)
      (write 6)
      (write 7)
      (write 8)
      (write 9))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    0
    1
    2
    3
    4
    5
    6
    7
    8
    9
    """

  Scenario: Apply rally function to a infinite list
    Given a file named "main.cloe" with:
    """
    (def (f) [42 ..(f)])
    (let a (f))
    (write (first a))
    (let b (rest a))
    (write (first b))
    (let c (rest b))
    (write (first c))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """
