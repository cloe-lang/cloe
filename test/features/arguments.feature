Feature: Arguments
  Scenario: Define a recursive function
    Given a file named "main.tisp" with:
    """
    (let (func x1 x2 (x3 0) (x4 0) ..args . y1 y2 (y3 0) (y4 1) ..kwargs)
         (+ x1 x2 x3 x4 ..args y1 y2 y3 y4))
    (let foo 1)

    (write (func 1 1 1 ..[1 foo foo] . y1 1 y3 1 foo 100000000 ..{"y2" 1}))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    10
    """
