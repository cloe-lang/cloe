Feature: clutil command
  Scenario: Install modules in a repository
    When I run the following commands:
    """
    CLOE_PATH=$PWD/.cloe clutil install https://github.com/cloe-lang/examples
    """
    Then I run the following commands:
    """
    PATH=$PWD/.cloe/bin:$PATH hello
    """
    And the exit status should be 0
    And the output should contain "Hello, world!"

  Scenario: Update repositories in a language directory
    Given I run the following commands:
    """
    CLOE_PATH=$PWD/.cloe clutil install https://github.com/cloe-lang/examples
    """
    When I run the following commands:
    """
    CLOE_PATH=$PWD/.cloe clutil update
    """
    And the exit status should be 0

  Scenario: Clean up a language directory
    Given I run the following commands:
    """
    CLOE_PATH=$PWD/.cloe clutil install https://github.com/cloe-lang/examples
    """
    And I run the following commands:
    """
    ls $PWD/.cloe
    """
    And the exit status should be 0
    When I run the following commands:
    """
    CLOE_PATH=$PWD/.cloe clutil clean
    """
    Then I run the following commands:
    """
    ls $PWD/.cloe
    """
    And the exit status should not be 0
