Feature: Others
  Scenario: Run Coel with an empty source
    Given a file named "main.coel" with:
    """
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly ""

  Scenario: Read a source from stdin
    Given a file named "main.coel" with:
    """
    (write "Hello, world!")
    """
    When I run the following script:
    """
    coel < main.coel
    """
    Then the stdout should contain exactly:
    """
    Hello, world!
    """

  Scenario: Run Coel script with shebang
    Given a file named "main.coel" with mode "0755" and with:
    """
    #!/usr/bin/env coel

    (write "Hello, world!")
    """
    When I successfully run `sh -c ./main.coel`
    Then the stdout should contain exactly "Hello, world!"

