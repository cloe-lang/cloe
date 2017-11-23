Feature: HTTP
  Scenario: Import HTTP module
    Given a file named "main.tisp" with:
    """
    (import "http")
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly ""

  Scenario: Send GET request
    Given a file named "main.tisp" with:
    """
    (import "http")

    (write ((http.Get "https://google.com") "status"))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly "200"
