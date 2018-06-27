Feature: Errors
  Scenario: Run an erroneous code
    Given a file named "main.cloe" with:
    """
    (write (+ 1 true))
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"
    And the stderr should contain "main.cloe"
    And the stderr should contain "(write (+ 1 true))"

  Scenario: Bind 2 values to an argument
    Given a file named "main.cloe" with:
    """
    (def (f x)
         (def (g y)
              (def (h z) (+ x y z))
              h)
         g)

    (write (((f 123) 456) . x 0))
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"

  Scenario: Catch an error
    Given a file named "main.cloe" with:
    """
    (write (catch (+ 1 true)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain "name"
    And the stdout should contain "message"

  Scenario: Catch an error passed by match expression
    Given a file named "main.cloe" with:
    """
    (write (@ (catch (match (error "FooError" "") x (error "BarError" ""))) "name"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain "FooError"
