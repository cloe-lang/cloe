Feature: Collections
  Scenario: Use collections as a function
    Given a file named "main.coel" with:
    """
    (seq
      (write ([123 [456 789] "foo" true nil false] 1))
      (write ({123 [456 789] "foo" "It's me." nil false} "foo"))
      (write ("Hello, world!" 5)))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [456 789]
    It's me.
    ,
    """

  Scenario: Chain indexing
    Given a file named "main.coel" with:
    """
    (write ({"foo" {"bar" 42}} "foo" "bar"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    42
    """
