Feature: Parallelism
  Scenario: Evaluate effects in parallel
    Given a file named "main.cloe" with:
    """
    (print (par [1 2 3] [4 5 6] [7 8 9]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "[7 8 9]"

  Scenario: Evaluate effects sequentially
    Given a file named "main.cloe" with:
    """
    (seq!
      (print 0)
      (print 1)
      (print 2)
      (print 3)
      (print 4)
      (print 5)
      (print 6)
      (print 7)
      (print 8)
      (print 9))
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

  Scenario: Apply rally function to an infinite list
    Given a file named "main.cloe" with:
    """
    (def (f) [42 ..(f)])
    (let a (f))
    (print (first a))
    (let b (rest a))
    (print (first b))
    (let c (rest b))
    (print (first c))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """
