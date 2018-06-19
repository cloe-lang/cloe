Feature: Others
  Scenario: Run Cloe with an empty source
    Given a file named "main.cloe" with:
    """
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Read a source from stdin
    Given a file named "main.cloe" with:
    """
    (write "Hello, world!")
    """
    When I run the following script:
    """
    cloe < main.cloe
    """
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Run Cloe script with shebang
    Given a file named "main.cloe" with mode "0755" and with:
    """
    #!/usr/bin/env cloe

    (write "Hello, world!")
    """
    When I successfully run `sh -c ./main.cloe`
    Then the stdout should contain exactly "Hello, world!"
