Feature: File System
  Scenario: Import file system module
    Given a file named "main.coel" with:
    """
    (import "fs")
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly ""

  Scenario: Create a directory
    Given a file named "main.coel" with:
    """
    (import "fs")

    (fs.createDirectory "foo")
    """
    When I successfully run `coel main.coel`
    Then I successfully run `rmdir foo`
