Feature: Regular Expression
  Scenario: Import regular expression module
    Given a file named "main.coel" with:
    """
    (import "re")
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly ""

  Scenario: Match pattern with strings
    Given a file named "main.coel" with:
    """
    (import "re")

    (eseq
      (write (re.match "ab*c" "ac"))
      (write (re.match "ab+c" "ac"))
      (write (re.match "ab*c" "abbbbbc")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    true
    false
    true
    """

  Scenario: Find pattern in strings
    Given a file named "main.coel" with:
    """
    (import "re")

    (eseq
      (write (re.find "abc" "ac"))
      (write (re.find "ab*c" "ac"))
      (write (re.find "a(bc+)" "abcccc"))
      (write (re.find "a(b*c)" "abbbbc")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    nil
    ["ac"]
    ["abcccc" "bcccc"]
    ["abbbbc" "bbbbc"]
    """

  Scenario: Replace pattern in strings
    Given a file named "main.coel" with:
    """
    (import "re")

    (eseq
      (write (re.replace "abc"    "foo"    "ac"))
      (write (re.replace "ab*c"   "x${0}z" "ac"))
      (write (re.replace "a(bc+)" "x${1}z" "abcccc"))
      (write (re.replace "a(b*c)" "x${1}z" "abbbbc")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    ac
    xacz
    xbccccz
    xbbbbcz
    """
