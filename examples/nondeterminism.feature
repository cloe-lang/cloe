Feature: Nondeterminism
  Scenario: Apply rally function to a infinite list
    Given a file named "main.coel" with:
    """
    (def (f) (prepend 42 (f)))
    (let a (f))
    (write (first a))
    (let b (rest a))
    (write (first b))
    (let c (rest b))
    (write (first c))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """
