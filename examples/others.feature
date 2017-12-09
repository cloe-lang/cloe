Feature: Others
  Scenario: Run Tisp with an empty source
    Given a file named "main.tisp" with:
    """
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly ""

  Scenario: Read a source from stdin
    Given a file named "main.tisp" with:
    """
    (write "Hello, world!")
    """
    When I run the following script:
    """
    tisp < main.tisp
    """
    Then the stdout should contain exactly:
    """
    Hello, world!
    """

  Scenario: Run Tisp script with shebang
    Given a file named "main.tisp" with mode "0755" and with:
    """
    #!/usr/bin/env tisp

    (write "Hello, world!")
    """
    When I successfully run `sh -c ./main.tisp`
    Then the stdout should contain exactly "Hello, world!"

