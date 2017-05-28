Feature: Parallelism
  Scenario: Evaluate outputs sequentially
    Given a file named "main.tisp" with:
    """
    (seq
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
    When I successfully run `tisp main.tisp`
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
